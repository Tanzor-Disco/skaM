import { useEffect, useState, type Dispatch, type SetStateAction } from 'react'
import './RoomList.css'

import RoomEntry from './components/room-entry/RoomEntry'
import RoomListHeader from './components/room-header/RoomListHeader'
import logo from '@/assets/logo.svg'
import type { Room } from '@/models/Room'

interface RoomListProps {
    setCurrentRoom: Dispatch<SetStateAction<Room | null>>
}

export default function RoomList({ setCurrentRoom }: RoomListProps) {
    const [rooms, setRooms] = useState<Room[]>([])
    useEffect(() => {
        async function getRooms() {
            const response = await fetch('api/main/rooms')
            if (!response.ok) {
                return []
            }
            const srvRooms = await response.json()
            if (srvRooms.data === null) {
                return
            }
            setRooms(srvRooms.data)
            setCurrentRoom(srvRooms.data[0])
        }
        getRooms()
    }, [])

    const roomComponents = rooms.map((room: Room) => {
        return (
            <RoomEntry
                key={room.ID}
                roomID={room.ID}
                roomIcon={logo}
                roomTitle={room.Name}
                setCurrentRoom={setCurrentRoom}
            />
        )
    })

    return (
        <section className="room-list">
            <RoomListHeader setRooms={setRooms} />
            <div className="room-entries">{roomComponents}</div>
        </section>
    )
}
