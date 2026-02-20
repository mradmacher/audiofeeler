package internal

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHome(t *testing.T) {
	teardown, _ := setupDbTest(t)
	defer teardown(t)

	dbEngine, err := NewDbEngine(os.Getenv("AUDIOFEELER_TEST_DATABASE_URL"))
	if err != nil {
		t.Errorf("Error creating the DbEngine: %v", err)
	}
	app := NewApp(NewTemplateEngine("../views"), dbEngine)

	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	app.router.ServeHTTP(res, req)

	if res.Code != http.StatusSeeOther {
		t.Errorf("Expected response code %d; got %d\n", http.StatusSeeOther, res.Code)
	}
}
