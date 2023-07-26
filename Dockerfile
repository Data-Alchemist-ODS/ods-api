# Langkah pertama, pilih base image yang sesuai untuk Go
FROM golang:1.20.5-alpine

# Set working directory di dalam wadahnya
WORKDIR /app

# Copy file go.mod dan go.sum terlebih dahulu agar dependensi dapat di-cache
COPY go.mod go.sum ./

# Download dependensi menggunakan go mod
RUN go mod download

# Copy seluruh file kode Go dari direktori lokal ke direktori kerja di wadah
COPY . .

# Build aplikasi Go dan simpan binary dengan nama "app" di root direktori ("/")
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Stage 2: Create a minimal container to run the application
FROM alpine:latest

# Install CA certificates to support HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=0 /app/app .

# Tetapkan variabel lingkungan yang diperlukan (sesuaikan dengan .env Anda)
ENV DATABASE_URL=mongodb+srv://AdrianBadjideh:vQPm8EgUsKlIeeT2@ods.rycue7a.mongodb.net/?retryWrites=true&w=majority
ENV API_KEY=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
ENV DATABASE=ODS
ENV HOST=localhost
ENV PORT=8000

# Tandai port 8000 yang akan digunakan oleh aplikasi
EXPOSE 8000

# Run the Go Fiber application
CMD ["./app"]
