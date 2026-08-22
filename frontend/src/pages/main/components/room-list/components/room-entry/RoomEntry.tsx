import type { Dispatch, SetStateAction } from 'react'
import type { Room } from '@/models/Room'
import './RoomEntry.css'

interface RoomEntryProps {
    roomID: number
    roomIcon: string
    roomTitle: string
    setCurrentRoom: Dispatch<SetStateAction<Room | null>>
}

export default function RoomEntry({
    roomID,
    roomIcon,
    roomTitle,
    setCurrentRoom,
}: RoomEntryProps) {
    function handleClick() {
        const room = {
            ID: roomID,
            Name: roomTitle,
        }
        setCurrentRoom(room)
    }
    return (
        <button className="room-entry" onClick={handleClick}>
            <img src={roomIcon} className="room-entry-icon" />
            <span className="room-entry-title">{roomTitle}</span>
        </button>
    )
}
