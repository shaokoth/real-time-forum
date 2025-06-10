document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const chatWindow = document.getElementById("chatWindow");
  const chatHeader = document.getElementById("chatHeader");
  const messagesContainer = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");
  const currentUserID = localStorage.getItem("CurrentUserID");

  let socket;
  let currentReceiver = null;
  // WebSocket functions
  function connectWebSocket() {
    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const wsUrl = protocol + window.location.host + "/ws";

    socket = new WebSocket(wsUrl);

    socket.onopen = function (event) {
      console.log("WebSocket connected");
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

  if (!userList) return;

  fetch("/users")
    .then((response) => {
      if (!response.ok) throw new Error("Failed to fetch users");
      return response.json();
    })
    .then((users) => {
      userList.innerHTML = ""; // Clear previous entries
      users.forEach((user) => {
        const userDiv = document.createElement("div");
        userDiv.className = "user-item";
        userDiv.textContent = user.nickname;
        userDiv.onclick = () => {
          currentReceiver = user.user_uuid; // Save the UUID of the clicked user
          chatHeader.textContent = `Chat with ${user.nickname}`; // Set chat header
          chatWindow.classList.remove("hidden"); // Show the chat window
          loadChatHistory(currentReceiver);
        };
        userList.appendChild(userDiv);
      });
    })
    .catch((err) => {
      console.error("Error fetching users:", err);
    });

  function loadChatHistory(otherUserID) {
    fetch(`/messages?with=${otherUserID}&offset=10`)
      .then((response) => {
        if (!response.ok) throw new Error("Failed to fetch messages");
        return response.json();
      })
      .then((messages) => {
        messagesContainer.innerHTML = ""; // Clear previous messages

        if (!Array.isArray(messages)) {
          console.error("Chat history is not an array:", messages);
          return;
        }
        console.log(messages)
        messages.forEach((message) => {
          displayMessage(message);
        });
      })
      .catch((error) => {
        console.error("Error loading chat history:", error);
      });
  }

  function displayMessage(message) {
    if (
      message.sender_id !== currentReceiver &&
      message.receiver_id !== currentReceiver
    ) {
      return; // Not part of current chat
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
    messagesContainer.scrollTop = messagesContainer.scrollHeight; // Scroll to bottom
  }

  function sendMessage() {
    const content = messageInput.value.trim();
    messageInput.value = "";
    if (!content || !currentReceiver) return;

    const message = {
      type: "message", // Optional: helpful for distinguishing message types
      sender_id: currentUserID, // <-- You must make sure currentUserID is set
      receiver_id: currentReceiver,
      content: content,
    };

    // Send via WebSocket
    socket.send(JSON.stringify(message));

    // Clear input
    messageInput.value = "";

    // Optimistically show the message (assumes it's sent)
    displayMessage({
      sender_id: currentUserID,
      receiver_id: currentReceiver,
      content: content,
      created_at: new Date().toISOString(),
    });
  }

  sendButton.addEventListener("click", sendMessage);
  messageInput.addEventListener("keypress", function (e) {
    if (e.key === "Enter") sendMessage();
  });
  // Initialize
  connectWebSocket();
});
