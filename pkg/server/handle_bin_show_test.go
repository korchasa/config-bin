package server_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleBinShow(t *testing.T) {
	srv, store, err := NewTestingServer("./test.sqlite")
	assert.NoError(t, err)

	binID := uuid.New()
	err = store.CreateBin(binID, "test", "test_content")
	assert.NoError(t, err)

	// Test case: Successful bin show
	t.Run("successful bin show", func(t *testing.T) {
		req := formRequestWithCookie(formRequestSpec{
			method:         "GET",
			path:           "/" + binID.String(),
			cookieBid:      binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "test_content")
	})

	// Test case: Invalid bin ID
	t.Run("invalid bin id", func(t *testing.T) {
		req := formRequestWithCookie(formRequestSpec{
			method:         "GET",
			path:           "/invalid",
			cookieBid:      binID,
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Contains(t, resp.Body.String(), "invalid_bin_id")
	})

	// Test case: Not existed bin ID
	t.Run("not existed bin id", func(t *testing.T) {
		req := formRequestWithCookie(formRequestSpec{
			method:         "GET",
			path:           "/9570f2e0-d5c8-4003-93fb-dbd60b54c2df",
			cookieBid:      uuid.MustParse("9570f2e0-d5c8-4003-93fb-dbd60b54c2df"),
			cookiePassword: "test",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.Contains(t, resp.Body.String(), "cant_get_bin")
	})

	// Test case: Missing password cookie
	t.Run("missing password cookie", func(t *testing.T) {
		req := formRequestWithCookie(formRequestSpec{
			method:         "GET",
			path:           "/" + binID.String(),
			cookieBid:      binID,
			cookiePassword: "",
		})
		resp := httptest.NewRecorder()

		srv.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Contains(t, resp.Body.String(), "Enter the bin password")
	})
}
