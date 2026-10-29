export default function Hero() {
  return (
    <section className="mb-7">
      <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-indigo-50/70 border border-indigo-200/60 mb-3.5">
        <span className="w-1.5 h-1.5 rounded-full bg-indigo-600" />
        <span className="font-mono text-[11px] font-medium tracking-wider text-indigo-950/80 uppercase">
          URL Redirection Engine
        </span>
      </div>

      <h1 className="text-3xl sm:text-[36px] font-semibold tracking-tight leading-[1.16] text-zinc-900 mb-3">
        Ringkas tautan panjang{' '}
        <span className="text-zinc-500 font-normal">tanpa distraksi.</span>
      </h1>
      <p className="text-[15px] text-zinc-600 leading-relaxed">
        Masukkan URL destinasi di bawah untuk menghasilkan tautan pendek yang bersih dan siap dibagikan.
      </p>
    </section>
  )
}
