FROM golang:1.21.11-alpine3.20 AS build-stage

WORKDIR /app

COPY ./ ./
RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o /api

FROM gcr.io/distroless/base-debian12:latest AS run-stage

WORKDIR /

COPY --from=build-stage /api /api

USER nonroot:nonroot
ENTRYPOINT ["/api"]