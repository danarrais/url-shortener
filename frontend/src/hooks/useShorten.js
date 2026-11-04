import { useState } from 'react'
import { shortenUrl } from '../lib/api'

/** Mengelola state form "Pendekkan": input, hasil, loading, dan pesan error. */
export function useShorten() {
  const [longUrl, setLongUrl] = useState('')
  const [shortUrl, setShortUrl] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')

  const submit = async (e) => {
    e.preventDefault()
    if (!longUrl.trim()) return

    setIsLoading(true)
    setError('')
    setShortUrl('')

    try {
      setShortUrl(await shortenUrl(longUrl))
      setLongUrl('')
    } catch (err) {
      setError(err.message)
    } finally {
      setIsLoading(false)
    }
  }

  const reset = () => {
    setLongUrl('')
    setShortUrl('')
    setError('')
  }

  return {
    longUrl,
    setLongUrl,
    shortUrl,
    isLoading,
    error,
    setError,
    submit,
    reset,
  }
}
