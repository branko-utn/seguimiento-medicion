# De issue a spec: guía para arrancar una historia

Cada historia empieza igual: **spec SDD → escenarios BDD → tests en rojo → código**. Esta guía cubre los dos primeros
pasos. El ejemplo completo de una historia de dominio es HU-12
([`specs/HU-12-registrar-historias-completadas.md`](../../specs/HU-12-registrar-historias-completadas.md)).

## 1. Preparar la rama (5')

```bash
git switch main && git pull
git switch -c feature/HU-XX-titulo-corto
```

En el tablero, mover el issue a **In progress**.

## 2. Pedir un borrador a la IA (10')

La consigna permite usar IA para escribir specs, pero **vos sos responsable de lo que se sube**: tenés que
entenderlo y poder explicarlo en la defensa. Copiá este prompt, completá lo que está entre `<…>` y pegá los
archivos que pide.

```text
Sos parte de un equipo que desarrolla en Go una aplicación de estimación y seguimiento de proyectos
(TPI de Ingeniería y Calidad de Software, UTN). Usamos SDD, BDD y TDD.

Necesito el borrador de la especificación SDD y de los escenarios BDD de esta historia de usuario:

<pegar el issue completo: historia, criterios de aceptación y notas>

Reglas:
1. Usá exactamente la plantilla de spec que pego abajo, en español. Completá: objetivo, entradas, salidas
   esperadas, reglas de negocio (RN-XX.n), restricciones, casos límite, condiciones de error y criterios de
   aceptación. Cada criterio CA-XX.n del issue debe aparecer y tener al menos un escenario.
2. Respetá los tipos y nombres del ADR de contratos de dominio que pego abajo. No inventes tipos de otras
   épicas: referencialas por ID.
3. No agregues funcionalidades que el issue no pide. Si algo es ambiguo, no lo decidas: listalo al final
   como "Dudas para el equipo".
4. Escribí el archivo .feature con "# language: es", etiquetas @HU-XX, @CA-XX.n y @normal / @alternativo /
   @limite / @error, cubriendo como mínimo un caso normal, uno de error y uno límite.
5. Usá pasos concretos y reutilizables (por ejemplo: Dado que existe un proyecto "Demo"), no pasos que
   describan la interfaz (clics, botones).

Plantilla de spec:
<pegar specs/_plantilla.md>

Contratos de dominio:
<pegar docs/adr/0002-contratos-de-dominio.md>

Ejemplo de spec ya aprobada:
<pegar specs/HU-12-registrar-historias-completadas.md>
```

## 3. Revisar el borrador antes de commitear (10')

No commitees nada que no puedas explicar. Revisá:

- [ ] **Cada criterio del issue** (`CA-XX.n`) está en la tabla de criterios y tiene al menos un escenario.
- [ ] Las **reglas de negocio** son verificables: se puede escribir un test que pase o falle.
- [ ] Hay **casos límite** reales (valores en el borde: 3 y 100 caracteres, cero, listas vacías, fechas iguales).
- [ ] Cada **condición de error** dice qué error se espera.
- [ ] Los **tipos y nombres** coinciden con el ADR 0002.
- [ ] La IA **no agregó** funcionalidades que el issue no pide.
- [ ] Resolviste o anotaste las **dudas** (si cambia un criterio, se actualiza también el issue).

## 4. Commitear y abrir el PR (5')

```bash
git add specs/HU-XX-*.md features/hu-xx-*.feature
git commit -m "docs(HU-XX): spec SDD y escenarios BDD"
git push -u origin feature/HU-XX-titulo-corto
gh pr create --draft --fill   # o desde la web
```

- Abrí el PR como **draft** desde este primer commit, con `Closes #N` en la descripción.
- Anotá el uso de IA en [`docs/ia/registro.md`](../ia/registro.md): fecha, qué pediste, qué usaste, cómo lo validaste.
- Sigue el ciclo TDD (sección 2 de `CONTRIBUTING.md`): `test(HU-XX): RED - …`, `feat(HU-XX): GREEN - …`, `refactor(HU-XX): …`.

## Si estás bloqueado

Si tu historia necesita un tipo de otra épica que todavía no está mergeado (ver la sección 6 del ADR 0002), avanzá
igual con la spec, los escenarios y los tests de **tu** dominio, y avisalo en la Daily.
