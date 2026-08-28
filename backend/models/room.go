package models

import ()

type RoomID int64

// Room represents a row in rooms table
// It can be used both for inserting or retrieving info from rooms table
type Room struct {
	ID   RoomID
	Name string
}

func NewRoom(name string) Room {
	return Room{
		Name: name,
	}
}
