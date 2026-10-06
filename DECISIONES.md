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
tiene un timeout de 5 segundos para no quedar bloqueada si Node no responde.

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
ataque. Todos los contenedores corren con usuarios sin privilegios. Node usa
la versión LTS (24) en el contenedor.

### 14. Node como servicio interno
En docker-compose solo Go y el frontend exponen puertos; Node únicamente es
accesible desde la red interna de Docker. Go espera a que Node esté sano
(healthcheck) antes de arrancar.

---

## Pruebas

### 15. Tests en Go: table-driven y basados en propiedades
Además de un caso con valores conocidos, se verifican las propiedades que
toda QR debe cumplir (A = Q·R, QᵀQ = I, R triangular superior con diagonal
no negativa) sobre matrices cuadradas, altas, anchas, diagonales, de una
fila y con columnas nulas.

### 16. Tests en Node: node:test + supertest
El test runner nativo de Node evita dependencias extra. supertest prueba los
endpoints en memoria gracias a la separación entre createApp() y listen().
Se cubren estadísticas, validaciones, endpoints y casos de seguridad (sin
token, firma inválida, token expirado). supertest es dependencia de
desarrollo y no entra en la imagen de Docker.

---

## Frontend

### 17. HTML, CSS y JavaScript sin framework
Una pantalla de login, un formulario y un resultado no justifican React o
Angular: sin build, sin dependencias.

### 18. nginx como servidor y proxy inverso
nginx sirve el frontend y redirige /api/ al contenedor de Go. El navegador
ve un solo origen, por lo que no hace falta configurar CORS. Se usa la
imagen nginx-unprivileged (non-root).

### 19. Token en memoria y CSP estricta
El token no se guarda en localStorage, para no dejar un token persistido
expuesto ante un XSS; al recargar se pide login. La Content Security Policy
solo permite scripts y estilos propios (JS y CSS en archivos separados) y
los resultados se insertan con textContent, nunca como HTML.

### 20. Redondeo solo visual
La interfaz muestra 4 decimales; la API conserva la precisión completa,
visible al pasar el cursor sobre cada número.

---

## Seguridad

### 21. Ocultar la tecnología del servidor
Se deshabilitó la cabecera X-Powered-By en Express y la versión de nginx
(server_tokens off), para no revelar información útil a un atacante.

### 22. JWT HS256 con secreto compartido; Go emite el token
POST /api/auth/login valida credenciales de demostración definidas por
variables de entorno y emite un token válido por 1 hora. En producción esta
función la cumpliría un proveedor de identidad (Azure AD, Auth0, Keycloak).
Solo se acepta HS256, para evitar ataques de confusión de algoritmo
("alg": "none"). La contraseña se compara en tiempo constante para evitar
ataques de temporización.

### 23. Go reenvía el token del usuario a Node (confianza cero)
Node valida el token aunque sea un servicio interno: no se confía en una
petición solo por provenir de la red interna. Alternativa considerada: un
token de servicio propio entre APIs.

### 24. Sin secretos por defecto
Si JWT_SECRET falta o tiene menos de 32 caracteres, las APIs se niegan a
arrancar: es preferible fallar a funcionar con un secreto inseguro. Los
secretos viven en .env (excluido de Git) y .env.example sirve de plantilla.

---

## Despliegue

### 25. Despliegue en la nube: propuesta con Azure Container Apps
No se incluyó en esta entrega por el plazo disponible. La propuesta es Azure
Container Apps porque permite reproducir el diseño de docker-compose: API Go
y frontend con ingreso externo, y API Node con ingreso solo interno. Los
secretos se gestionarían como secretos de la plataforma.
