package server

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/project"
	"github.com/branko-utn/seguimiento-medicion/internal/service"
	"github.com/branko-utn/seguimiento-medicion/internal/store/memory"
)

// escenarioEsfuerzo arma dos proyectos: "Demo" (P1), con la integrante Ana Pérez (I1) y la historia HU-1, y "Otro"
// (P2), con el integrante Beto Gómez (I2). Fija hoy en el 24/10/2026.
func escenarioEsfuerzo(t *testing.T) (http.Handler, Servicios) {
	t.Helper()
	original := hoy
	hoy = func() time.Time { return time.Date(2026, 10, 24, 0, 0, 0, 0, time.Local) }
	t.Cleanup(func() { hoy = original })

	proyectos, _ := memory.NuevoProyectos("")
	integrantes, _ := memory.NuevoIntegrantes("")
	historias, _ := memory.NuevoHistorias("")
	sprints, _ := memory.NuevoSprints("")
	estimaciones, _ := memory.NuevoEstimaciones("")
	registros, _ := memory.NuevoRegistros("")
	sprintsSvc := service.NuevoSprints(sprints, proyectos, historias)
	svc := Servicios{
		Proyectos:    service.NuevoProyectos(proyectos),
		Integrantes:  service.NuevoIntegrantes(integrantes, proyectos),
		Backlog:      service.NuevoBacklog(historias, proyectos, sprintsSvc),
		Sprints:      sprintsSvc,
		Estimaciones: service.NuevoEstimaciones(estimaciones, historias),
		Registros:    service.NuevoRegistros(registros, integrantes, historias, proyectos),
	}
	dia := func(d int) time.Time { return time.Date(2026, 10, d, 0, 0, 0, 0, time.Local) }
	for _, p := range []struct{ proyecto, integrante string }{{"Demo", "Ana Pérez"}, {"Otro", "Beto Gómez"}} {
		creado, err := svc.Proyectos.Crear(p.proyecto, "", dia(21), dia(31))
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Integrantes.Registrar(creado.ID, p.integrante, project.ProductBuilder); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := svc.Backlog.Crear("P1", "Crear proyecto", "", backlog.PrioridadAlta); err != nil {
		t.Fatal(err)
	}
	router, err := NewRouter(svc)
	if err != nil {
		t.Fatal(err)
	}
	return router, svc
}

func registroValido() url.Values {
	return url.Values{"integrante": {"I1"}, "historia": {"HU-1"}, "fecha": {"2026-10-24"}, "actividad": {"Desarrollo"}, "horas": {"3"}}
}

func TestEsfuerzo_seccionImplementada(t *testing.T) {
	for _, s := range secciones {
		if s.Ruta == "/esfuerzo" && s.EnConstruccion {
			t.Fatal("la sección /esfuerzo sigue marcada en construcción")
		}
	}
}

func TestEsfuerzo_pantallaVacia(t *testing.T) {
	router, _ := escenarioEsfuerzo(t)

	rec := ver(router, "/esfuerzo")

	if rec.Code != http.StatusOK {
		t.Fatalf("código = %d, se esperaba 200", rec.Code)
	}
	cuerpo := rec.Body.String()
	contiene(t, cuerpo, "<h1>Esfuerzo</h1>", "Todavía no hay horas registradas.")
	contiene(t, formularioDe(t, cuerpo, "/esfuerzo"),
		`<option value="I1">Ana Pérez · Demo</option>`, `<option value="I2">Beto Gómez · Otro</option>`,
		`<option value="HU-1">HU-1 · Crear proyecto</option>`,
		`name="fecha" value="2026-10-24"`, // hoy, precargado
		`name="actividad"`, `name="horas"`,
	)
	contiene(t, formularioDe(t, cuerpo, "/esfuerzo/historias/HU-1/estimacion"), `name="horas"`)
	contiene(t, cuerpo, "Sin estimar")
}

func TestEsfuerzo_registrar(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)

	rec := enviarFormulario(router, "/esfuerzo", registroValido(), false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/esfuerzo" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /esfuerzo", rec.Code, rec.Header().Get("Location"))
	}
	if n := len(svc.Registros.DeHistoria("HU-1")); n != 1 {
		t.Fatalf("HU-1 tiene %d registros, se esperaba 1", n)
	}
	contiene(t, ver(router, "/esfuerzo").Body.String(), "24/10/2026", "Ana Pérez", "HU-1", "Desarrollo", "<td>3</td>")
}

