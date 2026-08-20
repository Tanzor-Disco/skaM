import './RoomListHeader.css'
import { useState, type Dispatch, type SetStateAction } from 'react'
import RoomRegister from './components/room-register/RoomRegister'
import type { Room } from '@/models/Room'

interface RoomListHeaderProps {
    setRooms: Dispatch<SetStateAction<Room[]>>
}

export default function RoomListHeader({ setRooms }: RoomListHeaderProps) {
    const [roomRegisterVisible, setRoomRegisterVisible] = useState(false)
    function handleClick() {
        setRoomRegisterVisible(true)
    }
    return (
        <>
            <header className="room-list-header">
                <h1>Home</h1>
                <button onClick={handleClick}>Add a room</button>
            </header>
            {roomRegisterVisible && (
                <RoomRegister
                    setRooms={setRooms}
                    setRoomRegisterVisible={setRoomRegisterVisible}
                />
            )}
        </>
    )
}
