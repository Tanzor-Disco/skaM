import './LoginForm.css'
import InputField from '@/components/input-field/InputField'
import type {Dispatch,SetStateAction} from "react"

interface LoginFormProps {
	setInvalidInput: Dispatch<SetStateAction<boolean>>
}

export default function LoginForm({setInvalidInput}:LoginFormProps) {

	async function handleLoginSubmit(formData:FormData) {
		const formObj = Object.fromEntries(formData)
		const headers = new Headers()
		headers.append("Content-Type","application/json")
		const servResponse = await fetch("api/login", {
			method: "POST",
			headers: headers,
			body: JSON.stringify(formObj),
		})
		const servResponseObj = await servResponse.json()
		switch (servResponseObj.error_kind) {
			case "ERR_WRONG_LOGIN_DATA":
				setInvalidInput(true)
				break
		}
	}
    return (
        <form className="items-align" action={handleLoginSubmit}>
            <InputField fieldName={'Email'} required={true} inputType={"email"}/>
            <InputField fieldName={'Password'} required={true} inputType={"password"}/>
            <button type="submit" className="submit-button">
                Submit
            </button>
        </form>
    )
}
