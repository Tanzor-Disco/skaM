import './RegisterForm.css'
import InputField from '@/components/input-field/InputField'
import {useState} from "react"

export default function RegisterForm() {
	const [warningVisible,setWarningVisible] = useState(false)

	async function handleSubmit(formData:FormData) {
		const objData = Object.fromEntries(formData)
		const response = await fetch("/api/register", {
			method:"POST",
			headers: {
				"Content-Type":"application/json"
			},
			body:JSON.stringify(objData),
		})
		if (!response.ok) {
			const responseObj = await response.json()
			if (responseObj.error_kind === "ERR_EMAIL_TAKEN") {
				setWarningVisible(true)
				return 
			}
		}
		setWarningVisible(false)
	}
	
    return (
        <form className="register-form" action={handleSubmit}>
			<div className={warningVisible ? "smaller-margin" : undefined}>
				<InputField fieldName={'Email'} inputType={"email"} required={true} />
				{warningVisible && <p className="error-message">The email is already taken</p>}
			</div>
            <InputField fieldName={'Username'} inputType={"text"} required={true} />
            <InputField fieldName={'Password'} inputType={"text"} required={true} />
            <button className="form-submit">Submit</button>
        </form>
    )
}
