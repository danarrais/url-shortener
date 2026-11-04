import { useState } from 'react'
import { extractCode, getStats } from '../lib/api'

/** Mengelola state form "Statistik": input kode/URL, data, loading, error. */
export function useStats() {
  const [input, setInput] = useState('')
  const [data, setData] = useState(null)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    if (!input.trim()) return

    // Pesan ini lebih jelas daripada membiarkan backend membalas 404
    // untuk kode yang panjangnya memang tidak mungkin valid.
    const code = extractCode(input)
    if (!code) {
      setData(null)
      setError('Kode tautan harus 6 karakter huruf atau angka.')
      return
    }

    setIsLoading(true)
    setError('')
    setData(null)

    try {
      setData(await getStats(code))
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  return { input, setInput, data, isLoading, error, setError, submit }
}
