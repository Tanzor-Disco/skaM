import './RoomChatHeader.css'
import type { Room } from '@/models/Room'

interface RoomChatHeaderProps {
    currentRoom: Room
}

export default function RoomChatHeader({ currentRoom }: RoomChatHeaderProps) {
    return (
        <header className="room-chat-header">
            <h1> {currentRoom.Name} </h1>
        </header>
    )
}
