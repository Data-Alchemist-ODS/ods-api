# Langkah pertama, pilih base image yang sesuai untuk Go
FROM golang:1.20.5-alpine

# Set working directory di dalam wadah
WORKDIR /app

# Copy file go.mod dan go.sum terlebih dahulu agar dependensi dapat di-cache
COPY go.mod ./
COPY go.sum ./

# Download dependensi menggunakan go mod
RUN go mod download

# Copy seluruh file kode Go dari direktori lokal ke direktori kerja di wadah
COPY *.go ./

# Build aplikasi Go dan simpan binary dengan nama "main" di root direktori ("/")
RUN go build -o /main

# Tetapkan variabel lingkungan yang diperlukan (sesuaikan dengan .env Anda)
ENV DATABASE_URL=mongodb+srv://AdrianBadjideh:vQPm8EgUsKlIeeT2@ods.rycue7a.mongodb.net/?retryWrites=true&w=majority
ENV API_KEY=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
ENV DATABASE=ODS
ENV HOST=localhost
ENV PORT=8000

# Tandai port 8080 yang akan digunakan oleh aplikasi
EXPOSE 8080

# Perintah CMD akan dijalankan ketika wadah berjalan
CMD ["/main"]
