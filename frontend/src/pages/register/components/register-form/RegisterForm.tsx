import './RegisterForm.css'
import InputField from '@/components/input-field/InputField'

export default function RegisterForm() {
    return (
        <form className="register-form">
            <InputField fieldName={'Email'} />
            <InputField fieldName={'Username'} />
            <InputField fieldName={'Password'} />
            <button className="form-submit">Submit</button>
        </form>
    )
}
