package project

import "time"

// EstadoAl calcula el estado del proyecto con sus fechas respecto de hoy (RN-04.1, TEC-20): antes del inicio está
// Planificado, entre el inicio y el fin (ambos incluidos) En curso, y después del fin Finalizado. Compara por día.
func (p Proyecto) EstadoAl(hoy time.Time) Estado {
	d := dia(hoy)
	switch {
	case d.Before(dia(p.Inicio)):
		return Planificado
	case d.After(dia(p.Fin)):
		return Finalizado
	}
	return EnCurso
}

// DiasRestantes devuelve los días que faltan desde hoy hasta la fecha de fin, sin contar hoy (RN-04.2, TEC-20). Si
// hoy es el día de fin o ya pasó, devuelve 0.
func (p Proyecto) DiasRestantes(hoy time.Time) int {
	// Se cuentan días de calendario en UTC para que un cambio de horario no corra el resultado.
	desde := time.Date(hoy.Year(), hoy.Month(), hoy.Day(), 0, 0, 0, 0, time.UTC)
	hasta := time.Date(p.Fin.Year(), p.Fin.Month(), p.Fin.Day(), 0, 0, 0, 0, time.UTC)
	if dias := int(hasta.Sub(desde).Hours() / 24); dias > 0 {
		return dias
	}
	return 0
}

// dia devuelve la fecha sin la hora.
func dia(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
