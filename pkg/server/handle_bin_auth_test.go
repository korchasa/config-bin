package server_test

import (
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandleBinAuth(t *testing.T) {
    srv, store, err := NewTestingServer("./test.sqlite")
    assert.NoError(t, err)

    binID := uuid.New().String()
    err = store.CreateBin(uuid.MustParse(binID), "test", "test_content")
    assert.NoError(t, err)
    log.Warnf("binID: %s", binID)

    // Test case: Successful bin authentication
    t.Run("successful bin authentication", func(t *testing.T) {
        req := formRequest(formRequestSpec{
            method:   "POST",
            path:     "/" + binID + "/auth",
            formData: "password=test",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusFound, resp.Code)
        assert.Equal(t, "/"+binID, resp.Header().Get("Location"))
    })

    // Test case: Invalid bin ID
    t.Run("invalid bin id", func(t *testing.T) {
        req := formRequest(formRequestSpec{
            method:   "POST",
            path:     "/invalid/auth",
            formData: "password=test&text=test_content",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "invalid_bin_id")
    })

    // Test case: Not existed bin ID
    t.Run("not existed bin id", func(t *testing.T) {
        req := formRequest(formRequestSpec{
            method:   "POST",
            path:     "/9570f2e0-d5c8-4003-93fb-dbd60b54c2df/auth",
            formData: "password=test&text=test_content",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusNotFound, resp.Code)
        assert.Contains(t, resp.Body.String(), "no_bin_by_id")
    })

    // Test case: Missing password
    t.Run("missing password", func(t *testing.T) {
        req := formRequest(formRequestSpec{
            method:   "POST",
            path:     "/" + binID + "/auth",
            formData: "text=test_content",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "empty_password")
    })
}
