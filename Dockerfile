# Usamos a imagem completa do Go como nossa base final para o ambiente de desenvolvimento.
# Isso garante que o compilador 'go' e outras ferramentas estejam disponíveis para o Air.
FROM golang:1.24-alpine

# Instala o Git, que pode ser necessário para algumas dependências.
RUN apk add --no-cache git

# Define o diretório de trabalho dentro do container.
WORKDIR /app

# Instala a ferramenta de hot-reloading Air usando o novo caminho oficial.
RUN go install github.com/air-verse/air@latest

# Copia os arquivos de gerenciamento de dependências.
COPY go.mod go.sum ./

# Baixa as dependências para o cache do container.
RUN go mod download

# Copia todo o código-fonte para o diretório de trabalho.
# O volume no docker-compose.yml irá sobrepor isso em tempo de execução,
# mas esta cópia é útil para o build inicial e para ter um estado base.
COPY . .

# Expõe a porta que a nossa aplicação usa.
EXPOSE 8080

# Comando final que inicia o container.
# Nós rodamos o Air, que por sua vez irá compilar e rodar nosso app,
# reiniciando-o automaticamente a cada alteração de arquivo.
CMD ["air", "-c", ".air.toml"]