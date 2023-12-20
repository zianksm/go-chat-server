package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	app := gin.Default()
	app.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			panic(err)
		}

		for {
			messageType, p, err := conn.ReadMessage()

			if err != nil {
				panic(err)
			}

			conn.WriteMessage(messageType, p)
		}
	})
}
