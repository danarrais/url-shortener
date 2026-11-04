import { useState } from 'react'
import Header from './components/Header'
import Hero from './components/Hero'
import UrlForm from './components/UrlForm'
import ErrorAlert from './components/ErrorAlert'
import ResultCard from './components/ResultCard'
import FeaturesFooter from './components/FeaturesFooter'
import StatsForm from './components/StatsForm'
import StatsResultCard from './components/StatsResultCard'
import { useShorten } from './hooks/useShorten'
import { useStats } from './hooks/useStats'

function App() {
  const [activeTab, setActiveTab] = useState('shorten')
  const shorten = useShorten()
  const stats = useStats()

  const switchTab = (tab) => {
    setActiveTab(tab)
    shorten.setError('')
    stats.setError('')
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
            {activeTab === 'shorten' && (
              <>
                <UrlForm
                  longUrl={shorten.longUrl}
                  setLongUrl={shorten.setLongUrl}
                  onSubmit={shorten.submit}
                  isLoading={shorten.isLoading}
                  onClearError={() => shorten.setError('')}
                  onPasteError={() => shorten.setError('Izin clipboard ditolak oleh browser.')}
                />
                <ErrorAlert message={shorten.error} />
                <ResultCard shortUrl={shorten.shortUrl} />
              </>
            )}

            {activeTab === 'stats' && (
              <>
                <StatsForm
                  statsInput={stats.input}
                  setStatsInput={stats.setInput}
                  onSubmit={stats.submit}
                  isLoading={stats.isLoading}
                  onClearError={() => stats.setError('')}
                />
                <ErrorAlert message={stats.error} />
                <StatsResultCard statsData={stats.data} />
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
