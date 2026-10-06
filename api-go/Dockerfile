# Etapa 1: compilar
FROM golang:1.27-alpine AS build
WORKDIR /src

# Primero las dependencias, para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Binario estático (sin dependencias de C)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/api-go .

# Etapa 2: ejecutar (imagen mínima, sin shell, usuario no root)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /app/api-go /api-go
EXPOSE 8080
ENTRYPOINT ["/api-go"]