// Comando server levanta la aplicación web en http://localhost:8080 (o en el puerto de la variable PORT).
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/server"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	router, err := server.NewRouter(server.Servicios{})
	if err != nil {
		log.Fatalf("no se pudo iniciar la aplicación: %v", err)
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Software Metrics & Estimation escuchando en http://localhost%s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
