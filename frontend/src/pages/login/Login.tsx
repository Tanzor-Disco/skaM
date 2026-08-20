import './Login.css'
import logo from '@/assets/logo.svg'
import LoginForm from './components/login-form/LoginForm'
import { useState } from 'react'

export default function Login() {
    const [invalidInput, setInvalidInput] = useState(false)
    return (
        <main>
            <div className="page-login">
                <img src={logo} className="site-logo" />
                <LoginForm setInvalidInput={setInvalidInput} />
                <div className="login-link">
                    <p> Don't have an account ?</p>
                    <a href="/register">Sign up</a>
                </div>
                {invalidInput && (
                    <p className="error-message">Wrong login data entered</p>
                )}
            </div>
        </main>
    )
}
