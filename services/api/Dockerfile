FROM golang:1.22-alpine AS build
WORKDIR /src
COPY services/api/go.mod services/api/go.sum* ./
RUN go mod download
COPY services/api/ ./
RUN go build -o /out/api ./cmd/api

FROM alpine:3.20
RUN addgroup -S tipdrop && adduser -S tipdrop -G tipdrop
USER tipdrop
WORKDIR /app
COPY --from=build /out/api /app/api
EXPOSE 8080
CMD ["/app/api"]
