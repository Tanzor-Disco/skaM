import './RoomChatMessages.css'
import type { Room } from '@/models/Room'
import Message from './components/Message/Message'

interface RoomChatMessagesProps {
    currentRoom: Room
}

export default function RoomChatMessages({
    currentRoom,
}: RoomChatMessagesProps) {
    return (
        <div className="messages">
            <Message username={'huylo'} messageText={'Putin Huylo blya'} />
        </div>
    )
}
