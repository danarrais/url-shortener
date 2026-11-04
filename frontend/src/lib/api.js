// Semua komunikasi HTTP ke backend ada di sini, supaya komponen tidak perlu
// tahu soal fetch, status code, atau bentuk pesan error.

const API_URL = import.meta.env?.VITE_API_URL ?? 'http://localhost:8000'

/** Membaca body error menjadi pesan siap tampil. Backend membalas {"error": "..."}. */
async function readError(response) {
  const text = await response.text()
  try {
    const json = JSON.parse(text)
    return json.error || text.trim()
  } catch {
    return text.trim() || response.statusText || 'Terjadi kesalahan pada server'
  }
}

/** Pesan khusus untuk kegagalan jaringan, dibedakan dari error HTTP. */
function toMessage(err) {
  if (err.name === 'TypeError' || err.message === 'Failed to fetch') {
    return 'Tidak dapat terhubung ke server. Pastikan backend aktif.'
  }
  return err.message
}

async function request(path, options) {
  let response
  try {
    response = await fetch(API_URL + path, options)
  } catch (err) {
    throw new Error(toMessage(err), { cause: err })
  }

  if (!response.ok) {
    throw new Error(await readError(response))
  }

  return response.json()
}

export async function shortenUrl(longUrl) {
  const data = await request('/shorten', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ long_url: longUrl.trim() }),
  })

  if (!data.short_url) {
    throw new Error('Gagal memendekkan URL. Pastikan format tautan valid.')
  }

  return data.short_url
}

export async function getStats(code) {
  return request(`/stats/${encodeURIComponent(code)}`)
}

/**
 * Mengambil kode 6 karakter dari input pengguna. Menerima kode polos
 * ("aB3xY9") maupun URL lengkap ("http://localhost:8000/aB3xY9").
 * Mengembalikan null kalau tidak ada kode yang valid.
 */
export function extractCode(input) {
  const candidate = input.trim().split('/').filter(Boolean).pop() ?? ''
  return /^[A-Za-z0-9]{6}$/.test(candidate) ? candidate : null
}
