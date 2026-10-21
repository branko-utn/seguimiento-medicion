# Definition of Ready y Definition of Done

## Definition of Ready: una historia puede entrar a un sprint si

- [ ] Está escrita en formato **Como / Quiero / Para**.
- [ ] Tiene **criterios de aceptación** numerados (`CA-XX.1`, `CA-XX.2`, …) y validados con los Product Architects.
- [ ] Está **estimada** en Story Points con Planning Poker.
- [ ] Sus **dependencias** con otras historias están identificadas.
- [ ] Tiene **responsable** asignado.

## Definition of Done: una historia está terminada si

- [ ] La **especificación SDD** está completa en `specs/` (objetivo, entradas, salidas, reglas, restricciones, casos límite, errores y criterios).
- [ ] Los **escenarios BDD** cubren caso normal, alternativo, límite y error, están en `features/` y pasan.
- [ ] Los **tests unitarios** se escribieron con TDD: el historial muestra commits `RED → GREEN → REFACTOR`.
- [ ] `gofmt`, `go vet` y `go test ./...` pasan en el CI.
- [ ] La cobertura de `internal/domain` e `internal/metrics` es de al menos **80 %** (meta interna del equipo).
- [ ] Los errores se manejan y se muestran al usuario con un mensaje claro.
- [ ] El **Pull Request** fue revisado y aprobado por otro integrante, con la tabla de trazabilidad completa.
- [ ] Se probó en la aplicación y cumple los criterios de aceptación.
- [ ] Si se usó IA, quedó registrado en `docs/ia/registro.md`.
- [ ] La tarjeta está en **Done** en el tablero.