func TestEsfuerzo_registrarHorasConComa(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)
	campos := registroValido()
	campos.Set("horas", "1,5")

	if rec := enviarFormulario(router, "/esfuerzo", campos, false); rec.Code != http.StatusSeeOther {
		t.Fatalf("código = %d, se esperaba 303", rec.Code)
	}
	if r := svc.Registros.DeHistoria("HU-1"); len(r) != 1 || r[0].Horas != 1.5 {
		t.Errorf("registros = %+v, se esperaba uno de 1,5 horas", r)
	}
}

func TestEsfuerzo_registrarConError(t *testing.T) {
	casos := []struct {
		nombre  string
		campo   string
		valor   string
		mensaje string
	}{
		{"actividad vacía", "actividad", "  ", "la actividad es obligatoria"},
		{"horas en cero", "horas", "0", "las horas deben ser mayores que cero"},
		{"horas que no son un número", "horas", "tres", "las horas deben ser mayores que cero"},
		{"fecha futura", "fecha", "2026-10-25", "la fecha no puede ser futura"},
		{"fecha vacía", "fecha", "", "la fecha es obligatoria"},
		{"integrante de otro proyecto", "integrante", "I2", "el integrante no pertenece al proyecto de la historia"},
		{"integrante inexistente", "integrante", "I9", "el integrante no existe"},
		{"historia inexistente", "historia", "HU-99", "la historia no existe"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			router, svc := escenarioEsfuerzo(t)
			campos := registroValido()
			campos.Set(c.campo, c.valor)

			rec := enviarFormulario(router, "/esfuerzo", campos, false)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("código = %d, se esperaba 400", rec.Code)
			}
			form := formularioDe(t, rec.Body.String(), "/esfuerzo")
			contiene(t, form, c.mensaje)
			if c.campo != "actividad" {
				contiene(t, form, `value="Desarrollo"`) // vuelve a mostrar lo que se escribió
			}
			if n := len(svc.Registros.DeHistoria("HU-1")); n != 0 {
				t.Errorf("se guardaron %d registros, no debía guardarse ninguno", n)
			}
		})
	}
}

func TestEsfuerzo_estimar(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)

	rec := enviarFormulario(router, "/esfuerzo/historias/HU-1/estimacion", url.Values{"horas": {"4,5"}}, false)

	if rec.Code != http.StatusSeeOther || rec.Header().Get("Location") != "/esfuerzo" {
		t.Fatalf("respuesta = %d (Location %q), se esperaba 303 a /esfuerzo", rec.Code, rec.Header().Get("Location"))
	}
	e, ok, err := svc.Estimaciones.Buscar("HU-1")
	if err != nil || !ok || e.Horas != 4.5 {
		t.Fatalf("estimación de HU-1 = %+v (encontrada %v, error %v), se esperaba 4,5 horas", e, ok, err)
	}
	contiene(t, formularioDe(t, ver(router, "/esfuerzo").Body.String(), "/esfuerzo/historias/HU-1/estimacion"), `value="4.5"`)
}

func TestEsfuerzo_estimarConError(t *testing.T) {
	router, svc := escenarioEsfuerzo(t)

	rec := enviarFormulario(router, "/esfuerzo/historias/HU-1/estimacion", url.Values{"horas": {"0"}}, false)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("código = %d, se esperaba 400", rec.Code)
	}
	contiene(t, formularioDe(t, rec.Body.String(), "/esfuerzo/historias/HU-1/estimacion"), "las horas deben ser mayores que cero", `value="0"`)
	if _, ok, _ := svc.Estimaciones.Buscar("HU-1"); ok {
		t.Error("la estimación rechazada no debe guardarse")
	}
}

func TestEsfuerzo_estimarHistoriaInexistente(t *testing.T) {
	router, _ := escenarioEsfuerzo(t)

	rec := enviarFormulario(router, "/esfuerzo/historias/HU-99/estimacion", url.Values{"horas": {"4"}}, false)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("código = %d, se esperaba 404", rec.Code)
	}
	contiene(t, rec.Body.String(), "la historia no existe")
}
