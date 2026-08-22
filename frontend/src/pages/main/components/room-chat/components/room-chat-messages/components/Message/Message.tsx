import './Message.css'

interface MessageProps {
    username: string
    messageText: string
}

export default function Message({ username, messageText }: MessageProps) {
    return (
        <div className="message">
            <h1 className="username">{username}</h1>
            <p className="message-text">{messageText}</p>
        </div>
    )
}
