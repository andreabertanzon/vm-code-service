FROM golang:1.21.6 as builder
WORKDIR /app
COPY go.mod go.sum env.env ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the built binary and the env.env file from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/env.env . 
CMD ["./main"]
EXPOSE 8080
