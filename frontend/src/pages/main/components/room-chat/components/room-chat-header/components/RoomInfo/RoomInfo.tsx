import type { Dispatch, SetStateAction } from 'react'
import './RoomInfo.css'
import logo from '@/assets/logo.svg'
import type { Room } from '@/models/Room'
import RoomLink from './components/room-link/RoomLink'

interface RoomInfoProps {
    currentRoom: Room
    setRoomInfoVisible: Dispatch<SetStateAction<boolean>>
}

export default function RoomInfo({
    currentRoom,
    setRoomInfoVisible,
}: RoomInfoProps) {
    function handleRoomExitClick() {
        setRoomInfoVisible(false)
    }
    return (
        <div className="overlay-container">
            <div className="room-info-container">
                <img className="room-info-img" src={logo} />
                <h1>{currentRoom.Name}</h1>
                <RoomLink currentRoom={currentRoom} />
                <button
                    className="room-info-exit"
                    onClick={handleRoomExitClick}
                >
                    Exit
                </button>
            </div>
        </div>
    )
}
