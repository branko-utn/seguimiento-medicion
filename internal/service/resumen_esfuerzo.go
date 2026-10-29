package service

import (
	"sort"
	"strconv"
	"strings"

	"github.com/branko-utn/seguimiento-medicion/internal/domain"
	"github.com/branko-utn/seguimiento-medicion/internal/domain/effort"
	"github.com/branko-utn/seguimiento-medicion/internal/metrics"
)

// ResumenEsfuerzo es el esfuerzo registrado en un proyecto (HU-21).
type ResumenEsfuerzo struct {
	Registros     []RegistroEsfuerzo
	PorHistoria   []TotalHistoria
	PorIntegrante []TotalIntegrante
	Total         float64
}

// RegistroEsfuerzo es un registro de esfuerzo con el nombre de su integrante.
type RegistroEsfuerzo struct {
	effort.Registro
	Integrante string
}

// TotalHistoria es el total de horas registradas en una historia.
type TotalHistoria struct {
	Historia domain.HistoriaID
	Horas    float64
}

// TotalIntegrante es el total de horas registradas por un integrante.
type TotalIntegrante struct {
	Integrante domain.IntegranteID
	Nombre     string
	Horas      float64
}

// DelProyecto devuelve el resumen del esfuerzo registrado en un proyecto (HU-21), o project.ErrNoExiste. Junta los
// registros de las historias del proyecto (RN-21.1) y arma los totales con metrics.HorasReales (RN-21.3). Las
// historias y los integrantes sin registros no aparecen (RN-21.4). No modifica nada.
func (s *Registros) DelProyecto(id domain.ProyectoID) (ResumenEsfuerzo, error) {
	if _, err := s.proyectos.Buscar(id); err != nil {
		return ResumenEsfuerzo{}, err
	}

	var resumen ResumenEsfuerzo
	horasPorIntegrante := map[domain.IntegranteID][]float64{}
	nombres := map[domain.IntegranteID]string{}
	var todas []float64
	for _, h := range s.historias.Listar() {
		if h.ProyectoID != id {
			continue
		}
		registros := s.registros.DeHistoria(h.ID)
		if len(registros) == 0 {
			continue
		}
		var horas []float64
		for _, r := range registros {
			nombre, ok := nombres[r.IntegranteID]
			if !ok {
				i, err := s.integrantes.Buscar(r.IntegranteID)
				if err != nil {
					return ResumenEsfuerzo{}, err
				}
				nombre = i.Nombre
				nombres[r.IntegranteID] = nombre
			}
			resumen.Registros = append(resumen.Registros, RegistroEsfuerzo{Registro: r, Integrante: nombre})
			horas = append(horas, r.Horas)
			horasPorIntegrante[r.IntegranteID] = append(horasPorIntegrante[r.IntegranteID], r.Horas)
		}
		resumen.PorHistoria = append(resumen.PorHistoria, TotalHistoria{Historia: h.ID, Horas: metrics.HorasReales(horas)})
		todas = append(todas, horas...)
	}
	for integrante, horas := range horasPorIntegrante {
		resumen.PorIntegrante = append(resumen.PorIntegrante,
			TotalIntegrante{Integrante: integrante, Nombre: nombres[integrante], Horas: metrics.HorasReales(horas)})
	}
	resumen.Total = metrics.HorasReales(todas)

	// RN-21.2: del más reciente al más antiguo; en la misma fecha, por orden de carga (el número del ID).
	sort.SliceStable(resumen.Registros, func(i, j int) bool {
		a, b := resumen.Registros[i], resumen.Registros[j]
		if !a.Fecha.Equal(b.Fecha) {
			return a.Fecha.After(b.Fecha)
		}
		return numeroDeID(string(a.ID)) < numeroDeID(string(b.ID))
	})
	// RN-21.3: las historias por su número (HU-2 antes que HU-10) y los integrantes por nombre.
	sort.Slice(resumen.PorHistoria, func(i, j int) bool {
		return numeroDeID(string(resumen.PorHistoria[i].Historia)) < numeroDeID(string(resumen.PorHistoria[j].Historia))
	})
	sort.Slice(resumen.PorIntegrante, func(i, j int) bool {
		return resumen.PorIntegrante[i].Nombre < resumen.PorIntegrante[j].Nombre
	})
	return resumen, nil
}

// numeroDeID devuelve el número de un ID generado por los repositorios ("HU-10" → 10, "E3" → 3), para ordenarlos
// como números y no como texto.
func numeroDeID(id string) int {
	n, _ := strconv.Atoi(strings.TrimLeft(id, "ABCDEFGHIJKLMNOPQRSTUVWXYZ-"))
	return n
}
