import './RoomChat.css'
import type { Room } from '@/models/Room'
import type { MessageData } from '@/models/MessageData'
import RoomChatHeader from './components/room-chat-header/RoomChatHeader'
import RoomChatMessages from './components/room-chat-messages/RoomChatMessages'
import RoomChatInput from './components/room-chat-input/RoomChatInput'
import { useState, useEffect } from 'react'
import { openWebsocket } from '@/websocket/websocket'

interface RoomChatProps {
    currentRoom: Room
}

export default function RoomChat({ currentRoom }: RoomChatProps) {
    const [messages, setMessages] = useState<never[] | MessageData[]>([])
    useEffect(() => {
        openWebsocket(setMessages)
    })
    return (
        <section className="chat">
            <RoomChatHeader currentRoom={currentRoom} />
            <RoomChatMessages
                currentRoom={currentRoom}
                messages={messages}
                setMessages={setMessages}
            />
            <RoomChatInput currentRoom={currentRoom} />
        </section>
    )
}
