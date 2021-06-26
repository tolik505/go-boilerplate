FROM alpine:latest

WORKDIR /app
COPY bin/ .
COPY migrations/ ./migrations

CMD ["/app/httpserver"]

