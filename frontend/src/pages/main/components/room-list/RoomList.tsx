import "./RoomList.css"

import RoomEntry from "./components/room-entry/RoomEntry"
import logo from "@/assets/logo.svg"

export default function RoomList() {
	return (
		<section className="room-list">
			<h1>Home</h1>
			<div className="room-entries">
				<RoomEntry roomIcon={logo} roomTitle={"Gei Mira"} />
				<RoomEntry roomIcon={logo} roomTitle={"Putinskiye Sokoli"} />
			</div>
		</section>
	)
}
