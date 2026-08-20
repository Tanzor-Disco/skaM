import './Main.css'

import RoomList from './components/room-list/RoomList'
import RoomChat from './components/room-chat/RoomChat'

export default function Main() {
    return (
        <main className="page-main">
            <RoomList />
            <RoomChat />
        </main>
    )
}
