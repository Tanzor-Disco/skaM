import type { Dispatch, SetStateAction } from 'react'
import type { MessageData } from '@/models/MessageData'

export function openWebsocket(
    setMessages: Dispatch<SetStateAction<MessageData[]>>
): WebSocket {
    const wsProtocol = window.location.protocol == 'https:' ? 'wss' : 'ws'
    const wsURI = `${wsProtocol}://${window.location.host}/api/ws`
    const ws = new WebSocket(wsURI)
    ws.addEventListener('message', (e: MessageEvent) => {
        const message = JSON.parse(e.data) as MessageData
        setMessages((prevMessages: MessageData[]) => [...prevMessages, message])
    })
    return ws
}
