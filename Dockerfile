# Usa una imagen base con Go instalado
FROM golang:latest

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /go/src/app

# Copia el código fuente al directorio de trabajo del contenedor
COPY . .

# Descarga las dependencias del módulo Go
RUN go mod download

# Compila la aplicación
RUN go build -o worker .

# Ejecuta la aplicación
CMD ["./worker"]
