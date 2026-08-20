import { useEffect, useState } from 'react'
import './RoomList.css'

import RoomEntry from './components/room-entry/RoomEntry'
import RoomListHeader from './components/room-header/RoomListHeader'
import logo from '@/assets/logo.svg'
import type { Room } from '@/models/Room'

export default function RoomList() {
    const [rooms, setRooms] = useState<Room[]>([])
    useEffect(() => {
        async function getRooms() {
            const response = await fetch('api/main/rooms')
            if (!response.ok) {
                return []
            }
            const srvRooms = await response.json()
            setRooms(srvRooms.data)
        }
        getRooms()
    }, [])

    const roomComponents = rooms.map((room: Room) => {
        return <RoomEntry key={room.ID} roomIcon={logo} roomTitle={room.Name} />
    })

    return (
        <section className="room-list">
            <RoomListHeader setRooms={setRooms} />
            <div className="room-entries">{roomComponents}</div>
        </section>
    )
}
