package server_test

import (
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestHandleBinCreate(t *testing.T) {
    srv, _, err := NewTestingServer("./test.sqlite")
    assert.NoError(t, err)

    // Test case: Successful bin creation
    t.Run("successful bin creation", func(t *testing.T) {
        id := uuid.New().String()
        form := strings.NewReader("uuid=" + id + "&password=test&content=test_content")
        req, _ := http.NewRequest("POST", "/create", form)
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusFound, resp.Code)
        assert.Contains(t, resp.Header().Get("Location"), id)
    })

    // Test case: Missing UUID
    t.Run("missing uuid", func(t *testing.T) {
        form := strings.NewReader("password=test&content=test_content")
        req, _ := http.NewRequest("POST", "/create", form)
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "empty_uuid")
    })

    // Test case: Invalid UUID
    t.Run("invalid uuid", func(t *testing.T) {
        form := strings.NewReader("uuid=invalid&password=test&content=test_content")
        req, _ := http.NewRequest("POST", "/create", form)
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "invalid_bin_id")
    })

    // Test case: Missing password
    t.Run("missing password", func(t *testing.T) {
        id := uuid.New().String()
        form := strings.NewReader("uuid=" + id + "&content=test_content")
        req, _ := http.NewRequest("POST", "/create", form)
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "empty_password")
    })

    // Test case: Missing content
    t.Run("missing content", func(t *testing.T) {
        id := uuid.New().String()
        form := strings.NewReader("uuid=" + id + "&password=test")
        req, _ := http.NewRequest("POST", "/create", form)
        req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "empty_content")
    })
}
