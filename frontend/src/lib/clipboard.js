/** Membaca clipboard. Mengembalikan string kosong kalau izin ditolak. */
export async function readClipboard() {
  try {
    return await navigator.clipboard.readText()
  } catch {
    return ''
  }
}
