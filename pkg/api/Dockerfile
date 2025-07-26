# build stage (optional: can be done on host)
FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# final minimal image
FROM alpine:3.20

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
# ENV TODO_PASSWORD= # задайте снаружи, если нужно

WORKDIR /app
COPY --from=builder /app/app /app/app
COPY web /app/web

EXPOSE 7540

# БД положим наружу
VOLUME ["/data"]

CMD ["/app/app"]
