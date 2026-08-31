export function emailValid(email: string): Boolean {
    const re =
        /^[a-zA-Z0-9_-]+@(gmail\.com|protonmail\.com|proton\.me|tutamail\.com|tuta\.com|tutanota\.com|tutanota\.de)$/
    if (email.match(re) === null) {
        return false
    }
    if (email.length > 300) {
        return false
    }
    return true
}

export function usernameValid(username: string): Boolean {
    if (username.length > 0 && username.length <= 20) {
        return true
    }
    return false
}

function isLatin(char: string) {
    if ((char > 'a' && char < 'z') || (char > 'A' && char < 'Z')) {
        return true
    }
    return false
}

function isNumber(char: string) {
    if (char > '0' && char < '9') {
        return true
    }
    return false
}

export function passwordValid(password: string): Boolean {
    for (const char of password) {
        if (!isLatin(char) && !isNumber(char)) {
            return false
        }
        if (password.length == 0 || password.length > 70) {
            return false
        }
    }
    return true
}
