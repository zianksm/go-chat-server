package main

import (
	"chatserver/room"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	room.InitRooms()
	app := gin.Default()

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	app.GET("/ws/", func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {
			panic(err)
		}

		chatRoom := room.CreateRoom()

		room.AddUser(conn, chatRoom.Id)

	})

	app.GET("/ws/:id", func(ctx *gin.Context) {
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		roomId := ctx.Param("id")

		if err != nil {
			panic(err)
		}

		uuid, err := uuid.Parse(roomId)

		if err != nil {
			ctx.JSON(
				400,
				gin.H{
					"error": "Invalid room id",
				},
			)
		}

		chatRoom := room.CreateRoomWithId(uuid)

		room.AddUser(conn, chatRoom.Id)
	})

	app.Run(":8080")
}
