FROM golang:1.22-alpine AS build
WORKDIR /src
COPY services/worker/go.mod services/worker/go.sum* ./
RUN go mod download
COPY services/worker/ ./
RUN go build -o /out/worker ./cmd/worker

FROM alpine:3.20
RUN addgroup -S tipdrop && adduser -S tipdrop -G tipdrop
USER tipdrop
WORKDIR /app
COPY --from=build /out/worker /app/worker
CMD ["/app/worker"]
