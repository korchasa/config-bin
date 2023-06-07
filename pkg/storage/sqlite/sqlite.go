package sqlite

import (
    "configBin/pkg"
    "configBin/pkg/storage"
    "database/sql"
    "fmt"
    "github.com/google/uuid"
    _ "github.com/mattn/go-sqlite3"
    log "github.com/sirupsen/logrus"
    "time"
)

type Storage struct {
    db        *sql.DB
    encryptor storage.Encryptor
}

func NewSqliteStorage(dbPath string, encryptor storage.Encryptor) (*Storage, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    return &Storage{
        db:        db,
        encryptor: encryptor,
    }, nil
}

func (s *Storage) CreateBin(id uuid.UUID, pass string, unencryptedData string) error {
    ed, err := s.encryptor.Encrypt(unencryptedData, pass)
    if err != nil {
        return fmt.Errorf("failed to encrypt data: %v", err)
    }

    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.Exec("INSERT INTO bins (uuid, current_version) VALUES (?, 1)", id)
    if err != nil {
        return err
    }

    _, err = tx.Exec("INSERT INTO configurations (uuid, data, version, created_at) VALUES (?, ?, 1, ?)", id, ed, time.Now())
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *Storage) GetBin(id uuid.UUID, pass string) (*pkg.Bin, error) {
    tx, err := s.db.Begin()
    if err != nil {
        return nil, err
    }
    defer tx.Rollback()

    bin, err := s.getAndDecrypt(id, pass, tx)
    if err != nil {
        return nil, err
    }

    err = tx.Commit()

    return bin, err
}

func (s *Storage) getAndDecrypt(id uuid.UUID, pass string, tx *sql.Tx) (*pkg.Bin, error) {
    bin := &pkg.Bin{}
    row := tx.QueryRow("SELECT current_version FROM bins WHERE uuid = ?", id)
    err := row.Scan(&bin.CurrentVersion)
    if err != nil {
        return nil, err
    }

    rows, err := tx.Query("SELECT data, version, created_at, format FROM configurations WHERE uuid = ? ORDER BY version", id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        config := pkg.Configuration{}
        err = rows.Scan(&config.EncryptedData, &config.Version, &config.CreatedAt, &config.Format)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        config.UnencryptedData, err = s.encryptor.Decrypt(config.EncryptedData, pass)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt data: %v", err)
        }
        bin.History = append(bin.History, config)
    }

    return bin, nil
}

func (s *Storage) UpdateBin(id uuid.UUID, pass string, unencryptedData string) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    bin, err := s.getAndDecrypt(id, pass, tx)
    if err != nil {
        return err
    }

    _, err = tx.Exec("UPDATE bins SET current_version = ? WHERE uuid = ?", bin.CurrentVersion+1, id)
    if err != nil {
        return err
    }

    ed, err := s.encryptor.Encrypt(unencryptedData, pass)
    if err != nil {
        return fmt.Errorf("failed to encrypt data: %v", err)
    }

    _, err = tx.Exec("INSERT INTO configurations (uuid, data, version, created_at) VALUES (?, ?, ?, ?)", id, ed, bin.CurrentVersion+1, time.Now())
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *Storage) RollbackBin(id uuid.UUID, pass string, version int) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    bin, err := s.getAndDecrypt(id, pass, tx)
    if err != nil {
        return err
    }
    if version > bin.CurrentVersion {
        return fmt.Errorf("version %d is greater than current version %d", version, bin.CurrentVersion)
    }

    _, err = tx.Exec("UPDATE bins SET current_version = ? WHERE uuid = ?", version, id)
    if err != nil {
        return err
    }

    _, err = s.db.Exec("DELETE FROM configurations WHERE uuid = ? AND version > ?", id, version)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *Storage) Close() {
    err := s.db.Close()
    if err != nil {
        log.Warnf("error closing database: %v", err)
    }
}

func (s *Storage) IsReady() bool {
    err := s.db.Ping()
    if err != nil {
        log.Errorf("failed to ping db: %v", err)
        return false
    }
    return true
}
