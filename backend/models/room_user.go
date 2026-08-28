package models

type RoomUserID int64

// RoomUser is used to describe a row in room_users table
// it can be used both for registering and retrieving data
type RoomUser struct {
	ID     RoomUserID
	RoomID RoomID
	UserID UserID
}

func NewRoomUser(roomID RoomID, userID UserID) RoomUser {
	return RoomUser{
		RoomID: roomID,
		UserID: userID,
	}
}
