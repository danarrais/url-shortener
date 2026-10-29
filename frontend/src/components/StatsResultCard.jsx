import { useState } from 'react'
import { BarChart3, MousePointerClick, Link as LinkIcon, Copy, Check, ExternalLink } from 'lucide-react'

export default function StatsResultCard({ statsData }) {
  const [copied, setCopied] = useState(false)

  if (!statsData) return null

  const handleCopy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Gagal menyalin:', err)
    }
  }

  return (
    <div className="mt-2 p-4 bg-gradient-to-br from-indigo-50/40 via-zinc-50/60 to-sky-50/30 border border-indigo-200/70 rounded-lg">
      {/* Header Badge */}
      <div className="flex justify-between items-center mb-3">
        <span className="inline-flex items-center gap-1.5 font-mono text-[11px] font-medium text-indigo-700 bg-indigo-100/70 border border-indigo-300/60 px-2 py-0.5 rounded">
          <BarChart3 size={12} strokeWidth={2.5} />
          DATA STATISTIK DITEMUKAN
        </span>
        <span className="hidden sm:inline text-xs text-zinc-400">
          Ringkasan performa & destinasi URL
        </span>
      </div>

      {/* Metrics Row */}
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 mb-3">
        {/* Total Clicks Card */}
        <div className="bg-white/80 border border-zinc-200/80 rounded-lg p-3 flex items-center justify-between shadow-2xs">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-md bg-indigo-50 border border-indigo-100/80 flex items-center justify-center text-indigo-600 shrink-0">
              <MousePointerClick size={16} />
            </div>
            <div>
              <p className="text-[11px] font-mono font-medium text-zinc-500 uppercase tracking-wider">
                Total Diklik
              </p>
              <div className="flex items-baseline gap-1">
                <span className="font-mono text-xl font-bold text-zinc-900">
                  {statsData.clicks ?? 0}
                </span>
                <span className="text-xs text-zinc-500 font-medium">kali</span>
              </div>
            </div>
          </div>
        </div>

        {/* Short Code Card */}
        <div className="bg-white/80 border border-zinc-200/80 rounded-lg p-3 flex items-center justify-between shadow-2xs">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-md bg-indigo-50 border border-indigo-100/80 flex items-center justify-center text-indigo-600 font-mono font-bold text-xs shrink-0">
              /
            </div>
            <div>
              <p className="text-[11px] font-mono font-medium text-zinc-500 uppercase tracking-wider">
                Kode Tautan
              </p>
              <span className="font-mono text-sm font-semibold text-zinc-900">
                {statsData.short_code}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Destination URL Card */}
      <div className="bg-white/90 border border-zinc-200/90 rounded-lg p-3 shadow-2xs">
        <div className="flex items-center gap-1.5 text-[11px] font-mono font-medium text-zinc-500 uppercase tracking-wider mb-1.5">
          <LinkIcon size={12} className="text-zinc-400" />
          <span>Tujuan Asli</span>
        </div>

        <div className="flex justify-between items-center gap-3 flex-wrap">
          <a
            href={statsData.long_url}
            target="_blank"
            rel="noopener noreferrer"
            className="font-mono text-[13.5px] font-medium text-indigo-950 hover:text-indigo-600 hover:underline underline-offset-4 transition-colors break-all flex-1 min-w-0"
            title={statsData.long_url}
          >
            {statsData.long_url}
          </a>

          <div className="flex items-center gap-1.5 shrink-0">
            <button
              type="button"
              onClick={() => handleCopy(statsData.long_url)}
              className={`inline-flex items-center gap-1.5 px-3 py-1.5 border rounded-md text-[13px] font-medium transition-all cursor-pointer ${
                copied
                  ? 'border-emerald-500 text-emerald-700 bg-emerald-50'
                  : 'bg-white border-zinc-200/90 text-zinc-800 hover:border-indigo-300 hover:text-indigo-700 shadow-2xs'
              }`}
            >
              {copied ? (
                <>
                  <Check size={14} strokeWidth={2.2} />
                  <span>Tersalin</span>
                </>
              ) : (
                <>
                  <Copy size={14} />
                  <span>Salin URL</span>
                </>
              )}
            </button>

            <a
              href={statsData.long_url}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center justify-center w-8 h-8 bg-white border border-zinc-200/90 hover:border-indigo-300 rounded-md text-zinc-600 hover:text-indigo-600 transition-colors shadow-2xs"
              title="Buka tautan asli di tab baru"
            >
              <ExternalLink size={15} />
            </a>
          </div>
        </div>
      </div>
    </div>
  )
}
