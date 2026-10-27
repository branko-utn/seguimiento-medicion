package sprint

import (
	"errors"
	"testing"
)

func TestCerrar(t *testing.T) {
	s := Sprint{ID: "S1", Estado: Activo}
	if err := s.Cerrar(10, 3); err != nil {
		t.Fatalf("Cerrar: %v", err)
	}
	if s.Estado != Cerrado || s.SPPlanificados != 10 || s.SPCompletados != 3 {
		t.Errorf("sprint = %+v, se esperaba Cerrado con 10 planificados y 3 completados", s)
	}
}

func TestCerrar_errores(t *testing.T) {
	casos := []struct {
		nombre  string
		estado  Estado
		quiere  error
		mensaje string
	}{
		{"ya cerrado (RN-13.4)", Cerrado, ErrSprintCerrado, "el sprint está cerrado y no se puede modificar"},
		{"planificado (RN-13.1)", Planificado, ErrNoActivo, "solo se puede cerrar un sprint activo"},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			s := Sprint{ID: "S1", Estado: c.estado, SPPlanificados: 1, SPCompletados: 1}
			err := s.Cerrar(10, 3)
			if !errors.Is(err, c.quiere) || err.Error() != c.mensaje {
				t.Fatalf("error = %v, se esperaba %q", err, c.mensaje)
			}
			if s.Estado != c.estado || s.SPPlanificados != 1 || s.SPCompletados != 1 {
				t.Errorf("con error el sprint no debe cambiar: %+v", s)
			}
		})
	}
}
