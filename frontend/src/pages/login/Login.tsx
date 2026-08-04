import './Login.css'
import logo from '@/assets/logo.svg'
import LoginForm from './components/login-form/LoginForm'

export default function Login() {
    return (
        <main>
            <div className="page-center">
                <img src={logo} className="site-logo" />
                <LoginForm />
                <div className="login-link">
                    <p> Don't have an account ?</p>
                    <a href="/register">Sign up</a>
                </div>
            </div>
        </main>
    )
}
