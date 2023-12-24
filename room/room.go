package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	Id         uuid.UUID
	connection *websocket.Conn
}

type Room struct {
	Users    map[uuid.UUID]*User
	Id       uuid.UUID
	chatPipe chan string
}

var rooms map[uuid.UUID]*Room
var users map[uuid.UUID]*User

type UserActions func(map[uuid.UUID]*User)
type RoomActions func(map[uuid.UUID]*User)

func generateId() uuid.UUID {
	return uuid.New()
}
func IdTryFromString(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

func CreateRoom() *Room {
	room := new(Room)
	room.Id = generateId()
	room.Users = make(map[uuid.UUID]*User)
	room.chatPipe = make(chan string)

	rooms[room.Id] = room

	sendChat(room)

	return room
}

func CreateRoomWithId(id uuid.UUID) *Room {
	room := CreateRoom()
	room.Id = id

	return room
}

func sendChat(room *Room) {
	go func() {
		for {
			message := <-room.chatPipe
			for _, user := range room.Users {
				user.connection.WriteMessage(1, []byte(message))
			}
		}
	}()

}

func AddUser(conn *websocket.Conn, roomId uuid.UUID) *User {
	user := new(User)
	user.Id = generateId()
	user.connection = conn

	users[user.Id] = user

	room, exist := rooms[roomId]

	if !exist {
		room = CreateRoomWithId(roomId)
	}

	room.Users[user.Id] = user

	listenChatFromUser(user, room)
	sendJoinMessage(room, user)

	return user
}

func sendJoinMessage(room *Room, user *User) {
	message := user.Id.String() + " joined the room"
	sendSystemChat(room, message)
}

func sendSystemChat(room *Room, message string) {
	msg := "System: " + message
	room.chatPipe <- msg

}

func listenChatFromUser(user *User, room *Room) {
	go func() {
		for {
			msgType, msg, err := user.connection.ReadMessage()
			message := string(msg)
			message = user.Id.String() + ": " + message

			if websocket.IsCloseError(err) {
				removeUser(user, room)
				return
			}

			if err != nil {
				fmt.Println("Error reading message: ", err)
				continue
			}

			if msgType != websocket.TextMessage {
				fmt.Println("Not a text message")
				continue
			}

			room.chatPipe <- message
		}
	}()

}

func removeUser(user *User, room *Room) {
	delete(room.Users, user.Id)
	delete(users, user.Id)
}

func InitRooms() {
	rooms = make(map[uuid.UUID]*Room)
	users = make(map[uuid.UUID]*User)
}
