import './LoginForm.css'
import InputField from '@/components/input-field/InputField'

export default function LoginForm() {
    return (
        <form className="items-align">
            <InputField fieldName={'Email'} />
            <InputField fieldName={'Password'} />
            <button type="submit" className="submit-button">
                Submit
            </button>
        </form>
    )
}
