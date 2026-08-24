import './RoomChatHeader.css'
import type { Room } from '@/models/Room'
import { useState } from 'react'
import RoomInfo from './components/RoomInfo/RoomInfo'

interface RoomChatHeaderProps {
    currentRoom: Room
}

export default function RoomChatHeader({ currentRoom }: RoomChatHeaderProps) {
    const [roomInfoVisible, setRoomInfoVisible] = useState(false)
    async function handleClick() {
        setRoomInfoVisible(true)
    }
    return (
        <>
            <button className="room-chat-header" onClick={handleClick}>
                {currentRoom.Name}
            </button>
            {roomInfoVisible && (
                <RoomInfo
                    currentRoom={currentRoom}
                    setRoomInfoVisible={setRoomInfoVisible}
                />
            )}
        </>
    )
}
