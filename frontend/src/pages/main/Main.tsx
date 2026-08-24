import './Main.css'

import RoomList from './components/room-list/RoomList'
import RoomChat from './components/room-chat/RoomChat'
import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router'
import type { Room } from '@/models/Room'

export default function Main() {
    const navigate = useNavigate()
    const [userID, setUserID] = useState(0)
    const [currentRoom, setCurrentRoom] = useState<Room | null>(null)
    useEffect(() => {
        async function getUserData() {
            const response = await fetch('api/main')
            if (!response.ok) {
                navigate('/login')
                return
            }
            const body = await response.json()
            setUserID(body.data[0])
        }
        getUserData()
    })
    return (
        <main className="page-main">
            <RoomList setCurrentRoom={setCurrentRoom} />
            {currentRoom && <RoomChat currentRoom={currentRoom} />}
        </main>
    )
}
