FROM golang:alpine as builder

WORKDIR /app
COPY . .
RUN go build -o niuwa -tags=jsoniter .

FROM alpine:latest
WORKDIR /app
COPY ./config.toml.example ./config.toml
COPY --from=builder /app/niuwa ./

EXPOSE 8080
ENTRYPOINT ["./niuwa"]
