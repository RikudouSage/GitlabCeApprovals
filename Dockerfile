FROM golang:1.23 as build

ENV CGO_ENABLED=0

WORKDIR /app
COPY . .
RUN go build -o app .

FROM ubuntu:24.04
COPY --from=build /app/app /app
ENTRYPOINT ["/app"]
