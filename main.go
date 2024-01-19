package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Bağlantı açıldığında bir mesaj gönder
	resourceID := "12345" // Burada gerçek resource ID'yi ayarlayın
	err = conn.WriteMessage(websocket.TextMessage, []byte(resourceID))
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		// Mesajları dinle
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Received message: %s\n", p)

		// İsteğe bağlı olarak burada gelen mesajlara göre işlemler yapabilirsiniz
		// Örneğin, başka bir client'a mesaj göndermek gibi
		_ = messageType
	}
}

func main() {
	http.HandleFunc("/", handleConnection)
	fmt.Println("WebSocket server started on :8080")
	http.ListenAndServe(":8888", nil)
}
