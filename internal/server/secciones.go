package server

// seccion es una entrada de la navegación principal.
type seccion struct {
	Ruta        string
	Nombre      string
	Descripcion string
	Epica       string
	Sprint      int
	// Icono es el dibujo SVG de la sección en la portada (el atributo d de un <path>, en una grilla de 24×24).
	// TEC-21: así la plantilla no tiene que comparar textos para elegirlo.
	Icono string
	// EnConstruccion muestra una página de aviso en Ruta. Al implementar la sección, ponerlo en false
	// y registrar las rutas reales en NewRouter (si no, la ruta queda registrada dos veces y el servidor no arranca).
	EnConstruccion bool
}

// secciones define la navegación principal, en el orden en que se muestra.
var secciones = []seccion{
	{Ruta: "/proyectos", Icono: "M3 7h6l2 2h10v10H3z", Nombre: "Proyectos", Descripcion: "Alta y modificación de proyectos, integrantes, fechas y estado.", Epica: "E1", Sprint: 1},
	{Ruta: "/backlog", Icono: "M9 6h12M9 12h12M9 18h12M4 6h.01M4 12h.01M4 18h.01", Nombre: "Backlog", Descripcion: "Historias con prioridad, estado, Story Points y criterios de aceptación.", Epica: "E2", Sprint: 1},
	{Ruta: "/sprints", Icono: "M4 5h16v15H4zM4 10h16M9 3v4M15 3v4", Nombre: "Sprints", Descripcion: "Sprint Goal, asignación de historias, cierre e historial.", Epica: "E3", Sprint: 1},
	{Ruta: "/planning-poker", Icono: "M7 3h10v18H7zM10 8h4M10 16h4", Nombre: "Planning Poker", Descripcion: "Estimación con votos ocultos, rondas y estimación acordada.", Epica: "E4", Sprint: 2, EnConstruccion: false},
	{Ruta: "/esfuerzo", Icono: "M12 3a9 9 0 1 0 0 18a9 9 0 1 0 0-18M12 7v5l3 3", Nombre: "Esfuerzo", Descripcion: "Horas trabajadas por integrante y actividad.", Epica: "E5", Sprint: 1, EnConstruccion: false},
	{Ruta: "/defectos", Icono: "M12 3l9 16H3zM12 9v4M12 16h.01", Nombre: "Defectos", Descripcion: "Severidad, estado y sprints de detección y resolución.", Epica: "E6", Sprint: 1, EnConstruccion: false},
	{Ruta: "/dashboard", Icono: "M4 20V11M10 20V4M16 20v-7M2 20h20", Nombre: "Dashboard", Descripcion: "Métricas del proyecto y gráficos.", Epica: "E8", Sprint: 3, EnConstruccion: false},
	{Ruta: "/reportes", Icono: "M6 2h9l5 5v15H6zM14 2v6h6M9 13h6M9 17h6", Nombre: "Reportes", Descripcion: "Reporte de proyecto o de sprint, exportable a PDF.", Epica: "E9", Sprint: 3, EnConstruccion: false},
}
