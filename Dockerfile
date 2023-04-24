# Build stage
FROM golang:1.20.3-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage 
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./migration
COPY start.sh .

EXPOSE 8181
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]