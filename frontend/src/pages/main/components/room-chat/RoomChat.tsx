import './RoomChat.css'
import type { Room } from '@/models/Room'
import RoomChatHeader from './components/room-chat-header/RoomChatHeader'
import RoomChatMessages from './components/room-chat-messages/RoomChatMessages'
import RoomChatInput from './components/room-chat-input/RoomChatInput'

interface RoomChatProps {
    currentRoom: Room
}

export default function RoomChat({ currentRoom }: RoomChatProps) {
    return (
        <section className="chat">
            <RoomChatHeader currentRoom={currentRoom} />
            <RoomChatMessages currentRoom={currentRoom} />
            <RoomChatInput />
        </section>
    )
}
