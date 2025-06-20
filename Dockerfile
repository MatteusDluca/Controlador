# --- FASE 1: BUILDER ---
# Usa a imagem oficial e mais recente do Go, que já tem tudo que precisamos.
FROM golang:1.23-alpine AS builder

# Instala o Git, que ainda pode ser útil.
RUN apk add --no-cache git

WORKDIR /app

# Copia todo o nosso contexto de projeto, incluindo a pasta "vendor".
COPY . .

# Compila o nosso aplicativo.
# A flag "-mod=vendor" é a ordem crucial: "Use apenas os pacotes da pasta vendor".
# A build é agora 100% offline e autossuficiente.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -o /usr/local/bin/controlador-server ./cmd/server

# --- FASE 2: FINAL ---
# Usa a imagem mínima do Alpine para um resultado final pequeno e seguro.
FROM alpine:latest

# Instala apenas os certificados, que são essenciais.
RUN apk add --no-cache ca-certificates

# Copia APENAS o binário compilado da fase de 'builder'.
COPY --from=builder /usr/local/bin/controlador-server /usr/local/bin/controlador-server

# Expõe a porta.
EXPOSE 8080

# O comando para iniciar nosso servidor.
CMD ["/usr/local/bin/controlador-server"]