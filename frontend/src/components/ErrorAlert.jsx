export default function ErrorAlert({ message }) {
  if (!message) return null

  return (
    <div
      role="alert"
      className="mt-2 px-3.5 py-2.5 bg-rose-50/80 border border-rose-200/80 rounded-lg flex items-center gap-2.5"
    >
      <div className="w-4 h-4 rounded-full bg-rose-600 text-white font-mono text-[11px] font-semibold flex items-center justify-center shrink-0">
        !
      </div>
      <div className="text-[13px] text-rose-700 font-medium">{message}</div>
    </div>
  )
}
