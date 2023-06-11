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

func (s *Storage) ApplySchema() error {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS bins
        (
            uuid       TEXT,
            data       TEXT,
            version    INTEGER,
            created_at TIMESTAMP
        );
        CREATE INDEX IF NOT EXISTS bins_uuid_version_index ON bins (uuid, version);
        `)
    if err != nil {
        return fmt.Errorf("failed to create schema: %v", err)
    }
    return nil
}

func (s *Storage) CreateBin(id uuid.UUID, pass string, unencryptedData string) error {
    ed, err := s.encryptor.Encrypt(unencryptedData, pass)
    if err != nil {
        return fmt.Errorf("failed to encrypt data: %v", err)
    }

    _, err = s.db.Exec("INSERT INTO bins (uuid, data, version, created_at) VALUES (?, ?, 0, ?)", id, ed, time.Now())
    if err != nil {
        return fmt.Errorf("failed to insert bin: %v", err)
    }
    return nil
}

func (s *Storage) GetBin(id uuid.UUID, pass string) (*pkg.Bin, error) {
    bin := &pkg.Bin{}
    rows, err := s.db.Query("SELECT uuid, data, version, created_at FROM bins WHERE uuid = ? ORDER BY version DESC", id)
    if err != nil {
        return nil, err
    }
    defer func() {
        err := rows.Close()
        if err != nil {
            log.Errorf("failed to close rows: %v", err)
        }
    }()

    if rows.Next() {
        var data string
        err = rows.Scan(&bin.ID, &data, &bin.Version, &bin.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan first row: %v", err)
        }
        bin.Data, err = s.encryptor.Decrypt(data, pass)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt first row data: %v", err)
        }
    } else {
        return nil, nil
    }

    bin.History = make([]pkg.Bin, 0)
    for rows.Next() {
        prevBin := pkg.Bin{}
        prevData := ""
        err = rows.Scan(&prevBin.ID, &prevData, &prevBin.Version, &prevBin.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %v", err)
        }
        prevBin.Data, err = s.encryptor.Decrypt(prevData, pass)
        if err != nil {
            return nil, fmt.Errorf("failed to decrypt data of version `%d`: %v", bin.Version, err)
        }
        bin.History = append(bin.History, prevBin)
    }

    return bin, nil
}

func (s *Storage) UpdateBin(id uuid.UUID, pass string, unencryptedData string) error {
    bin, err := s.GetBin(id, pass)
    if err != nil {
        return fmt.Errorf("failed to get bin: %v", err)
    }
    if bin == nil {
        return fmt.Errorf("bin not found")
    }

    data, err := s.encryptor.Encrypt(unencryptedData, pass)
    if err != nil {
        return fmt.Errorf("failed to encrypt data: %v", err)
    }

    _, err = s.db.Exec("INSERT INTO bins (uuid, data, version, created_at) VALUES (?, ?, (SELECT COALESCE(MAX(version), 0) + 1 FROM bins WHERE uuid = ?), ?)", id, data, id, time.Now())
    return err
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
