import './RoomLink.css'
import { useEffect, useState } from 'react'
import type { Room } from '@/models/Room'

interface RoomLinkProps {
    currentRoom: Room
}

export default function RoomLink({ currentRoom }: RoomLinkProps) {
    const [link, setLink] = useState('')
    useEffect(() => {
        async function getInviteURL() {
            const response = await fetch(
                `/api/main/room/invite?room_id=${currentRoom.ID}`
            )
            if (!response.ok) {
                console.warn(
                    `response returned with an error: ${response.status}`
                )
                return
            }
            const respObj = await response.json()
            setLink(respObj.data[0])
        }
        getInviteURL()
    }, [currentRoom])
    const [buttonText, setButtonText] = useState('Copy Link')
    async function handleClick() {
        await navigator.clipboard.writeText(link)
        setButtonText('Link copied')
        setTimeout(() => {
            setButtonText('Copy link')
        }, 2000)
    }
    return (
        <div className="room-link">
            <p className="room-link-text">{link}</p>
            <button className="link-copy" onClick={handleClick}>
                {buttonText}
            </button>
        </div>
    )
}
