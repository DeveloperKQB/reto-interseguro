# Decisiones técnicas

Este documento registra las decisiones tomadas durante el reto, las
alternativas consideradas y el motivo de cada elección.

---

## Interpretación del enunciado

### 1. Factorización QR en lugar de rotación
El enunciado menciona "rotación de la matriz" en la sección de arquitectura,
pero la funcionalidad requerida pide explícitamente la factorización QR, y la
segunda API habla de "las matrices devueltas" (plural), lo que corresponde a
Q y R. Elegí implementar QR.
Nota: la QR puede calcularse mediante rotaciones de Givens, por lo que ambos
conceptos están relacionados.

---

## API en Go (Fiber)

### 2. Fiber v2
Es la versión estable y más documentada de Fiber.

### 3. Householder implementado a mano
Más estable numéricamente que Gram-Schmidt, que pierde ortogonalidad por
errores de redondeo. No usé una librería (gonum) para demostrar dominio del
algoritmo; los tests verifican sus propiedades: A = Q·R, QᵀQ = I y R
triangular superior.

### 4. QR reducida con diagonal de R no negativa
Para A de m×n, con k = min(m, n): Q es m×k y R es k×n. Es la convención por
defecto de NumPy y envía menos datos a la API de Node. Forzar la diagonal de
R no negativa hace que el resultado sea único (la QR admite cambios de signo).

### 5. Limpieza de residuos numéricos
Los valores con |x| < 1e-12 se devuelven como 0 para evitar ruido de punto
flotante (por ejemplo 1e-17 o -0).

### 6. Límite de 100x100
La QR tiene costo O(m·n²); el límite evita que peticiones abusivas saturen
el servidor.

---

## API en Node.js (Express)

### 7. API genérica: recibe {"matrices": [...]}
No depende de que los datos vengan de una QR; acepta cualquier cantidad de
matrices. Así la API de Node queda desacoplada de la de Go.

### 8. Matriz diagonal con tolerancia 1e-10 y definición generalizada
Por el punto flotante, |x| <= 1e-10 cuenta como cero. Se aceptan matrices
rectangulares (a_ij = 0 para todo i != j), necesario porque la Q reducida
puede no ser cuadrada.

### 9. Estadísticas en una sola pasada; límite de body de 1 MB
Máximo, mínimo y suma se calculan recorriendo los datos una sola vez:
O(total de elementos). El límite por defecto de Express (100 KB) no alcanza
para Q y R de 100x100, por eso se amplió a 1 MB.

---

## Comunicación entre APIs

### 10. Go orquesta la llamada a Node
El cliente solo llama a POST /api/qr y recibe q, r y stats en una sola
respuesta, cumpliendo la arquitectura del enunciado (Go → Node). La URL de
Node es configurable con la variable de entorno STATS_API_URL, y la llamada
tiene un timeout de 5 segundos.

### 11. Si Node falla: 502 Bad Gateway
Alternativa considerada: devolver la QR sin estadísticas (degradación
parcial). Elegí 502 por claridad, ya que es el código HTTP correcto cuando
falla un servicio del que se depende. El detalle técnico se registra en el
log y al cliente se le devuelve un mensaje genérico, para no exponer
información interna.

### 12. Inyección de dependencias vía interfaz StatsProvider
El handler recibe el cliente de Node a través de una interfaz. Esto permite
probarlo con un fake, sin levantar la API de Node.

---

## Contenedores (Docker)

### 13. Docker multi-stage, distroless y non-root
Una etapa compila y otra, mínima, solo ejecuta. La imagen de Go usa
distroless (sin shell ni gestor de paquetes), lo que reduce la superficie de
ataque. Ambos contenedores corren con usuarios sin privilegios. Node usa la
versión LTS (24) en el contenedor.

### 14. Node como servicio interno
En docker-compose solo Go expone su puerto; Node únicamente es accesible
desde la red interna de Docker. Go espera a que Node esté sano
(healthcheck) antes de arrancar.

---

## Seguridad

### 15. Cabecera X-Powered-By deshabilitada en Express
Evita revelar la tecnología del servidor a posibles atacantes.