import { Hero } from './components/marketing/Hero'
import { ManifestoSection } from './components/marketing/ManifestoSection'
import { Features } from './components/marketing/Features'
import { CoOpSection } from './components/marketing/CoOpSection'
import { CTA } from './components/marketing/CTA'
import { Navbar } from './components/marketing/Navbar'
import { Footer } from './components/marketing/Footer'

function App() {
  return (
    <div className="min-h-screen bg-gable-bg text-slate-200">
      <Navbar />
      <Hero />
      <ManifestoSection />
      <Features />
      <CoOpSection />
      <CTA />
      <Footer />
    </div>
  )
}

export default App
