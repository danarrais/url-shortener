import { useState, useRef } from 'react'
import Header from './components/Header'
import Hero from './components/Hero'
import UrlForm from './components/UrlForm'
import ErrorAlert from './components/ErrorAlert'
import ResultCard from './components/ResultCard'
import FeaturesFooter from './components/FeaturesFooter'
import StatsForm from './components/StatsForm'
import StatsResultCard from './components/StatsResultCard'

function App() {
  // State Navigasi
  const [activeTab, setActiveTab] = useState('shorten') // 'shorten' atau 'stats'

  // State Shorten
  const [longUrl, setLongUrl] = useState('')
  const [shortUrl, setShortUrl] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const [errorMsg, setErrorMsg] = useState('')
  const [copied, setCopied] = useState(false)
  const inputRef = useRef(null)

  // State Stats
  const [statsInput, setStatsInput] = useState('')
  const [statsData, setStatsData] = useState(null)
  const [isStatsLoading, setIsStatsLoading] = useState(false)
  const [statsError, setStatsError] = useState('')

  // Handle Paste & Copy
  const handlePaste = async () => {
    try {
      const text = await navigator.clipboard.readText()
      if (text) {
        setLongUrl(text)
        setErrorMsg('')
        inputRef.current?.focus()
      }
    } catch {
      setErrorMsg('Izin clipboard ditolak oleh browser.')
    }
  }

  const handleCopy = async (text) => {
    try {
      await navigator.clipboard.writeText(text)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch (err) {
      console.error('Gagal menyalin:', err)
    }
  }

  // Submit Shorten
  const handleShortenSubmit = async (e) => {
    e.preventDefault()
    if (!longUrl.trim()) return

    setIsLoading(true)
    setErrorMsg('')
    setShortUrl('')

    try {
      const response = await fetch('http://localhost:8000/shorten', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ long_url: longUrl.trim() })
      })
      const data = await response.json()

      if (response.ok && data.short_url) {
        setShortUrl(data.short_url)
        setLongUrl('')
      } else {
        setErrorMsg(data.detail || 'Gagal memendekkan URL. Pastikan format tautan valid.')
      }
    } catch {
      setErrorMsg('Tidak dapat terhubung ke server. Pastikan backend aktif.')
    } finally {
      setIsLoading(false)
    }
  }

  // Submit Stats
  const handleStatsSubmit = async (e) => {
    e.preventDefault()
    if (!statsInput.trim()) return

    setIsStatsLoading(true)
    setStatsError('')
    setStatsData(null)

    try {
      // Ambil kode unik jika user menempelkan full URL
      const code = statsInput.trim().split('/').filter(Boolean).pop()
      
      const response = await fetch(`http://localhost:8000/stats/${code}`)
      if (response.ok) {
        const data = await response.json()
        setStatsData(data)
      } else {
        setStatsError("Tautan pendek tidak ditemukan di database.")
      }
    } catch {
      setStatsError("Tidak dapat terhubung ke server backend.")
    } finally {
      setIsStatsLoading(false)
    }
  }

  // Reset function saat pindah tab
  const switchTab = (tab) => {
    setActiveTab(tab)
    setErrorMsg('')
    setStatsError('')
  }

  return (
    <div className="relative min-h-screen bg-[#fafaf9] bg-[radial-gradient(#e4e4e7_1px,transparent_1px)] [background-size:24px_24px] text-zinc-900 font-sans antialiased flex flex-col overflow-hidden selection:bg-indigo-100 selection:text-indigo-900">
      
      {/* Ambient Glow */}
      <div className="pointer-events-none fixed inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-64" aria-hidden="true">
        <div className="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-indigo-200/40 via-sky-100/30 to-emerald-100/30 sm:left-[calc(50%-15rem)] sm:w-[64rem]" />
      </div>

      <Header />

      <main className="w-full max-w-[640px] mx-auto px-5 flex-1 flex flex-col justify-center pb-16">
        <Hero />

        {/* Command Surface */}
        <section className="relative bg-white border border-zinc-200/90 rounded-xl shadow-[0_2px_8px_-2px_rgba(24,24,27,0.05),0_1px_2px_rgba(24,24,27,0.03)] transition-all">
          <div className="absolute inset-x-8 -top-px h-px bg-gradient-to-r from-transparent via-indigo-500/40 to-transparent" />

          {/* Tab Navigation Menu */}
          <div className="flex p-1.5 border-b border-zinc-100 bg-zinc-50/50 rounded-t-xl">
            <button 
              onClick={() => switchTab('shorten')}
              className={`flex-1 py-2 px-3 text-sm font-medium rounded-lg transition-all ${activeTab === 'shorten' ? 'bg-white text-zinc-900 shadow-sm border border-zinc-200/50' : 'text-zinc-500 hover:text-zinc-700 hover:bg-zinc-100/50 border border-transparent'}`}
            >
              Pendekkan
            </button>
            <button 
              onClick={() => switchTab('stats')}
              className={`flex-1 py-2 px-3 text-sm font-medium rounded-lg transition-all ${activeTab === 'stats' ? 'bg-white text-zinc-900 shadow-sm border border-zinc-200/50' : 'text-zinc-500 hover:text-zinc-700 hover:bg-zinc-100/50 border border-transparent'}`}
            >
              Statistik
            </button>
          </div>

          {/* Tab Content */}
          <div className="p-2">
            {/* Render Mode: Shorten */}
            {activeTab === 'shorten' && (
              <>
                <UrlForm
                  longUrl={longUrl}
                  setLongUrl={setLongUrl}
                  onSubmit={handleShortenSubmit}
                  onPaste={handlePaste}
                  isLoading={isLoading}
                  inputRef={inputRef}
                  onClearError={() => setErrorMsg('')}
                />
                <ErrorAlert message={errorMsg} />
                <ResultCard
                  shortUrl={shortUrl}
                  copied={copied}
                  onCopy={handleCopy}
                />
              </>
            )}

            {/* Render Mode: Stats */}
            {activeTab === 'stats' && (
              <>
                <StatsForm
                  statsInput={statsInput}
                  setStatsInput={setStatsInput}
                  onSubmit={handleStatsSubmit}
                  isLoading={isStatsLoading}
                  onClearError={() => setStatsError('')}
                />
                <ErrorAlert message={statsError} />
                <StatsResultCard statsData={statsData} />
              </>
            )}
          </div>

        </section>

        <FeaturesFooter />
      </main>
    </div>
  )
}

export default App