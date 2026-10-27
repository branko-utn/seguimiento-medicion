# ADR 0004 · Identidad sin autenticación

- **Estado:** aceptada en la planning del Sprint 2
- **Responsable:** Branko Almeira

## Contexto

El Planning Poker (HU-16 a HU-18) necesita saber **quién vota**, y en el futuro otras pantallas podrían querer saber
quién registra horas o defectos. La aplicación no tiene registro de usuarios ni login: quedaron fuera de alcance en el
ADR 0001 porque la consigna no los pide y nos llevarían días.

## Decisión

**Opción A · Identidad elegida.** Cuando una pantalla necesita saber quién actúa, la persona **elige su nombre de la lista
de integrantes del proyecto** (HU-03), por ejemplo "Votás como: Ana Pérez". Ese `IntegranteID` es el que se registra.

- Cada uno puede usar la aplicación desde su navegador (la app corre en una máquina de la red) o por turnos en una sola pantalla.
- Las reglas de negocio se aplican igual: solo pueden actuar integrantes del proyecto (RN-16.1).
- Los roles (`Rol` de HU-03) pueden habilitar o no acciones en la interfaz (por ejemplo, revelar votos), pero es una
  regla de interfaz, no de seguridad.

## Alternativas descartadas

- **Login con usuario y contraseña:** sesiones, roles y permisos son una épica entera; no entra antes del 2/11.
- **Un facilitador que carga todos los votos:** queda como plan B si la red falla en la demo.

## Consecuencias

- **Limitación aceptada:** cualquiera puede elegir el nombre de otra persona. Es una herramienta interna de un equipo
  que confía entre sí.
- Si en el futuro se agrega login, solo cambia cómo se obtiene el `IntegranteID`: el dominio y los servicios no cambian.
