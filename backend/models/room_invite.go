package models

type RoomInviteID int64

type RoomInvite struct {
	ID     RoomInviteID
	RoomID RoomID
	Token  string
}

func NewRoomInvite(roomID RoomID, token string) RoomInvite {
	return RoomInvite{
		RoomID: roomID,
		Token:  token,
	}
}
