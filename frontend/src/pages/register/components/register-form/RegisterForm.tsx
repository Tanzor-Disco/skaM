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
			if (responseObj.message === "The email is already taken") {
				setWarningVisible(prevWarning => !prevWarning)
				return 
			}
		}
		setWarningVisible(prevWarning => !prevWarning)
	}
	
    return (
        <form className="register-form" action={handleSubmit}>
			<div className={warningVisible ? "smaller-margin" : undefined}>
				<InputField fieldName={'Email'} required={true} />
				{warningVisible && <p className="error-message">The email is already taken</p>}
			</div>
            <InputField fieldName={'Username'} required={true} />
            <InputField fieldName={'Password'} required={true} />
            <button className="form-submit">Submit</button>
        </form>
    )
}
