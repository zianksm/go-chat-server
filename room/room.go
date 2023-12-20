package room

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	Id         uuid.UUID
	connection websocket.Conn
}

type Room struct {
	Users map[uuid.UUID]*User
	Id    uuid.UUID
}

var rooms map[uuid.UUID]*Room
var users map[uuid.UUID]*User

type UserActions func(map[uuid.UUID]*User)
type RoomActions func(map[uuid.UUID]*User)

func GenerateId() uuid.UUID {
	return uuid.New()
}

func InitRooms() {
	rooms = make(map[uuid.UUID]*Room)
	users = make(map[uuid.UUID]*User)
}

func DoUsersInRoom(callback UserActions, roomId uuid.UUID) {
	users := GetUsersInRoom(roomId)
	callback(users)
}

func CreateRoom() *Room {
	room := new(Room)
	room.Id = GenerateId()
	rooms[room.Id] = room

	return room
}

func GetRoom(id uuid.UUID) *Room {
	return rooms[id]
}

func GetUsersInRoom(id uuid.UUID) map[uuid.UUID]*User {
	return users
}

func CreateUser(conn websocket.Conn) *User {
	user := new(User)
	user.Id = GenerateId()
	user.connection = conn
	users[user.Id] = user

	return user
}

func AddUserToRoom(roomId uuid.UUID, user *User) {
}
