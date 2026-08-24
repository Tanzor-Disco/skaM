import './Register.css'
import RegisterForm from './components/register-form/RegisterForm'
import logo from '@/assets/logo.svg'
import { useState } from 'react'

export default function Register() {
    const [registerComplete, setRegisterComplete] = useState(false)
    return (
        <main className="page-register">
            <img src={logo} className="register-site-logo" />
            <RegisterForm setRegisterComplete={setRegisterComplete} />
            {registerComplete && <p>A letter was sent to your email</p>}
        </main>
    )
}
