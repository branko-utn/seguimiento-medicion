// Comando server levanta la aplicación web en http://localhost:8080 (o en el puerto de la variable PORT). Los datos se
// guardan en archivos JSON en la carpeta data/ (o en la de la variable DATA_DIR).
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/app"
	"github.com/branko-utn/seguimiento-medicion/internal/server"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	dirDatos := "data"
	if d := os.Getenv("DATA_DIR"); d != "" {
		dirDatos = d
	}
	svc, err := app.NuevosServicios(dirDatos)
	if err != nil {
		log.Fatalf("no se pudieron abrir los datos en %s: %v", dirDatos, err)
	}

	router, err := server.NewRouter(svc)
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
