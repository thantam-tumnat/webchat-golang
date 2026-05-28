import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { UsernameGate } from './components/UsernameGate'
import { RoomsPage } from './pages/RoomsPage'
import { ChatPage } from './pages/ChatPage'

export default function App() {
  return (
    // UsernameGate บังคับตั้งชื่อก่อน แล้วค่อยเข้าถึง routes
    <UsernameGate>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<RoomsPage />} />
          <Route path="/rooms/:id" element={<ChatPage />} />
        </Routes>
      </BrowserRouter>
    </UsernameGate>
  )
}
