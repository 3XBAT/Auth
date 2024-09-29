FROM golang:1.23-alpine AS Builder
# Это файл для аутфывадь
WORKDIR /app

COPY . .

RUN go mod download 

RUN go build -o /bin/application cmd/auth/main.go

FROM alpine:latest AS Runner

COPY --from=builder /bin/application ./

COPY  cmd/config/local.yaml /config.yaml

CMD [ "/application" ]