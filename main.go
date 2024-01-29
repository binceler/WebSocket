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
	"os"
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
	Action         string `json:"action"`
	ThisAgent      string `json:"thisAgent"`
	SessionID      string `json:"session_id"`
	IsAdmin        bool   `json:"is_admin"`
	Location       string `json:"location"`
	Username       string `json:"username"`
	SignLangStatus string `json:"sign_lang_status"`
	AgentName      string `json:"agent_name"`
	Browser        string `json:"browser"`
	UserID         string `json:"user_id"`
	AvailableTime  int64  `json:"-"`
	AssignStatus   bool   `json:"assign_status"`
	ProjectID      string `json:"project_id"`
	Room           string `json:"room"`
}

var onlineAgentSessions = make(map[string]string)
var rooms = make(map[string]interface{})
var queue = make(map[string]interface{})
var queueListeners = make(map[string]interface{})
var admins = make(map[string]interface{})
var deviceInfo = make(map[string]interface{})
var assignReqs = make(map[string]interface{})
var interactionIds = make(map[string]interface{})
var identAssignCounts = make(map[string]interface{})

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
		var data map[string]interface{}
		from := &ReceivedMessage{}

		// Gelen JSON verisinden alanları okuma ve bağlantı arayüzüne atama
		from.Location = getString(data, "location")
		from.Username = getString(data, "username")
		from.SignLangStatus = getString(data, "sign_lang_status")
		from.AgentName = getString(data, "agent_name")
		from.Browser = getString(data, "browser")
		from.SessionID = getString(data, "session_id")
		from.UserID = getString(data, "user_id")
		from.ProjectID = getString(data, "project_id")
		from.Room = getString(data, "room")

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

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func Connect() (*gorm.DB, error) {

	dbHost := os.Getenv("CONFIG_PHP_DBHOST")
	dbUser := os.Getenv("CONFIG_PHP_DBUSER")
	dbPass := os.Getenv("CONFIG_PHP_DBPASS")
	dbPort := os.Getenv("CONFIG_PHP_DBPORT")
	dbName := os.Getenv("CONFIG_PHP_DBNAME")

	// Construct the DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbHost, dbUser, dbPass, dbName, dbPort)

	// Use the DSN for your PostgreSQL connection

	// PostgreSQL'e bağlan
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
