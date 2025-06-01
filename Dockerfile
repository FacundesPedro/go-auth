# Estágio de compilação
FROM golang:1.24.3-alpine AS builder

# Configura variáveis de ambiente para compilação
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

# Instala dependências de compilação
RUN apk add --no-cache git ca-certificates

# Configura o diretório de trabalho
WORKDIR /app

# Copia os arquivos de dependência primeiro para melhor cache
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código-fonte
COPY . .

# Compila o binário estático
RUN go build -ldflags='-w -s' -o /app/go-auth

# Estágio de execução
FROM alpine:3.19

# Diretório de trabalho para execução
WORKDIR /dist

# Instala certificados SSL
RUN apk add --no-cache ca-certificates

# Copia o binário compilado
COPY --from=builder /app/go-auth /usr/local/bin/go-auth
# Copia pasta public com arquivos estaticos para dentro do alpine
COPY --from=builder /app/public /dist/public

# Diretório de trabalho para execução
# WORKDIR /app

# Expõe a porta da aplicação
EXPOSE 5000

# Define variáveis de ambiente necessárias
#ENV OAUTH_CLIENT_ID="" \
#    OAUTH_CLIENT_SECRET=""

# Comando de inicialização
ENTRYPOINT ["/usr/local/bin/go-auth"]
