# Builder
FROM golang:1.22.2 as build

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o cloudrun ./cmd/server


# Runner
FROM scratch

ENV HOST=0.0.0.0
ENV PORT=8080

WORKDIR /app

COPY --from=build /app/cloudrun .

ENTRYPOINT ["./cloudrun"]
