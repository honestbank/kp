FROM golang:1.18 as builder
WORKDIR /app
COPY . .

RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go
RUN ldd main

FROM gcr.io/distroless/base

WORKDIR /app
# Copy the Pre-built binary file from the previous stage
COPY --from=builder --chown=nonroot:nonroot /app/main .
EXPOSE 8888

USER nonroot

# Command to run the executable
CMD ["./main", "server"]
