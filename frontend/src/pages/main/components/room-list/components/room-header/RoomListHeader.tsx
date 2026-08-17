import "./RoomListHeader.css"
import {useState} from "react"
import RoomRegister from "./components/room-register/RoomRegister"

export default function RoomListHeader() {
	const [roomRegisterVisible, setRoomRegisterVisible] = useState(false)
	function handleClick() {
		setRoomRegisterVisible(true)
	}
	return (
		<>
		<header className="room-list-header">
				<h1>Home</h1>
				<button onClick={handleClick}>Add a room</button>
		</header>
		{roomRegisterVisible && <RoomRegister setRoomRegisterVisible={setRoomRegisterVisible} />}
		</>
	)
}

