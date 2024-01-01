# Stage 1: Build binary
FROM golang:alpine as builder
WORKDIR /go/src/be-service-auth
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o be-service-auth app/main.go

# Stage 2: Create a minimal image to run the application
FROM alpine:latest
EXPOSE 8882

# Install Git
RUN apk --no-cache add git

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app/
COPY --from=builder /go/src/be-service-auth/be-service-auth .
COPY --from=builder /go/src/be-service-auth/db/migration ./db/migration
COPY openapi-submodule/openapi-auth.yaml /app/openapi-auth.yaml
CMD ["./be-service-auth"]
