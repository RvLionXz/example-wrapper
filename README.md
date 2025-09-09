# Go Gemini API Proxy & Client Library

## Ringkasan

Proyek ini adalah sebuah sistem client-server lengkap yang ditulis dalam Go. Tujuannya adalah untuk menyediakan sebuah backend proxy yang aman dan cerdas untuk Google Gemini API, serta sebuah library (wrapper) yang mudah digunakan untuk berinteraksi dengannya.

## Fitur Utama

-   **Backend Proxy Aman**: Menyembunyikan API Key Google Gemini utama di sisi server.
-   **Dukungan Hybrid (Streaming & Non-Streaming)**: Satu endpoint cerdas yang bisa memberikan respons streaming atau respons tunggal (non-streaming) tergantung dari parameter request.
-   **Parameter Lanjutan**: Meneruskan parameter tambahan seperti `temperature` ke Google Gemini API.
-   **Client Library (`omnic`)**: Sebuah wrapper Go yang menyederhanakan interaksi dengan backend, menyediakan fungsionalitas siap pakai untuk aplikasi lain.

---

## Struktur Proyek

```
goclientside/
├── backend/              # Kode server backend
│   └── main.go
├── omnic/                # Kode library client (wrapper)
│   └── omnic.go
├── example-client/       # Contoh program yang menggunakan library omnic
│   └── main.go
├── go.mod                # File manajemen modul Go
└── README.md             # Dokumentasi ini
```

---

## Cara Menjalankan Proyek

### 1. Konfigurasi Modul Go

(Langkah ini tidak perlu diulangi jika sudah dilakukan sebelumnya). Pastikan file `go.mod` di root proyek Anda berisi baris `replace` untuk library lokal kita:
```mod
module goclientside

replace goclientside/omnic => ./omnic
```

### 2. Menjalankan Backend Server

1.  **Buka Terminal 1**.
2.  Pindah ke direktori `backend`: `cd backend`
3.  Set environment variable untuk API Key Gemini Anda:
    *   Di PowerShell: `$env:GEMINI_API_KEY="YOUR_GEMINI_API_KEY"`
    *   Di bash: `export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"`
4.  Jalankan server: `go run main.go`
5.  Biarkan terminal ini berjalan.

### 3. Menjalankan Example Client

1.  **Buka Terminal 2**.
2.  Pindah ke direktori `example-client`: `cd example-client`
3.  Jalankan client: `go run main.go`

---

## Dokumentasi API Backend

-   **Endpoint**: `POST /v1/chat/completions`
-   **Request Body (JSON)**:
    ```json
    {
      "model": "gemini-1.5-flash-latest",
      "messages": [
        {
          "role": "user",
          "content": "Tulis sebuah cerita pendek."
        }
      ],
      "stream": true,
      "temperature": 0.8
    }
    ```
-   **Response Body**: Responsnya adalah **respons mentah dari Google Gemini API**. Bentuknya tergantung pada nilai `stream`:
    -   Jika `"stream": false`, responsnya adalah satu objek JSON besar: `{"candidates":[...]}`.
    -   Jika `"stream": true`, responsnya adalah aliran data (stream) Server-Sent Events: `data: {"candidates":[...]}`.

---

## Penjelasan Wrapper Library (`omnic`)

Library `omnic` bertujuan untuk menyembunyikan kerumitan komunikasi HTTP dan parsing stream dari pengguna akhir.

### Cara Menggunakan

**1. Impor Library**
```go
import "goclientside/omnic" // Sesuaikan dengan nama modul Anda
```

**2. Buat Client Baru**
`Client` adalah objek utama yang menyimpan konfigurasi (seperti URL backend).
```go
client := omnic.NewClient("http://localhost:8080")
```

**3. Buat Request**
Buat sebuah `struct` `APIRequest` untuk mendefinisikan apa yang Anda inginkan.
```go
request := omnic.APIRequest{
    Model:  "gemini-1.5-flash-latest",
    Stream: true,
    Messages: []omnic.Message{
        {Role: "user", Content: "Halo dunia!"},
    },
}
```

**4. Panggil Method `ChatCompletionCreate`**
Ini adalah method utama yang cerdas. Ia selalu mengembalikan sebuah *channel*, tidak peduli apakah Anda meminta streaming atau tidak.
```go
responseChan, err := client.ChatCompletionCreate(request)
if err != nil {
    log.Fatal(err)
}
```

**5. Baca Hasil dari Channel**
Gunakan `for range` loop untuk membaca hasil. Loop ini secara ajaib bekerja untuk kedua kasus:
-   Jika streaming, ia akan berjalan berkali-kali, mencetak setiap potongan teks.
-   Jika non-streaming, ia hanya akan berjalan satu kali, mencetak seluruh teks jawaban.

```go
for textChunk := range responseChan {
    fmt.Print(textChunk)
}
```

Dengan menggunakan wrapper ini, program utama Anda (`example-client`) menjadi sangat bersih dan tidak perlu tahu tentang `HTTP`, `JSON`, atau kerumitan `streaming`. Itulah kekuatan dari sebuah wrapper yang baik.