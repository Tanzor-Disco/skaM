import './InviteInfo.css'
import logo from '@/assets/logo.svg'
import { useEffect, useState } from 'react'
import type { Room } from '@/models/Room'
import { useNavigate } from 'react-router'

export default function InviteInfo() {
    const token = new URLSearchParams(window.location.search).get('token')
    const navigate = useNavigate()
    const [room, setRoom] = useState<Room | null>(null)
    useEffect(() => {
        async function fetchRoomData() {
            const response = await fetch(`api/info/room?token=${token}`)
            if (!response.ok) {
                navigate('/login')
                return
            }
            const respObj = await response.json()
            setRoom(respObj.data[0])
        }
        fetchRoomData()
    }, [])

    async function handleClick() {
        const response = await fetch('api/invite/add', {
            method: 'POST',
            headers: { 'Context-Type': 'application/json' },
            body: JSON.stringify({
                token: token,
            }),
        })
        const respObj = await response.json()

        if (!response.ok && respObj.error_kind !== 'ERR_UNIQUE_VIOLATION') {
            navigate('/login')
            return
        }
        navigate('/main')
    }
    return (
        <main className="invite-info">
            <img src={logo} className="invite-img" />
            {room && <h1 className="room-name">{room.Name}</h1>}
            <button className="join-button" onClick={handleClick}>
                Join
            </button>
        </main>
    )
}
