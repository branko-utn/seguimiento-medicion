package server

// seccion es una entrada de la navegación principal.
type seccion struct {
	Ruta        string
	Nombre      string
	Descripcion string
	Epica       string
	Sprint      int
	// EnConstruccion muestra una página de aviso en Ruta. Al implementar la sección, ponerlo en false
	// y registrar las rutas reales en NewRouter (si no, la ruta queda registrada dos veces y el servidor no arranca).
	EnConstruccion bool
}

// secciones define la navegación principal, en el orden en que se muestra.
var secciones = []seccion{
	{Ruta: "/proyectos", Nombre: "Proyectos", Descripcion: "Alta y modificación de proyectos, integrantes, fechas y estado.", Epica: "E1", Sprint: 1},
	{Ruta: "/backlog", Nombre: "Backlog", Descripcion: "Historias con prioridad, estado, Story Points y criterios de aceptación.", Epica: "E2", Sprint: 1},
	{Ruta: "/sprints", Nombre: "Sprints", Descripcion: "Sprint Goal, asignación de historias, cierre e historial.", Epica: "E3", Sprint: 1},
	{Ruta: "/planning-poker", Nombre: "Planning Poker", Descripcion: "Estimación con votos ocultos, rondas y estimación acordada.", Epica: "E4", Sprint: 2, EnConstruccion: false},
	{Ruta: "/esfuerzo", Nombre: "Esfuerzo", Descripcion: "Horas trabajadas por integrante y actividad.", Epica: "E5", Sprint: 1, EnConstruccion: false},
	{Ruta: "/defectos", Nombre: "Defectos", Descripcion: "Severidad, estado y sprints de detección y resolución.", Epica: "E6", Sprint: 1, EnConstruccion: false},
	{Ruta: "/dashboard", Nombre: "Dashboard", Descripcion: "Métricas del proyecto y gráficos.", Epica: "E8", Sprint: 3, EnConstruccion: true},
	{Ruta: "/reportes", Nombre: "Reportes", Descripcion: "Reporte de proyecto o de sprint, exportable a PDF.", Epica: "E9", Sprint: 3, EnConstruccion: true},
}
