package server_test

import (
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleBinUpdate(t *testing.T) {
	srv, store, err := NewTestingServer("./test.sqlite")
	assert.NoError(t, err)

	binID := uuid.New()
	err = store.CreateBin(binID, "test", "test_content")
	assert.NoError(t, err)
	log.Warnf("binID: %s", binID)

	// Test case: Successful bin update
	t.Run("successful bin update", func(t *testing.T) {
		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/" + binID.String() + "/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusFound, resp.Code)
		assert.Equal(t, "/"+binID.String(), resp.Header().Get("Location"))
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/invalid/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/00000000-0000-0000-0000-000000000000/update",
			formData:       "content=updated_content",
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_update_bin")
	})

	// Test case: Missing content
	t.Run("missing content", func(t *testing.T) {
		req := requestWithFormAndCookie(formRequestSpec{
			method:         "POST",
			path:           "/" + binID.String() + "/update",
			formData:       "",
			cookieBinID:    binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "no_content")
	})
}
