package server

import (
	"net/http/httptest"
	"testing"
)

// pedir hace un pedido al router completo y devuelve la respuesta grabada.
func pedir(t *testing.T, metodo, ruta string) *httptest.ResponseRecorder {
	t.Helper()
	router, err := NewRouter(Servicios{})
	if err != nil {
		t.Fatalf("NewRouter: %v", err)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(metodo, ruta, nil))
	return rec
}
