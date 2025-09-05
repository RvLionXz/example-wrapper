# Proyek Wrapper API Gemini (Format OpenAI)

## Ringkasan

Versi terbaru dari proyek ini berfungsi sebagai **server proxy yang meniru antarmuka (interface) dari OpenAI API**, namun di belakang layar tetap menggunakan Google Gemini sebagai model pemrosesnya. Tujuannya adalah untuk menyediakan sebuah "gerbang tunggal" dengan format API yang standar dan populer (OpenAI) untuk berinteraksi dengan berbagai model AI.

Proyek ini terdiri dari tiga komponen utama:

1.  **Backend**: Server HTTP yang menerima request dalam format OpenAI, "menerjemahkannya" untuk Gemini, mengirim request ke Google, lalu "menerjemahkan" kembali responsnya ke format OpenAI sebelum dikirim ke client.
2.  **Omnic Library**: Library Go (`wrapper`) yang akan kita sesuaikan untuk berinteraksi dengan backend baru ini.
3.  **Example Client**: Program Go sederhana untuk menunjukkan cara menggunakan library `omnic`.

---

## Struktur Proyek

Struktur folder tidak berubah:
```
goclientside/
├── backend/              # Folder berisi kode server backend
│   └── main.go
├── omnic/                # Folder berisi kode library (wrapper)
│   └── omnic.go
├── example-client/       # Folder berisi contoh program client
│   └── main.go
├── go.mod                # File utama untuk manajemen modul Go
└── README.md             # Dokumentasi ini
```

---

## Cara Menjalankan Proyek

### 1. Konfigurasi Modul Go

(Langkah ini tidak perlu diulangi jika sudah dilakukan sebelumnya)
Pastikan file `go.mod` di root proyek Anda berisi baris berikut:
```mod
module goclientside

replace goclientside/omnic => ./omnic
```

### 2. Menjalankan Backend Server

1.  **Buka Terminal 1**.
2.  Pindah ke direktori `backend`:
    `cd backend`
3.  Set environment variable untuk API Key Gemini Anda:
    *   Di PowerShell: `$env:GEMINI_API_KEY="YOUR_GEMINI_API_KEY"`
    *   Di bash: `export GEMINI_API_KEY="YOUR_GEMINI_API_KEY"`
4.  Jalankan server:
    `go run main.go`
5.  Server sekarang berjalan dan siap menerima request di endpoint `/v1/chat/completions`.

---

## Dokumentasi API Backend (Format OpenAI)

Untuk berinteraksi langsung dengan backend (misalnya via Postman atau cURL).

-   **Endpoint**: `/v1/chat/completions`
-   **Method**: `POST`
-   **Headers**:
    -   `Content-Type`: `application/json`
    -   *(Otentikasi `X-Client-Api-Key` sudah dihapus di versi ini)*

-   **Request Body** (JSON):
    Struktur body sekarang harus mengikuti format OpenAI. Field yang paling penting adalah `model` dan `messages`.
    ```json
    {
      "model": "gemini-1.5-flash-latest",
      "messages": [
        {
          "role": "user",
          "content": "Tulis sebuah lagu tentang bahasa pemrograman Go."
        }
      ]
    }
    ```

-   **Success Response** (JSON):
    Jika sukses, backend akan memberikan respons dalam format OpenAI.
    ```json
    {
      "id": "chatcmpl-1716823888",
      "object": "chat.completion",
      "created": 1716823888,
      "model": "gemini-1.5-flash-latest",
      "choices": [
        {
          "index": 0,
          "message": {
            "role": "assistant",
            "content": "(Verse 1)\nKetik `package main`, dunia menyapa..."
          },
          "finish_reason": "stop"
        }
      ]
    }
    ```

-   **Error Responses**:
    -   `400 Bad Request`: Jika body JSON yang dikirim tidak sesuai format.
    -   `500 Internal Server Error`: Jika terjadi kesalahan saat backend berkomunikasi dengan Google Gemini.