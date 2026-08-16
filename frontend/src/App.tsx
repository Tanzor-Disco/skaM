import './App.css'
import { BrowserRouter, Routes, Route, Navigate} from 'react-router'

import Login from './pages/login/Login'
import Register from './pages/register/Register'
import Main from "./pages/main/Main"

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
				<Route path="/" element={<Navigate to="/login" replace />} />
				<Route path="/main" element={<Main />} />
            </Routes>
        </BrowserRouter>
    )
}

export default App
