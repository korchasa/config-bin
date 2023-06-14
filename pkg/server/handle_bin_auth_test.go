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
        req := requestWithForm(formRequestSpec{
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
        req := requestWithForm(formRequestSpec{
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
        req := requestWithForm(formRequestSpec{
            method:   "POST",
            path:     "/00000000-0000-0000-0000-000000000000/auth",
            formData: "password=test&text=test_content",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusNotFound, resp.Code)
        assert.Contains(t, resp.Body.String(), "no_bin_by_id")
    })

    // Test case: Missing password
    t.Run("missing password", func(t *testing.T) {
        req := requestWithForm(formRequestSpec{
            method:   "POST",
            path:     "/" + binID + "/auth",
            formData: "text=test_content",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "empty_password")
    })

    // Test case: Wrong password
    t.Run("wrong password", func(t *testing.T) {
        req := requestWithForm(formRequestSpec{
            method:   "POST",
            path:     "/" + binID + "/auth",
            formData: "password=wrong_password",
        })
        resp := httptest.NewRecorder()

        srv.ServeHTTP(resp, req)

        assert.Equal(t, http.StatusBadRequest, resp.Code)
        assert.Contains(t, resp.Body.String(), "cant_get_bin")
    })
}
