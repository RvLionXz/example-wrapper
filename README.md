# Proyek Wrapper API Gemini (Format Gemini)

## Ringkasan

Versi terbaru dari proyek ini berfungsi sebagai **server proxy yang meniru antarmuka (interface) dari Gemini API dan mendukung streaming**. Di belakang layar, ia tetap menggunakan Google Gemini sebagai model pemrosesnya.

---

## Cara Menjalankan Proyek

(Cara menjalankan backend dan client tidak berubah, silakan lihat versi README sebelumnya jika perlu)

---

## Dokumentasi API Backend (Format Gemini)

Untuk berinteraksi langsung dengan backend (misalnya via Postman atau cURL).

-   **Endpoint**: `/v1/chat/completions`
-   **Method**: `POST`
-   **Headers**:
    -   `Content-Type`: `application/json`

### Request Body (JSON)

Struktur body mengikuti format Gemini. Untuk mengaktifkan streaming, tambahkan field `"stream": true`.

**Contoh Request Streaming:**
```json
{
  "model": "gemini-1.5-flash-latest",
  "messages": [
    {
      "role": "user",
      "content": "Tulis sebuah puisi singkat tentang hujan."
    }
  ],
  "stream": true
}
```

### Respons

Respons dari server akan berbeda tergantung pada nilai field `stream`.

#### 1. Respons Non-Streaming (`stream: false` atau tidak ada)

Jika `stream` tidak di-set ke `true`, Anda akan menerima **satu objek JSON besar** di akhir, setelah AI selesai berpikir.

**Contoh:**
```json
{
  "id": "chatcmpl-...
  "object": "chat.completion",
  "created": 1716823888,
  "model": "gemini-1.5-flash-latest",
  "choices": [
    {
      "index": 0,
      "message": {
        "role": "assistant",
        "content": "Tetes air jatuh membasahi bumi,\nMembawa lagu rindu dalam sunyi."
      },
      "finish_reason": "stop"
    }
  ]
}
```

#### 2. Respons Streaming (`stream: true`)

Jika `stream` di-set ke `true`, Anda akan menerima **aliran data (stream)** yang berkelanjutan. Koneksi akan tetap terbuka dan server akan mengirim potongan-potongan data setiap kali AI menghasilkan teks baru.

Header respons akan berisi `Content-Type: text/event-stream`.

Body respons akan terlihat seperti ini:

```
data: {"choices":[{"delta":{"content":"Tetes"}}]}

data: {"choices":[{"delta":{"content":" air"}}]}

data: {"choices":[{"delta":{"content":" jatuh"}}]}

data: {"choices":[{"delta":{"content":" membasahi"}}]}

data: {"choices":[{"delta":{"content":" bumi"}}]}

... 

data: [DONE]
```

**Penjelasan Potongan Data (Chunk):**
-   Setiap baris yang diawali dengan `data: ` adalah satu potongan JSON yang terpisah.
-   Client Anda perlu membaca stream ini baris per baris, mengambil JSON setelah `data: `, mem-parsing-nya, dan mengambil teks dari dalam `delta.content` untuk ditampilkan.
-   Aliran akan diakhiri dengan pesan `data: [DONE]` (tergantung implementasi server, namun ini adalah pola umum).