document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const chatWindow = document.getElementById("chatWindow");
  const chatHeader = document.getElementById("chatHeader");
  const messagesContainer = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");

  let socket;
  let currentReceiver = null;
  // WebSocket functions
  function connectWebSocket() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = (protocol + window.location.host + '/ws');

    socket = new WebSocket(wsUrl);

    socket.onopen = function (event) {
      console.log("WebSocket connected");
      fetchOnlineUsers();
    };

    socket.onmessage = function (event) {
      const message = JSON.parse(event.data);
      if (message.type === "user_list") {
        updateUserList(message.users);
      } else if (message.type === "message") {
        displayMessage(message);
      }
    };

    socket.onclose = function (event) {
      console.log("WebSocket disconnected");
      setTimeout(connectWebSocket, 3000); // Reconnect after 3 seconds
    };

    socket.onerror = function (error) {
      console.error("WebSocket error:", error);
    };
  }
  // Fetch online users
  function fetchOnlineUsers() {
    socket.send(JSON.stringify({ type: "get_users" }));
  }

  // Update user list in the UI
  function updateUserList(users) {
    userList.innerHTML = "";
    users.forEach((user) => {
      const userElement = document.createElement("div");
      userElement.className = "user-item";
      userElement.textContent = user.Nickname;
      userElement.dataset.userId = user.UUID;

      userElement.addEventListener("click", () => {
        currentReceiver = user.UUID;
        chatHeader.textContent = `Chat with ${user.Nickname}`;
        chatWindow.classList.remove("hidden");
        loadChatHistory(user.UUID);
      });

      userList.appendChild(userElement);
    });
  }

  // Load chat history
  function loadChatHistory(receiverId) {
    fetch(`/messages${receiverId}`)
      .then((response) => response.json())
      .then((messages) => {
        messagesContainer.innerHTML = "";
        messages.forEach((message) => {
          displayMessage(message);
        });
      });
  }

  // Display a message in the chat window
  function displayMessage(message) {
    if (
      message.sender_id !== currentReceiver &&
      message.receiver_id !== currentReceiver
    ) {
      return; // Not relevant to current chat
    }

    const messageElement = document.createElement("div");
    messageElement.className =
      message.sender_id === currentReceiver ? "received" : "sent";
    messageElement.innerHTML = `
            <div class="message-content">${message.content}</div>
            <div class="message-time">${new Date(
              message.created_at
            ).toLocaleTimeString()}</div>
        `;
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  // Send message
  function sendMessage() {
    const content = messageInput.value.trim();
    if (!content || !currentReceiver) return;

    const message = {
      type: "message",
      receiver_id: currentReceiver,
      content: content,
    };

    socket.send(JSON.stringify(message));
    messageInput.value = "";

    // Optimistically display the message
    displayMessage({
      sender_id: "me", // This will be replaced with actual ID from server
      receiver_id: currentReceiver,
      content: content,
      created_at: new Date().toISOString(),
    });
  }

  // Event listeners
  sendButton.addEventListener("click", sendMessage);
  messageInput.addEventListener("keypress", function (e) {
    if (e.key === "Enter") sendMessage();
  });

  // Initialize
    connectWebSocket();
});
