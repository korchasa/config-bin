package server_test

import (
	"configBin/pkg/server/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleAPIGetBin(t *testing.T) {
	srv, store, err := NewTestingServer("./test.sqlite")
	assert.NoError(t, err)

	binID := uuid.New()
	err = store.CreateBin(binID, "test", "test_content")
	assert.NoError(t, err)

	// Test case: Successful bin retrieval
	t.Run("successful bin retrieval", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/"+binID.String(), nil)
		req.Header.Add(utils.PasswordHeader, "test")
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "test_content")
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/invalid", nil)
		req.Header.Add(utils.PasswordHeader, "test")
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/00000000-0000-0000-0000-000000000000", nil)
		req.Header.Add(utils.PasswordHeader, "test")
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_get_bin")
	})

	// Test case: Missing password
	t.Run("missing password", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/"+binID.String(), nil)
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "empty_password")
	})
}
