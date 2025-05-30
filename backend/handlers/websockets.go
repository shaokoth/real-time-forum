package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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
		return true
	},
}

// HandleWebSocket handles WebSocket connections
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Not logged in", http.StatusUnauthorized)
		return
	}
	_, err = utils.GetUserFromSession(cookie.Value)
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
	client := &Client{UserID: user.UUID, Conn: conn, Send: make(chan []byte)}

	mu.Lock()
	Clients[conn.RemoteAddr().String()] = client
	mu.Unlock()

	// Notify all other clients that a new user is online
	broadcast <- models.Message{
		Type: "user_status",
		Content: map[string]interface{}{
			"userId":   user.ID,
			"nickname": user.Nickname,
			"status":   "online",
		},
		CreatedAt: time.Now(),
	}

	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()
}

// Reads messages from the client
func (c *Client) readPump() {
	defer func() {
		// Get user details
		var nickname string
		err := database.Db.QueryRow("SELECT nickname FROM users WHERE id = ?", c.UserID).Scan(&nickname)
		if err != nil {
			// Notify all clients that a user is offline
			broadcast <- models.Message{
				Type: "user_status",
				Content: map[string]interface{}{
					"userId":   c.UserID,
					"nickname": nickname,
					"status":   "offline",
				},
				CreatedAt: time.Now(),
			}
		}
		c.Conn.Close()
	}()
	for {
		// Read message from the client
		_, message, err := c.Conn.ReadMessage() // ignore client messages
		if err != nil {
			break
		}
		// Parse message and determine receiver
		var msg models.Message
		if err = json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}
		// Set the sender ID
		msg.SenderID = c.UserID
		msg.CreatedAt = time.Now()

		// Set the sender ID
		msg.SenderID = c.UserID
		msg.CreatedAt = time.Now()

		// Handle the message based on its type
		switch msg.Type {
		case "private_message":
			// Store the message in the database
			_, err := database.Db.Exec(
				"INSERT INTO private_messages (content, sender_id, receiver_id) VALUES (?, ?, ?)",
				msg.Content.(map[string]interface{})["text"], msg.SenderID, msg.ReceiverID,
			)
			if err != nil {
				log.Printf("Error storing message: %v", err)
				continue
			}

			// Forward the message to the receiver
			broadcast <- msg

		case "fetch_messages":
			// Get previous messages between users
			receiverIDI := int(msg.Content.(map[string]interface{})["receiverId"].(float64))
			receiverID := strconv.Itoa(receiverIDI)
			limit := int(msg.Content.(map[string]interface{})["limit"].(float64))
			offset := int(msg.Content.(map[string]interface{})["offset"].(float64))

			// Query for messages
			rows, err := database.Db.Query(`
				SELECT id, content, sender_id, receiver_id, created_at, is_read
				FROM private_messages
				WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)
				ORDER BY created_at DESC
				LIMIT ? OFFSET ?
			`, c.UserID, msg.SenderID, receiverID, c.UserID, limit, offset)
			if err != nil {
				log.Printf("Error fetching messages: %v", err)
				continue
			}

			// Parse the messages
			var messages []map[string]interface{}
			for rows.Next() {
				var msg struct {
					ID         int
					Content    string
					SenderID   string
					ReceiverID string
					CreatedAt  time.Time
					IsRead     bool
				}
				err := rows.Scan(&msg.ID, &msg.Content, &msg.SenderID, &msg.ReceiverID, &msg.CreatedAt, &msg.IsRead)
				if err != nil {
					log.Printf("Error parsing message: %v", err)
					continue
				}

				// Add to the result
				messages = append(messages, map[string]interface{}{
					"id":         msg.ID,
					"content":    msg.Content,
					"senderId":   msg.SenderID,
					"receiverId": msg.ReceiverID,
					"createdAt":  msg.CreatedAt,
					"isRead":     msg.IsRead,
				})
			}
			rows.Close()

			// Mark messages as read
			_, err = database.Db.Exec(`
				UPDATE private_messages
				SET is_read = TRUE
				WHERE receiver_id = ? AND sender_id = ? AND is_read = FALSE
			`, c.UserID, receiverID)
			if err != nil {
				log.Printf("Error marking messages as read: %v", err)
			}

			// Send the messages back to the client
			response := models.Message{
				Type:       "messages_history",
				Content:    messages,
				ReceiverID: c.UserID,
				SenderID:   receiverID,
				CreatedAt:  time.Now(),
			}

			// Marshal and send
			responseJSON, err := json.Marshal(response)
			if err != nil {
				log.Printf("Error marshaling response: %v", err)
				continue
			}

			c.Send <- responseJSON
		}

		// Don't broadcast private messages to all clients
		if msg.Type != "private_message" && msg.Type != "fetch_messages" {
			broadcast <- msg
		}
	}
}

// Writes messages to the client
func (c *Client) writePump() {
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
