import './RegisterForm.css'
import FormField from '@/components/form-field/FormField'
import { useState } from 'react'
import type { Dispatch, SetStateAction, SyntheticEvent } from 'react'

import { emailValid, usernameValid, passwordValid } from '@/utils/validate'

interface RegisterFormProps {
    setRegisterComplete: Dispatch<SetStateAction<boolean>>
}

interface RegisterSubmit {
    email: string
    username: string
    password: string
}

export default function RegisterForm({
    setRegisterComplete,
}: RegisterFormProps) {
    const [emailErr, setEmailErr] = useState('')
    const [usernameErr, setUsernameErr] = useState('')
    const [passwordErr, setPasswordErr] = useState('')

    function toggleState(errorKind: string) {
        switch (errorKind) {
            case 'ERR_EMAIL_TAKEN':
                setEmailErr('Email is already taken')
                break
            case 'ERR_INVALID_USERNAME_LENGTH':
                setUsernameErr(
                    'The maximum length of the username is 20 characters'
                )
                break
            case 'ERR_FORBIDDEN_PASSWORD_CHARS':
                setPasswordErr(
                    'The password should contain only latin characters and digits'
                )
                break
            case 'ERR_INVALID_PASSWORD_LENGTH':
                setPasswordErr(
                    'The maximum length of the password is 70 characters'
                )
        }
    }
    function toggleAllOff() {
        setEmailErr('')
        setUsernameErr('')
        setPasswordErr('')
    }

    function registerSubmitValid(submit: RegisterSubmit): boolean {
        toggleAllOff()
        if (!emailValid(submit.email)) {
            setEmailErr(
                'Invalid Email: the site accepts proton, tuta and gmail'
            )
            return false
        }
        if (!usernameValid(submit.username)) {
            setUsernameErr(
                'The maximum length should be no more than 20 characters'
            )
            return false
        }
        if (!passwordValid(submit.password)) {
            setPasswordErr(
                'Password should be no more than 70 characters and consist only of latin character and numbers'
            )
            return false
        }
        return true
    }

    async function handleSubmit(event: SyntheticEvent<HTMLFormElement>) {
        toggleAllOff()
        event.preventDefault()
        const form = event.currentTarget
        const formData = new FormData(form)
        const registerSubmit: RegisterSubmit = {
            email: formData.get('email') as string,
            username: formData.get('username') as string,
            password: formData.get('password') as string,
        }
        form.reset()
        if (!registerSubmitValid(registerSubmit)) {
            return
        }
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(registerSubmit),
        })

        if (!response.ok) {
            const responseObj = await response.json()
            toggleState(responseObj.error_kind)
            return
        }
        setRegisterComplete(true)
    }

    return (
        <form className="register-form" onSubmit={handleSubmit}>
            <FormField
                warning={emailErr}
                fieldName={'Email'}
                inputType={'email'}
            />
            <FormField
                warning={usernameErr}
                fieldName={'Username'}
                inputType={'text'}
            />
            <FormField
                warning={passwordErr}
                fieldName={'Password'}
                inputType={'text'}
            />
            <button className="form-submit">Submit</button>
        </form>
    )
}
