import "./RoomEntry.css"

interface RoomEntryProps {
	roomIcon:string;
	roomTitle:string;
}

export default function RoomEntry({roomIcon,roomTitle}:RoomEntryProps) {
	return (
		<button className="room-entry">
			<img src={roomIcon} className="room-entry-icon" />
			<span className="room-entry-title">{roomTitle}</span>
		</button>
	)
}
