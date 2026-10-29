import { Check, Copy, ExternalLink } from 'lucide-react'

export default function ResultCard({ shortUrl, copied, onCopy }) {
  if (!shortUrl) return null

  return (
    <div className="mt-2 p-4 bg-gradient-to-br from-emerald-50/40 via-zinc-50/60 to-indigo-50/30 border border-emerald-200/70 rounded-lg">
      <div className="flex justify-between items-center mb-2.5">
        <span className="inline-flex items-center gap-1.5 font-mono text-[11px] font-medium text-emerald-700 bg-emerald-100/70 border border-emerald-300/60 px-2 py-0.5 rounded">
          <Check size={12} strokeWidth={2.5} />
          TAUTAN SIAP DIGUNAKAN
        </span>
        <span className="hidden sm:inline text-xs text-zinc-400">
          Klik tombol salin atau buka langsung
        </span>
      </div>

      <div className="flex justify-between items-center gap-4 flex-wrap">
        <a
          href={shortUrl}
          target="_blank"
          rel="noopener noreferrer"
          className="font-mono text-[17px] font-medium text-indigo-950 hover:text-indigo-600 hover:underline underline-offset-4 transition-colors break-all"
        >
          {shortUrl}
        </a>

        <div className="flex items-center gap-1.5">
          <button
            type="button"
            onClick={() => onCopy(shortUrl)}
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
            href={shortUrl}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center justify-center w-8 h-8 bg-white border border-zinc-200/90 hover:border-indigo-300 rounded-md text-zinc-600 hover:text-indigo-600 transition-colors shadow-2xs"
            title="Buka di tab baru"
          >
            <ExternalLink size={15} />
          </a>
        </div>
      </div>
    </div>
  )
}
