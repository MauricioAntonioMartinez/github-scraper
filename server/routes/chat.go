package routes

import (
	"encoding/json"
	"fmt"

	"github.com/antoniodipinto/ikisocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type MessageObject struct {
	Data string `json:"data"`
	From string `json:"from"`
	To   string `json:"to"`
}



// ChatRoutes handling websockets
func ChatRoutes(app *fiber.App){
	clients := make(map[string]string)
	router := app.Group("/ws")

	router.Use( func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			c.Locals("user_id", c.Query("user_id"))
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})



	router.Get("/client", ikisocket.New(func(kws *ikisocket.Websocket) {
		// Retrieve user id from the middleware (optional)
		userID := fmt.Sprintf("%v", kws.Locals("user_id"))

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("user_id", userID)

		// On connect event. Notify when comes a new connection
		kws.OnConnect = func() {
			// Add the connection to the list of the connected clients
			// The UUID is generated randomly
			clients[userID] = kws.UUID
			//Broadcast to all the connected users the newcomer
			kws.Broadcast([]byte("New user connected: "+userID+" and UUID: "+kws.UUID), true)
			//Write welcome message
			kws.Emit([]byte("Hello user: " + userID + " and UUID: " + kws.UUID))
		}

		// On message event
		kws.OnMessage = func(data []byte) {

			message := MessageObject{}
			json.Unmarshal(data, &message)
			fmt.Println(string(data))
			// Emit the message directly to specified user
			err := kws.EmitTo(clients[message.To], data)
			if err != nil {
				fmt.Println(err)
			}
		}
	}))

	// router.Get("/:id", websocket.New(func(c *websocket.Conn) {
	// 	// c.Locals is added to the *websocket.Conn
	// 	log.Println(c.Locals("allowed"))  // true
	// 	log.Println(c.Params("id"))       // 123
	// 	log.Println(c.Query("v"))         // 1.0
	// 	log.Println(c.Cookies("session")) // ""

	// 	// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
	// 	var (
	// 		mt  int
	// 		msg []byte
	// 		err error
	// 	)
		

	// 	for {
	// 		if mt, msg, err = c.ReadMessage(); err != nil {
	// 			log.Println("read:", err)
	// 			break
	// 		}
	// 		log.Printf("recv: %s", msg)

	// 		if err = c.WriteMessage(mt, append( msg,23,43,34)); err != nil {
	// 			log.Println("write:", err)
	// 			break
	// 		}
	// 	}

	// }))

}