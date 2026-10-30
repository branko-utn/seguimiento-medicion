package service

import (
	"math"
	"sort"
	"time"

	"github.com/branko-utn/seguimiento-medicion/internal/dashboard"
	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/backlog"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/defect"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/sprint"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// Dashboard es el caso de uso del dashboard del proyecto (HU-28). Solo lee: junta los datos de los repositorios y
// los servicios que ya existen, y internal/dashboard arma la vista.
type Dashboard struct {
	proyectos    RepoProyectos
	sprints      RepoSprints
	historias    RepoHistorias
	estimaciones RepoEstimaciones
	registros    RepoRegistros
	defectos     RepoDefectos
}

func NuevoDashboard(proyectos RepoProyectos, sprints RepoSprints, historias RepoHistorias,
	estimaciones RepoEstimaciones, registros RepoRegistros, defectos RepoDefectos) *Dashboard {
	return &Dashboard{proyectos: proyectos, sprints: sprints, historias: historias,
		estimaciones: estimaciones, registros: registros, defectos: defectos}
}

// Del arma el dashboard de un proyecto (HU-28). Devuelve project.ErrNoExiste si el proyecto no existe.
func (s *Dashboard) Del(proyectoID domain.ProyectoID, hoy time.Time) (dashboard.Dashboard, error) {
	p, err := s.proyectos.Buscar(proyectoID)
	if err != nil {
		return dashboard.Dashboard{}, err
	}

	res := dashboard.Dashboard{
		ProyectoID: p.ID,
		Nombre:     p.Nombre,
		Inicio:     p.Inicio,
		Fin:        p.Fin,
		Estado:     p.EstadoAl(hoy).String(),
		DiasRestantes: p.DiasRestantes(hoy),
	}

	// Historias del proyecto
	todas := s.historias.Listar()
	var historiasProyecto []backlog.Historia
	for _, h := range todas {
		if h.ProyectoID == proyectoID {
			historiasProyecto = append(historiasProyecto, h)
		}
	}
	totalHistorias := len(historiasProyecto)
	completadasProyecto := 0
	for _, h := range historiasProyecto {
		if h.Estado == backlog.Completada {
			completadasProyecto++
		}
	}
	if totalHistorias == 0 {
		res.PorcentajeCompletadas = 0
	} else {
		if pct, err := metrics.PorcentajeCompletadas(completadasProyecto, totalHistorias); err == nil {
			res.PorcentajeCompletadas = math.Round(pct*10) / 10
		} else {
			res.PorcentajeCompletadas = 0
		}
	}

	// Sprints
	sprintsList := s.sprints.Listar()
	// sprint activo
	var sprintActivo *sprint.Sprint
	var sprintsCerrados []sprint.Sprint
	for _, sp := range sprintsList {
		if sp.ProyectoID != proyectoID {
			continue
		}
		switch sp.Estado {
		case sprint.Activo:
			c := sp
			sprintActivo = &c
		case sprint.Cerrado:
			sprintsCerrados = append(sprintsCerrados, sp)
		}
	}
	if sprintActivo != nil {
		res.SprintActivo = sprintActivo.ID
	}

	// Velocidad: promedio de SP completados de sprints cerrados
	completadosVel := make([]int, 0, len(sprintsCerrados))
	for _, sp := range sprintsCerrados {
		completadosVel = append(completadosVel, sp.SPCompletados)
	}
	res.Velocidad = metrics.Velocidad(completadosVel)

	// Métricas del proyecto (HU-26, HU-27)
	// Horas estimadas y reales de todas las historias del proyecto
	var horasEstimadasProyecto float64
	historiaSet := map[domain.HistoriaID]backlog.Historia{}
	for _, h := range historiasProyecto {
		historiaSet[h.ID] = h
		if e, ok := s.estimaciones.Buscar(h.ID); ok {
			horasEstimadasProyecto += e.Horas
		}
	}
	// Reales: todos los registros de historias del proyecto
	var realesProyecto []metrics.RegistroHoras
	for _, reg := range s.registros.Listar() {
		if _, ok := historiaSet[reg.HistoriaID]; ok {
			realesProyecto = append(realesProyecto, metrics.RegistroHoras{Fecha: reg.Fecha, Horas: reg.Horas})
		}
	}
	res.HorasEstimadas = horasEstimadasProyecto
	res.HorasReales = metrics.HorasReales(func() []float64 {
		var hs []float64
		for _, r := range realesProyecto {
			hs = append(hs, r.Horas)
		}
		return hs
	}())
	if d, err := metrics.Desviacion(res.HorasEstimadas, res.HorasReales); err == nil {
		res.Desviacion = math.Round(d*10) / 10
		res.DesviacionCalculable = true
	} else {
		res.DesviacionCalculable = false
	}

	// Defectos del proyecto
	defsProyecto := s.defectosDelProyecto(proyectoID)
	var ds []metrics.DefectoSprint
	for _, d := range defsProyecto {
		ds = append(ds, metrics.DefectoSprint{SprintDeteccion: string(d.SprintDeteccion), SprintResolucion: string(d.SprintResolucion)})
	}
	res.DefectosDetectados = metrics.TotalDetectados(ds)
	res.DefectosResueltos = metrics.TotalResueltos(ds)
	res.DefectosAbiertos = res.DefectosDetectados - res.DefectosResueltos
	if res.DefectosAbiertos < 0 {
		res.DefectosAbiertos = 0
	}

	// Series por sprint cerrado (ordenados por fecha de inicio)
	sort.Slice(sprintsCerrados, func(i, j int) bool {
		if sprintsCerrados[i].Inicio.Equal(sprintsCerrados[j].Inicio) {
			return sprintsCerrados[i].ID < sprintsCerrados[j].ID
		}
		return sprintsCerrados[i].Inicio.Before(sprintsCerrados[j].Inicio)
	})

	if len(sprintsCerrados) == 0 {
		res.Mensajes = dashboard.Mensajes{
			Velocidad: dashboard.MensajeSinSprintsCerrados,
			SP:        dashboard.MensajeSinSprintsCerrados,
			Horas:     dashboard.MensajeSinSprintsCerrados,
			Defectos:  dashboard.MensajeSinSprintsCerrados,
		}
		return res, nil
	}

	res.Puntos = make([]dashboard.PuntoSprint, 0, len(sprintsCerrados))
	// Calcular por sprint cerrado
	for _, sp := range sprintsCerrados {
		pt := dashboard.PuntoSprint{
			Sprint:         sp.ID,
			SPPlanificados: sp.SPPlanificados,
			SPCompletados:  sp.SPCompletados,
		}
		// Horas del sprint: historias planificadas en el sprint
		horasEstimadasSprint := 0.0
		// Mapear historias del sprint
		histS := map[domain.HistoriaID]bool{}
		for _, hid := range sp.Historias {
			histS[hid] = true
			if e, ok := s.estimaciones.Buscar(hid); ok {
				horasEstimadasSprint += e.Horas
			}
		}
		pt.HorasEstimadas = horasEstimadasSprint
		// Horas reales en periodo del sprint
		var regsSprint []metrics.RegistroHoras
		for _, reg := range s.registros.Listar() {
			if histS[reg.HistoriaID] {
				regsSprint = append(regsSprint, metrics.RegistroHoras{Fecha: reg.Fecha, Horas: reg.Horas})
			}
		}
		pt.HorasReales = metrics.HorasEnPeriodo(regsSprint, sp.Inicio, sp.Fin)
		// Defectos del sprint
		pt.DefectosDetectados = metrics.DefectosDetectados(dsSprint(defsProyecto, sp.ID), string(sp.ID))
		pt.DefectosResueltos = metrics.DefectosResueltos(dsSprint(defsProyecto, sp.ID), string(sp.ID))
		res.Puntos = append(res.Puntos, pt)
	}

	// Mensajes por gráfico según reglas
	// Sin horas estimadas ni registradas en los sprints cerrados
	sinHoras := true
	for _, p := range res.Puntos {
		if p.HorasEstimadas != 0 || p.HorasReales != 0 {
			sinHoras = false
			break
		}
	}
	if sinHoras {
		res.Mensajes.Horas = dashboard.MensajeSinHoras
	}
	// Sin defectos en los sprints cerrados
	sinDefectos := true
	for _, p := range res.Puntos {
		if p.DefectosDetectados != 0 || p.DefectosResueltos != 0 {
			sinDefectos = false
			break
		}
	}
	if sinDefectos {
		res.Mensajes.Defectos = dashboard.MensajeSinDefectos
	}

	return res, nil
}

func (s *Dashboard) defectosDelProyecto(proyecto domain.ProyectoID) []defect.Defecto {
	delProyecto := make(map[domain.HistoriaID]bool)
	for _, h := range s.historias.Listar() {
		if h.ProyectoID == proyecto {
			delProyecto[h.ID] = true
		}
	}
	var ds []defect.Defecto
	for _, d := range s.defectos.Listar() {
		if delProyecto[d.HistoriaID] {
			ds = append(ds, d)
		}
	}
	return ds
}

func dsSprint(defs []defect.Defecto, sprintID domain.SprintID) []metrics.DefectoSprint {
	var out []metrics.DefectoSprint
	for _, d := range defs {
		if d.SprintDeteccion == sprintID || d.SprintResolucion == sprintID {
			out = append(out, metrics.DefectoSprint{SprintDeteccion: string(d.SprintDeteccion), SprintResolucion: string(d.SprintResolucion)})
		}
	}
	return out
}
