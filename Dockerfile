# Build aşaması
FROM golang:1.21.8-alpine AS builder

# Çalışma dizinini ayarlayın
WORKDIR /app

# Go mod dosyalarını kopyalayın
COPY go.mod go.sum ./

# Modülleri indirin
RUN go mod download

# Uygulama kaynak kodlarını kopyalayın
COPY . .

# Go uygulamasını derleyin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o differ-template-engine .

# Run aşaması
FROM alpine:latest

# Sertifikaları ekleyin
RUN apk --no-cache add ca-certificates

# Çalışma dizinini ayarlayın
WORKDIR /root/

# Derlenmiş uygulamayı kopyalayın
COPY --from=builder /app/differ-template-engine .

# Uygulamanızı çalıştırın
CMD ["./differ-template-engine"]
