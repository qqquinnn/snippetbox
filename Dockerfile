# --- STAGE 1: Build Go binary ---
FROM golang:1.26.2-alpine AS builder

# Set working directory inside container.
WORKDIR /app

# Copy dependency files.
# Unless go.mod/go.sum change, modules are not re-downloaded.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code.
COPY . .

# Build the application.
# -o web: name the output binary 'web'
# ./cmd/web: path to the main package
# -ldflags="-s -w": Strips debug info to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o web ./cmd/web

# --- STAGE 2: Create lightweight runtime image ---
FROM alpine:latest

# Create non-privileged user to run application.
RUN adduser -D snippetbox
USER snippetbox

WORKDIR /home/snippetbox

# Copy only compiled binary from builder stage.
COPY --from=builder /app/web .

# Copy TLS certs.
# In prod, these are likely mounted as secrets or volumes.
# For local dev, we copy the folder.
COPY --from=builder /app/tls ./tls

# Application listens on :4000 by default.
EXPOSE 4000

# Run the binary.
CMD [ "./web" ]