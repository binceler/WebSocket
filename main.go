package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ReceivedMessage struct {
	Action    string `json:"action"`
	ThisAgent string `json:"thisAgent"`
	SessionID string `json:"session_id"`
	IsAdmin   bool   `json:"is_admin"`
}

var onlineAgentSessions = make(map[string]string)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	connectionID := uuid.New().String()

	mapD := map[string]string{"action": "sysMsg", "content": "Welcome: " + connectionID}
	mapB, _ := json.Marshal(mapD)

	// Bağlantı açıldığında bir mesaj gönder
	resourceID := string(mapB) // Burada gerçek resource ID'yi ayarlayın
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
		var receivedMessage ReceivedMessage
		err = json.Unmarshal(p, &receivedMessage)

		if receivedMessage.Action == "checkAgentOnlineList" {
			if _, ok := onlineAgentSessions[receivedMessage.ThisAgent]; ok {
				if onlineAgentSessions[receivedMessage.ThisAgent] != receivedMessage.SessionID {
					messageMap, _ := json.Marshal(map[string]string{"action": "logOutAgent"})
					message := string(messageMap)
					err = conn.WriteMessage(websocket.TextMessage, []byte(message))
				}
			} else {
				onlineAgentSessions[receivedMessage.ThisAgent] = receivedMessage.SessionID
			}
		}

		if err != nil {
			fmt.Println(err)
			return
		}

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
