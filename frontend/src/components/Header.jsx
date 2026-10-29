export default function Header() {
  return (
    <header className="sticky top-0 z-10 bg-[#fafaf9]/80 backdrop-blur-md border-b border-zinc-200/80">
      <div className="max-w-[640px] mx-auto px-5 py-3.5 flex items-center justify-between">
        <div className="flex items-center gap-1.5 font-mono text-sm font-medium">
          <span className="text-indigo-600 font-bold">/</span>
          <span className="text-zinc-900 tracking-tight">pendek.in</span>
          <span className="text-[11px] px-1.5 py-0.5 bg-indigo-50/80 border border-indigo-200/60 rounded text-indigo-600 ml-1">
            v1.0
          </span>
        </div>
      </div>
    </header>
  )
}
