import { Routes, Route } from 'react-router-dom'
import { AppShell } from './components/layout/AppShell'
import { Dashboard } from './pages/Dashboard'
import { Inventory } from './pages/Inventory'
import { Locations } from './pages/Locations'

function App() {
  return (
    <AppShell>
      <Routes>
        <Route path="/" element={<Dashboard />} />
        <Route path="/inventory" element={<Inventory />} />
        <Route path="/locations" element={<Locations />} />
      </Routes>
    </AppShell>
  )
}

export default App
