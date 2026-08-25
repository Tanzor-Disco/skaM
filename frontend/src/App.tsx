import './App.css'
import { BrowserRouter, Routes, Route } from 'react-router'

import Login from './pages/login/Login'
import Register from './pages/register/Register'
import Main from './pages/main/Main'
import Root from './pages/root/Root'
import InviteInfo from './pages/invite/InviteInfo'

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/" element={<Root />} />
                <Route path="/main" element={<Main />} />
                <Route path="invite" element={<InviteInfo />} />
            </Routes>
        </BrowserRouter>
    )
}

export default App
