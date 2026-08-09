import './RegisterForm.css'
import InputField from '@/components/input-field/InputField'
import FormField from "./components/FormField"
import {useState} from "react"


export default function RegisterForm() {
	const [emailErr,setEmailErr] = useState("")
	const [usernameErr,setUsernameErr] = useState("")
	const [passwordErr,setPasswordErr] = useState("")


	function toggleState(errorKind:string) {
		switch (errorKind) {
			case "ERR_EMAIL_TAKEN":
				setEmailErr("Email is already taken")
				break
			case "ERR_INVALID_USERNAME_LENGTH":
				setUsernameErr("The maximum length of the username is 20 characters")
				break
			case "ERR_FORBIDDEN_PASSWORD_CHARS":
				setPasswordErr("The password should contain only latin characters and digits")
				break
			case "ERR_INVALID_PASSWORD_LENGTH":
				setPasswordErr("The maximum length of the password is 70 characters")
		}
	}
	function toggleAllOff() {
		setEmailErr("")
		setUsernameErr("")
	}

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
			toggleAllOff()
			toggleState(responseObj.error_kind)
			return
		}
	}
	
    return (
        <form className="register-form" action={handleSubmit}>
			<FormField 
			warning={emailErr} 
			fieldName={"Email"} 
			inputType={"email"}
			/>
			<FormField 
			warning={usernameErr}
			fieldName={"Username"}
			inputType={"text"}
			/>
			<FormField 
			warning={passwordErr}
			fieldName={"Password"}
			inputType={"text"}
			/>
        	<button className="form-submit">Submit</button>
        </form>
    )
}
