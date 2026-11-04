import { useEffect, useRef, useState } from 'react'

/**
 * Menyalin teks ke clipboard dan menandai "tersalin" selama `duration` ms.
 * Timer dibersihkan saat komponen dilepas agar tidak memanggil setState
 * pada komponen yang sudah tidak ada.
 */
export function useClipboard(duration = 2000) {
  const [copied, setCopied] = useState(false)
  const timer = useRef(null)

  useEffect(() => () => clearTimeout(timer.current), [])

  const copy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      clearTimeout(timer.current)
      timer.current = setTimeout(() => setCopied(false), duration)
    } catch (err) {
      console.error('Gagal menyalin:', err)
    }
  }

  return { copied, copy }
}
