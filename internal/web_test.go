package audiofeeler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHome(t *testing.T) {
	teardown, _ := setupTest(t)
	defer teardown(t)

	app, err := NewApp("../views", os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
	if err != nil {
		t.Errorf("Error creating the app: %v", err)
	}

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	app.router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Expected response code %d; got %d\n", http.StatusOK, res.Code)
	}
}
