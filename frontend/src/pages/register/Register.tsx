import './Register.css'
import RegisterForm from './components/register-form/RegisterForm'
import logo from "@/assets/logo.svg"

export default function Register() {
    return (
        <main className="page-center">
			<img src={logo} className="site-logo"/>
            <RegisterForm />
        </main>
    )
}
