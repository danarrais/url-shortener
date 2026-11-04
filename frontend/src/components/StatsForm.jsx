import { useRef } from 'react'
import { BarChart3, ArrowRight } from 'lucide-react'
import { readClipboard } from '../lib/clipboard'

export default function StatsForm({
  statsInput,
  setStatsInput,
  onSubmit,
  isLoading,
  onClearError,
}) {
  const inputRef = useRef(null)

  const handlePaste = async () => {
    const text = await readClipboard()
    if (!text) return
    setStatsInput(text)
    if (onClearError) onClearError()
    inputRef.current?.focus()
  }

  const handleClear = () => {
    setStatsInput('')
    if (onClearError) onClearError()
    inputRef.current?.focus()
  }

  return (
    <form onSubmit={onSubmit} className="flex flex-col sm:flex-row gap-2 items-stretch">
      <div className="group flex-1 flex items-center gap-2.5 px-3 bg-zinc-50/70 border border-zinc-200/90 rounded-lg focus-within:bg-white focus-within:border-indigo-500 focus-within:ring-3 focus-within:ring-indigo-500/10 transition-all">
        <BarChart3 size={16} className="text-zinc-400 group-focus-within:text-indigo-600 transition-colors shrink-0" />

        <input
          ref={inputRef}
          type="text"
          className="w-full border-none bg-transparent py-3 text-[14.5px] text-zinc-900 placeholder:text-zinc-400 focus:outline-none"
          placeholder="Masukkan kode atau link pendek (misal: ClJBRq)..."
          value={statsInput}
          onChange={(e) => {
            setStatsInput(e.target.value)
            if (onClearError) onClearError()
          }}
          required
        />

        <div className="flex items-center shrink-0">
          {statsInput ? (
            <button
              type="button"
              onClick={handleClear}
              className="bg-white border border-zinc-200 hover:border-zinc-300 text-zinc-500 hover:text-zinc-800 font-mono text-[11px] px-2 py-1 rounded transition-colors cursor-pointer"
              title="Bersihkan input"
            >
              Clear
            </button>
          ) : (
            <button
              type="button"
              onClick={handlePaste}
              className="bg-white border border-zinc-200 hover:border-indigo-200 hover:bg-indigo-50/40 text-zinc-600 hover:text-indigo-700 font-mono text-[11px] px-2 py-1 rounded transition-colors cursor-pointer"
              title="Tempel dari clipboard"
            >
              Paste
            </button>
          )}
        </div>
      </div>

      <button
        type="submit"
        disabled={isLoading || !statsInput.trim()}
        className="inline-flex items-center justify-center gap-2 bg-indigo-600 hover:bg-indigo-700 text-white border border-indigo-600 hover:border-indigo-700 px-5 py-3 sm:py-0 rounded-lg text-sm font-medium whitespace-nowrap shadow-xs transition-all active:translate-y-px disabled:bg-zinc-100 disabled:border-zinc-200 disabled:text-zinc-400 disabled:shadow-none disabled:cursor-not-allowed cursor-pointer"
      >
        {isLoading ? (
          <>
            <span className="w-3.5 h-3.5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
            <span>Mencari</span>
          </>
        ) : (
          <>
            <span>Lihat Statistik</span>
            <ArrowRight size={14} />
          </>
        )}
      </button>
    </form>
  )
}
