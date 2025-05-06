# Etapa de build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependência e baixar módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo o código do projeto
COPY . .

# Compilar a aplicação
RUN go build -o bloGo ./cmd/web

# Etapa final: imagem leve para produção
FROM alpine:latest

WORKDIR /root/

# Copiar o binário da aplicação
COPY --from=builder /app/bloGo .

# Copiar os diretórios necessários
COPY --from=builder /app/posts ./posts
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/.env ./.env

# Expor a porta usada pela aplicação
EXPOSE 8080

# Comando padrão
CMD ["./bloGo"]
