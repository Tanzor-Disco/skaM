import InputField from "@/components/input-field/InputField"
import "./RoomRegister.css"
import type {Dispatch,SetStateAction} from "react"

interface RoomRegisterProps {
	setRoomRegisterVisible:Dispatch<SetStateAction<boolean>>
}

export default function RoomRegister({setRoomRegisterVisible}:RoomRegisterProps) {
	function handleCancelClick() {
		setRoomRegisterVisible(false)
	}
	async function handleSubmit(formData:FormData) {
		const formObj = Object.fromEntries(formData)
		const headers = new Headers
		headers.append("Content-Type", "application/json")
		const response = await fetch("/api/main/new/room", {
			method:"POST",
			headers:headers,
			body:JSON.stringify(formObj)
		})

	}

	return (
		<div className="page-overlay">
			<form className="room-register-form" action={handleSubmit}>
				<h1>Create a room</h1>
				<InputField fieldName={"Name"} inputType={"text"} required={true}/>
				<div className="register-form-buttons">
					<button onClick={handleCancelClick} type="button">
						Cancel
					</button>
					<button>Create</button>
				</div>

			</form>
		</div>
	)
}
