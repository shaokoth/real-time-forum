package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"real-time-forum/backend/database"
	"real-time-forum/backend/models"
	"real-time-forum/backend/utils"

	"github.com/gorilla/websocket"
)

var (
	user      models.User
	Clients   map[string]*Client
	broadcast chan models.Message
)

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

// InitWebSocket initializes the WebSocket system
func InitWebSocket() {
	go handleBroadcast()
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	user, err := utils.GetUserFromSession(cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}
	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connections", err)
		return
	}
	client := &Client{
		UserID: user.UUID,
		Conn:   conn,
		Send:   make(chan []byte),
	}

	// Update user status to online
	_, err = database.Db.Exec("UPDATE users SET is_online = TRUE, last_seen = CURRENT_TIMESTAMP WHERE id = ?", client.UserID)
	if err != nil {
		log.Println(err)
	}

	mu.Lock()
	Clients[user.UUID] = client
	mu.Unlock()

	

	go readMessages(client)
	go writeMessages(client)
}

// Reads messages from the client
func readMessages(c *Client) {
	defer func() {
		mu.Lock()

		// Update user status to offline
		_, err := database.Db.Exec("UPDATE users SET is_online = FALSE, last_seen = CURRENT_TIMESTAMP WHERE id = ?", c.UserID)
		if err != nil {
			log.Println(err)
		}
		delete(Clients, c.UserID)
		mu.Unlock()
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break // Disconnect client
		}
		var msg models.Message
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Validate message
		if msg.ReceiverID == "" || msg.Content == "" {
			log.Println("Invalid message format")
			continue
		}

		msg.SenderID = c.UserID
		msg.CreatedAt = time.Now()

		_, err = database.Db.Exec(
			"INSERT INTO private_messages (content, sender_id, receiver_id, created_at) VALUES (?, ?, ?, ?)",
			msg.Content, msg.SenderID, msg.ReceiverID, msg.CreatedAt,
		)
		if err != nil {
			log.Printf("Error storing message: %v", err)
			continue
		}
		// Forward the message to the receiver
		broadcast <- msg
	}
}

func writeMessages(c *Client) {
	defer func() {
		c.Conn.Close()
	}()
	for {
		message, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		// Write the message to the client
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// handleBroadcast routes messages from the broadcast channel to the appropriate clients
func handleBroadcast() {
	for {
		msg := <-broadcast
		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Broadcast marshal error: %v", err)
			continue
		}

		mu.Lock()
		// Send to receiver
		if receiver, ok := Clients[msg.ReceiverID]; ok {
			receiver.Send <- data
		}
		// Also echo back to sender if connected
		if sender, ok := Clients[msg.SenderID]; ok {
			sender.Send <- data
		}
		mu.Unlock()
	}
}
