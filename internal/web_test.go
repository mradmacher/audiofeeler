package audiofeeler

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHome(t *testing.T) {
    app, err := NewApp("../views")
    if err != nil {
      t.Errorf("Error creating the app: %v", err)
    }

    app.MountHandlers()

    req := httptest.NewRequest("GET", "/", nil)
    res := httptest.NewRecorder()
    app.router.ServeHTTP(res, req)

    if res.Code != http.StatusOK {
      t.Errorf("Expected response code %d; got %d\n", http.StatusOK, res.Code)
    }
}
