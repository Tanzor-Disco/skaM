import "./RoomList.css"

import RoomEntry from "./components/room-entry/RoomEntry"
import RoomListHeader from "./components/room-header/RoomListHeader"
import logo from "@/assets/logo.svg"

export default function RoomList() {
	const response = fetch("/api/main/rooms")
	return (
		<section className="room-list">
			<RoomListHeader />
			<div className="room-entries">
				<RoomEntry roomIcon={logo} roomTitle={"Gei Mira"} />
				<RoomEntry roomIcon={logo} roomTitle={"Putinskiye Sokoli"} />
			</div>
		</section>
	)
}
