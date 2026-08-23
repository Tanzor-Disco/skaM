import type { Dispatch, SetStateAction } from 'react'
import type { MessageData } from '@/models/MessageData'
import './RoomChatInput.css'
import type { Room } from '@/models/Room'

interface RoomChatInputProps {
    setMessages: Dispatch<SetStateAction<Array<MessageData>>>
    currentRoom: Room
}

export default function RoomChatInput({
    setMessages,
    currentRoom,
}: RoomChatInputProps) {
    async function handleSubmit(formData: FormData) {
        const data = {
            ...Object.fromEntries(formData),
            roomID: currentRoom.ID,
        }
        const response = await fetch('/api/main/new/message', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        })
        if (!response.ok) {
            console.warn("Didn't recieve StatusOK", response)
            return
        }
        const respObj = await response.json()
        const newMessage = respObj.data[0]
        setMessages((prevMessages) => [...prevMessages, newMessage])
    }
    return (
        <form className="room-chat-form" action={handleSubmit}>
            <input className="room-chat-input" name="text" required={true} />
            <button> Send </button>
        </form>
    )
}
