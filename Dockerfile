# ETAPA 1: Construcción (Builder)
# Usamos una imagen oficial de Go con Alpine (ligera)
FROM golang:1.25-alpine AS builder

# Establecemos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos de dependencias primero (para aprovechar el caché de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos la aplicación
# -o main: nombre del binario de salida
# ./cmd/api: ruta donde está tu main.go
RUN go build -o main ./cmd/api

# ETAPA 2: Ejecución (Runner)
# Usamos una imagen vacía de Alpine Linux para producción
FROM alpine:latest

WORKDIR /root/

# Copiamos solo el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Copiamos las migraciones y seeds
COPY --from=builder /app/internal/database/seeds ./internal/database/seeds

# Exponemos el puerto interno (debe coincidir con el de tu .env)
EXPOSE 8080

# Comando para ejecutar la app
CMD ["./main"]