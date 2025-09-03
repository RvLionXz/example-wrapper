# Proyek Wrapper API Gemini dengan Go

## Ringkasan

Proyek ini adalah sebuah contoh lengkap tentang cara membuat server backend dengan bahasa Go yang berfungsi sebagai "wrapper" atau "proxy" yang aman untuk Google Gemini API. Tujuannya adalah untuk menyembunyikan API Key Gemini utama di sisi server dan memberikan API key yang berbeda kepada client.

Proyek ini terdiri dari tiga komponen utama:

1.  **Backend**: Sebuah server HTTP yang menerima permintaan dari client, menambahkan API Key Gemini, meneruskan permintaan ke Google, dan mengembalikan responsnya.
2.  **Omnic Library**: Sebuah library Go (`wrapper`) yang memudahkan interaksi dengan server backend kita.
3.  **Example Client**: Sebuah program Go sederhana yang menunjukkan cara menggunakan library `omnic`.

---

## Struktur Proyek

```
goclientside/
├── backend/              # Folder berisi kode server backend
│   └── main.go
├── omnic/                # Folder berisi kode library (wrapper)
│   └── omnic.go
├── example-client/       # Folder berisi contoh program client
│   └── main.go
├── go.mod                # File utama untuk manajemen modul Go
├── go.sum
└── README.md             # Dokumentasi ini
```

---

## Persyaratan

-   Bahasa Go (versi 1.21 atau lebih baru direkomendasikan).
-   API Key dari Google AI Studio (untuk Gemini).

---

## Cara Menjalankan Proyek

Ikuti langkah-langkah ini dari awal untuk menjalankan keseluruhan proyek.

### 1. Konfigurasi Modul Go

Proyek ini menggunakan Go Modules untuk mengelola dependensi dan paket lokal. Jika Anda memulai dari awal, Anda perlu menginisialisasi modul dan memberitahu Go di mana menemukan library `omnic` lokal kita.

*   Buka terminal di direktori root proyek (`goclientside`).
*   Inisialisasi modul:
    ```sh
    go mod init goclientside
    ```
*   Tambahkan referensi untuk library `omnic` lokal:
    ```sh
    go mod edit -replace=goclientside/omnic=./omnic
    ```

### 2. Menjalankan Backend Server

Server backend harus dijalankan terlebih dahulu.

1.  **Buka Terminal 1**.
2.  Pindah ke direktori `backend`:
    ```sh
    cd backend
    ```
3.  Set environment variable untuk API Key Gemini Anda. Ganti `YOUR_GEMINI_API_KEY` dengan key Anda yang sebenarnya.
    *   Di PowerShell (Windows):
        ```powershell
        $env:GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
        ```
    *   Di bash (Linux/macOS):
        ```sh
        export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"
        ```
4.  Jalankan server:
    ```sh
    go run main.go
    ```
5.  Biarkan terminal ini berjalan. Anda akan melihat log bahwa server aktif di port 8080.

### 3. Menjalankan Example Client

Setelah backend berjalan, buka terminal baru untuk menjalankan client.

1.  **Buka Terminal 2**.
2.  Pindah ke direktori `example-client`:
    ```sh
    cd example-client
    ```
3.  Jalankan program client:
    ```sh
    go run main.go
    ```

Jika semua berjalan lancar, Anda akan melihat prompt dikirim dan jawaban dari Gemini dicetak di terminal ini.

---

## Dokumentasi API Backend

Anda juga bisa berinteraksi dengan backend secara langsung menggunakan tool seperti Postman atau cURL.

-   **Endpoint**: `/api/generate`
-   **Method**: `POST`
-   **Headers**:
    -   `Content-Type`: `application/json`
    -   `X-Client-Api-Key`: Kunci yang valid (contoh: `supersecret-client-key-123`, didefinisikan di `backend/main.go`).
-   **Request Body** (JSON):
    ```json
    {
        "prompt": "jelaskan apa itu blockchain dalam satu kalimat"
    }
    ```
-   **Contoh cURL**:
    ```sh
    curl -X POST -H "Content-Type: application/json" -H "X-Client-Api-Key: supersecret-client-key-123" -d "{\"prompt\": \"jelaskan apa itu blockchain\"}" http://localhost:8080/api/generate
    ```

### Respons

-   **200 OK**: Jika sukses, Anda akan menerima respons JSON langsung dari Gemini API.
-   **401 Unauthorized**: Jika `X-Client-Api-Key` salah atau tidak ada.
-   **400 Bad Request**: Jika body JSON salah format atau field `prompt` kosong.
-   **500 Internal Server Error**: Jika backend gagal menghubungi Google atau ada masalah internal lainnya.

---

## Menggunakan Library `omnic`

Untuk menggunakan library ini di proyek Go lain (dalam modul yang sama), Anda bisa mengimpor dan menggunakannya seperti ini:

```go
package main

import (
    "fmt"
    "log"
    "goclientside/omnic" // Sesuaikan dengan nama modul Anda
)

func main() {
    // Konfigurasi client
    backendURL := "http://localhost:8080"
    clientAPIKey := "supersecret-client-key-123"

    // Buat client baru
    client := omnic.NewClient(backendURL, clientAPIKey)

    // Panggil methodnya
    text, err := client.GenerateContent("prompt Anda di sini")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }

    fmt.Println(text)
}
```
