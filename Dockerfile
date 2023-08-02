FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . .
RUN go build

FROM alpine AS release

WORKDIR /app

COPY --from=build /app/meteor-server .
COPY config.json .

ENTRYPOINT [ "/app/meteor-server" ]