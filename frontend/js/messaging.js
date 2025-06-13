document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const chatWindow = document.getElementById("chatWindow");
  const chatHeader = document.getElementById("chatHeader");
  const messagesContainer = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");
  const currentUserID = localStorage.getItem("CurrentUserID");
  const chatAvatar = document.getElementById("chatAvatar");
  const chatUserName = document.getElementById("chatUserName");
  const chatUserStatus = document.getElementById("chatUserStatus");
  const typingIndicator = document.getElementById("typingIndicator");

  let socket;
  let currentReceiver = null;
  let currentReceiverName = null;
  let typingTimeout;
  let messageOffset = 0;
  let isLoadingMessages = false;

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
      } else if (message.type === "status") {
        console.log(message);
        updateUserStatus(message.sender_id, message.online);
      } else if (
        message.type === "typing" &&
        message.sender_id === currentReceiver
      ) {
        typingIndicator.textContent = "Typing...";
      } else if (
        message.type === "stop_typing" &&
        message.sender_id === currentReceiver
      ) {
        typingIndicator.textContent = "";
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

  function updateUserList(users) {
    if (!userList) return;
    userList.innerHTML = ""; // Clear previous entries
    users.forEach((user) => {
      createUserItem(user);
    });
  }

  function updateUserStatus(userId, isOnline) {
    const userItems = document.querySelectorAll(".user-item");
    userItems.forEach((item) => {
      if (item.dataset.userId === userId) {
        const statusDiv = item.querySelector(".user-status");
        if (statusDiv) {
          statusDiv.textContent = isOnline ? "Online" : "Offline";
          statusDiv.style.color = isOnline ? "green" : "gray";
        }
      }
    });
  }

  function createUserItem(user) {
    const userDiv = document.createElement("div");
    userDiv.className = "user-item";
    userDiv.dataset.userId = user.user_uuid; // camelCase used here
    console.log(user);

    const onlineStatus = user.isOnline ? "Online" : "Offline";
    const statusColor = user.isOnline ? "green" : "gray";

    userDiv.innerHTML = `
    <div class="user-avatar">${user.nickname.charAt(0).toUpperCase()}</div>
    <div class="user-info">
      <div class="user-name">${user.nickname}</div>
      <div class="user-status" style="color: ${statusColor}">${onlineStatus}</div>
    </div>
  `;

    userDiv.onclick = () => {
      document.getElementById("messagesTitle").style.display = "none";
      openChat(user.user_uuid, user.nickname, onlineStatus);
    };

    userList.appendChild(userDiv);
  }

  function openChat(userId, userName, status) {
    currentReceiver = userId;
    currentReceiverName = userName;
    messageOffset = 0;

    // Update chat header
    chatAvatar.textContent = userName.charAt(0).toUpperCase();
    chatUserName.textContent = userName;
    chatUserStatus.textContent = status;

    // Show chat window, hide user list
    userList.classList.add("hidden");
    chatWindow.classList.remove("hidden");

    // Load chat history
    loadChatHistory(currentReceiver, false);

    // Focus on input
    messageInput.focus();
  }

  if (userList) {
    fetch("/users")
      .then((response) => {
        if (!response.ok) throw new Error("Failed to fetch users");
        return response.json();
      })
      .then((users) => {
        updateUserList(users);
      })
      .catch((err) => {
        console.error("Error fetching users:", err);
      });
  }

  function loadChatHistory(otherUserID, append = false) {
    if (isLoadingMessages) return;

    isLoadingMessages = true;
    fetch(`/messages?with=${otherUserID}&offset=${messageOffset}`)
      .then((response) => {
        if (!response.ok) throw new Error("Failed to fetch messages");
        return response.json();
      })
      .then((messages) => {
        if (!Array.isArray(messages)) {
          console.error("Chat history is not an array:", messages);
          return;
        }

        if (!append) {
          messagesContainer.innerHTML = ""; // On initial open
        }

        const scrollPositionBefore = messagesContainer.scrollHeight;

        messages.reverse().forEach((message) => {
          const messageElement = createMessageElement(message);
          if (append) {
            messagesContainer.prepend(messageElement);
          } else {
            messagesContainer.appendChild(messageElement);
          }
        });

        // Adjust scroll position to maintain view if loading older messages
        if (append) {
          messagesContainer.scrollTop =
            messagesContainer.scrollHeight - scrollPositionBefore;
        } else {
          messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }

        messageOffset += 10;
      })
      .catch((error) => {
        console.error("Error loading chat history:", error);
      })
      .finally(() => {
        isLoadingMessages = false;
      });
  }

  function createMessageElement(message) {
    const messageElement = document.createElement("div");
    messageElement.className = `message ${
      message.sender_id === currentUserID ? "sent" : "received"
    }`;

    const messageContent = document.createElement("div");
    messageContent.className = "message-content";
    messageContent.textContent = message.content;

    const messageTime = document.createElement("div");
    messageTime.className = "message-time";
    messageTime.textContent = new Date(message.created_at).toLocaleTimeString();
    messageTime.style.fontSize = "0.75rem";
    messageTime.style.opacity = "0.7";
    messageTime.style.marginTop = "4px";

    messageElement.appendChild(messageContent);
    messageElement.appendChild(messageTime);

    return messageElement;
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

  messageInput.addEventListener("input", () => {
    //Send "typing" message
    socket.send(
      JSON.stringify({
        type: "typing",
        receiver_id: currentReceiver,
      })
    );
    // Clear existing timeout
    clearTimeout(typingTimeout);

    // Send "stop_typing" after 1 second of inactivity
    typingTimeout = setTimeout(() => {
      socket.send(
        JSON.stringify({
          type: "stop_typing",
          receiver_id: currentReceiver,
        })
      );
    }, 5000);
  });

  function throttle(func, limit) {
    let inThrottle;
    return function () {
      if (!inThrottle) {
        func.apply(this, arguments);
        inThrottle = true;
        setTimeout(() => (inThrottle = false), limit);
      }
    };
  }

  messagesContainer.addEventListener(
    "scroll",
    throttle(() => {
      if (messagesContainer.scrollTop < 50 && currentReceiver) {
        loadChatHistory(currentReceiver, true);
      }
    }, 1000) // 1 second throttle
  );

  messageInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      socket.send(
        JSON.stringify({
          type: "stop_typing",
          receiver_id: currentReceiver,
        })
      );
    }
  });

  sendButton.addEventListener("click", sendMessage);
  messageInput.addEventListener("keypress", function (e) {
    if (e.key === "Enter") sendMessage();
  });

  chatHeader.addEventListener("click", () => {
    chatWindow.classList.add("hidden");
    userList.classList.remove("hidden");
    currentReceiver = null;
    currentReceiverName = null;
    document.getElementById("messagesTitle").style.display = "block";
  });
  // Initialize
  connectWebSocket();
});
