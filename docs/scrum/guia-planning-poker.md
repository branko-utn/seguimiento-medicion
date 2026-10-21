# Guía de Planning Poker

Para leer **antes** de la reunión de arranque del jueves 22/10, así no perdemos tiempo explicando cómo se juega.

## Para qué sirve

Estimar el tamaño de una historia **comparándola con otras**, no adivinando horas. La unidad son los
**Story Points**: miden cuánto trabajo tiene, qué tan compleja es y cuánta incertidumbre hay.

## La baraja

Sucesión de Fibonacci: **1, 2, 3, 5, 8, 13, 21**, más dos cartas especiales:

- **?**: no entiendo la historia, necesito que me la expliquen.
- **☕**: es demasiado grande, hay que partirla en historias más chicas.

Los números se separan a propósito: cuanto más grande la historia, menos precisión tiene sentido.

## La historia de referencia

Antes de estimar, elegimos **una historia chica y conocida** y le asignamos **3 puntos**.
Todas las demás se comparan con esa: "¿es más grande o más chica que la de referencia?".

## Cómo se juega cada historia

1. Alguien lee la historia y sus criterios de aceptación.
2. Se hacen las preguntas que hagan falta (es la parte que más valor tiene).
3. Cada uno elige su carta **sin decirla**.
4. Se muestran todas **al mismo tiempo** (en remoto: cada uno escribe su número en el chat a la cuenta de tres).
5. Si todos coinciden, se anota y se pasa a la siguiente.
6. Si hay diferencias, explican primero **el que puso el número más alto y el que puso el más bajo**: no para convencer, sino para contar qué vieron que el resto no vio.
7. Se vuelve a votar. Como máximo **3 rondas**.
8. Si después de 3 rondas no hay acuerdo, se toma el valor más alto o se parte la historia.

El resultado se carga en el campo **Estimate** del issue en el tablero.

> Además de usarlo, lo tenemos que **programar**: es el requerimiento 4 de la consigna (épica E4).
> Jugarlo antes nos ayuda a entender el dominio.
