export interface Room {
    ID: number
    Name: string
}

export function newRoom(ID: number, name: string): Room {
    return {
        ID: ID,
        Name: name,
    }
}
