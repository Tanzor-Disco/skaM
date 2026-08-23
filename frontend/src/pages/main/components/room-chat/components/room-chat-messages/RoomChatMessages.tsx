import './RoomChatMessages.css'
import type { Room } from '@/models/Room'
import Message from './components/Message/Message'
import { useEffect, type Dispatch, type SetStateAction } from 'react'
import type { MessageData } from '@/models/MessageData'

interface RoomChatMessagesProps {
    currentRoom: Room
    messages: MessageData[]
    setMessages: Dispatch<SetStateAction<Array<MessageData>>>
}

export default function RoomChatMessages({
    currentRoom,
    messages,
    setMessages,
}: RoomChatMessagesProps) {
    useEffect(() => {
        async function getMessages() {
            const response = await fetch(
                `/api/main/messages?room_id=${currentRoom.ID}`
            )
            const respObj = await response.json()
            if (respObj.data === null) {
                setMessages([])
                return
            }
            setMessages(respObj.data)
        }
        getMessages()
    }, [currentRoom])
    const messagesComponents = messages.map((message) => {
        return (
            <Message
                key={message.ID}
                username={message.Author}
                messageText={message.Text}
            />
        )
    })
    return <div className="messages">{messagesComponents}</div>
}
