package server

import "net/http"

// health informa que la aplicación está levantada (TEC-01).
func health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"estado":"ok"}`))
}
