# Build stage
FROM golang:1.19 AS build

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pvc-admission-controller

# Final stage
FROM ubuntu

# Copy the Go binary from the build stage and grant execute permissions
COPY --from=build /src/pvc-admission-controller /app/pvc-admission-controller
RUN chmod +x /app/pvc-admission-controller

EXPOSE 443
CMD ["/app/pvc-admission-controller"]
