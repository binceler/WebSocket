package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	db *gorm.DB
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

	/*var result string
	err = db.Raw("UPDATE tbl_users SET first_name='abcdefgo' where id=1").Scan(&result).Error*/

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
	var err error
	db, err = Connect()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", handleConnection)
	fmt.Println("WebSocket server started on :8080")
	http.ListenAndServe(":8888", nil)
}

func Connect() (*gorm.DB, error) {
	// PostgreSQL bağlantı bilgileri
	var dsn string = "host=host.docker.internal user=busrai password=123 dbname=laravel_last port=5432"

	// PostgreSQL'e bağlan
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
