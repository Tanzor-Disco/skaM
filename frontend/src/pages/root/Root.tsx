import { useNavigate } from 'react-router'
import { useEffect } from 'react'

export default function Root() {
    const navigate = useNavigate()
    useEffect(() => {
        async function redirect() {
            const response = await fetch('api/root')
            if (!response.ok) {
                navigate('/login')
                return
            }
            navigate('/main')
        }
        redirect()
    })
    return null
}
