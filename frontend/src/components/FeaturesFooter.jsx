export default function FeaturesFooter() {
  return (
    <div className="mt-5 flex flex-wrap items-center justify-between gap-4 px-1 text-xs text-zinc-500">
      <div className="flex items-center gap-5">
        <span className="inline-flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 rounded-full bg-emerald-500" />
          Redirect instan
        </span>
        <span className="inline-flex items-center gap-1.5">
          <span className="w-1.5 h-1.5 rounded-full bg-indigo-400" />
          Bebas iklan & pelacak
        </span>
      </div>
      <span className="font-mono text-[11px] text-zinc-400 bg-zinc-100/80 px-2 py-0.5 rounded border border-zinc-200/60">
        HTTP 302
      </span>
    </div>
  )
}
