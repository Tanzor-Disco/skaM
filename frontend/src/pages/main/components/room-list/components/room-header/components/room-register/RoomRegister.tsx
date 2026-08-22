import FormField from '@/components/form-field/FormField'
import './RoomRegister.css'
import { useState, type Dispatch, type SetStateAction } from 'react'
import type { Room } from '@/models/Room'
import { newRoom } from '@/models/Room'

interface RoomRegisterProps {
    setRoomRegisterVisible: Dispatch<SetStateAction<boolean>>
    setRooms: Dispatch<SetStateAction<Room[]>>
}

export default function RoomRegister({
    setRoomRegisterVisible,
    setRooms,
}: RoomRegisterProps) {
    const [warning, setWarning] = useState('')

    function handleCancelClick() {
        setRoomRegisterVisible(false)
    }

    async function handleSubmit(formData: FormData) {
        const formObj = Object.fromEntries(formData)
        const headers = new Headers()
        headers.append('Content-Type', 'application/json')
        const response = await fetch('/api/main/new/room', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(formObj),
        })
        const responseObj = await response.json()
        if (!response.ok) {
            switch (responseObj.error_kind) {
                case 'ERR_INVALID_ROOM_NAME_LENGTH':
                    setWarning(
                        'The length of the name should be no more than 20 characters'
                    )
            }
            return
        }
        if (typeof formObj.name == 'string') {
            const room = newRoom(responseObj.data[0], formObj.name)
            setRooms((prevRooms) => [...prevRooms, room])
        } else {
            console.warn("new entry couldn't be created: wrong type")
        }
        setWarning('')
        setRoomRegisterVisible(false)
    }

    return (
        <div className="page-overlay">
            <form className="room-register-form" action={handleSubmit}>
                <h1>Create a room</h1>
                <FormField
                    warning={warning}
                    fieldName={'Name'}
                    inputType={'text'}
                />
                <div className="register-form-buttons">
                    <button onClick={handleCancelClick} type="button">
                        Cancel
                    </button>
                    <button>Create</button>
                </div>
            </form>
        </div>
    )
}
