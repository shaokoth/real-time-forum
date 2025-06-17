document.addEventListener("DOMContentLoaded", function () {
  const userList = document.getElementById("userList");
  const chatWindow = document.getElementById("chatWindow");
  const chatHeader = document.getElementById("chatHeader");
  const messagesContainer = document.getElementById("messages");
  const messageInput = document.getElementById("messageInput");
  const sendButton = document.getElementById("sendButton");
  const currentUserID = localStorage.getItem("CurrentUserID");
  const currentUserName = localStorage.getItem("CurrentUserName");

  const chatAvatar = document.getElementById("chatAvatar");
  const chatUserName = document.getElementById("chatUserName");
  const chatUserStatus = document.getElementById("chatUserStatus");
  const typingIndicator = document.getElementById("typingIndicator");
  const emptyState = document.getElementById("emptyState");
  const messageInputContainer = document.getElementById(
    "messageInputContainer"
  );
  const searchInput = document.getElementById("searchInput");
  const backButton = document.getElementById("backButton");
  const usersSidebar = document.getElementById("usersSidebar");

  let socket;
  let currentReceiver = null;
  let currentReceiverName = null;
  let typingTimeout;
  let messageOffset = 0;
  let isLoadingMessages = false;
  let allUsers = [];

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
        showTypingIndicator();
      } else if (
        message.type === "stop_typing" &&
        message.sender_id === currentReceiver
      ) {
        hideTypingIndicator();
      }
    };

    socket.onclose = function (event) {
      console.log("WebSocket disconnected");
      setTimeout(connectWebSocket, 3000);
    };

    socket.onerror = function (error) {
      console.error("WebSocket error:", error);
    };
  }

  function updateUserList(users) {
    if (!userList) return;
    allUsers = users;
    renderFilteredUsers();
  }

  function renderFilteredUsers() {
    const searchTerm = searchInput.value.toLowerCase();
    const filteredUsers = allUsers.filter((user) =>
      user.nickname.toLowerCase().includes(searchTerm)
    );

    userList.innerHTML = "";
    filteredUsers.forEach((user) => {
      createUserItem(user);
    });
  }

  function updateUserStatus(userId, isOnline) {
    const userItems = document.querySelectorAll(".user-item");
    userItems.forEach((item) => {
      if (item.dataset.userId === userId) {
        const statusIndicator = item.querySelector(".status-indicator");
        const statusText = item.querySelector(".user-status");
        if (statusIndicator && statusText) {
          statusIndicator.className = `status-indicator ${
            isOnline ? "online" : "offline"
          }`;
          statusText.textContent = isOnline ? "Online" : "Offline";
        }
      }
    });

    // Update chat header if this is the current conversation
    if (currentReceiver === userId) {
      chatUserStatus.textContent = isOnline ? "Online" : "Offline";
      const chatStatusIndicator = chatAvatar.querySelector(".status-indicator");
      if (chatStatusIndicator) {
        chatStatusIndicator.className = `status-indicator ${
          isOnline ? "online" : "offline"
        }`;
      }
    }
  }

  function createUserItem(user) {
    const userDiv = document.createElement("div");
    userDiv.className = "user-item";
    userDiv.dataset.userId = user.user_uuid;

    const onlineStatus = user.isOnline ? "Online" : "Offline";
    const statusClass = user.isOnline ? "online" : "offline";

    userDiv.innerHTML = `
                    <div class="user-avatar">
                        ${user.nickname.charAt(0).toUpperCase()}
                        <div class="status-indicator ${statusClass}"></div>
                    </div>
                    <div class="user-info">
                        <div class="user-name">${user.nickname}</div>
                        <div class="user-status">${onlineStatus}</div>
                    </div>
                `;

    userDiv.onclick = () => {
      openChat(user.user_uuid, user.nickname, onlineStatus);
    };

    userList.appendChild(userDiv);
  }

  function openChat(userId, userName, status) {
    currentReceiver = userId;
    currentReceiverName = userName;
    messageOffset = 0;

    // Update active user in sidebar
    document.querySelectorAll(".user-item").forEach((item) => {
      item.classList.remove("active");
    });
    document
      .querySelector(`[data-user-id="${userId}"]`)
      .classList.add("active");

    // Update chat header
    chatAvatar.innerHTML = `
                    ${userName.charAt(0).toUpperCase()}
                    <div class="status-indicator ${
                      status === "Online" ? "online" : "offline"
                    }"></div>
                `;
    chatUserName.textContent = userName;
    chatUserStatus.textContent = status;

    // Show chat interface, hide empty state
    emptyState.style.display = "none";
    chatHeader.style.display = "flex";
    messagesContainer.style.display = "flex";
    messageInputContainer.style.display = "block";

    // On mobile, hide sidebar and show chat
    if (window.innerWidth <= 768) {
      usersSidebar.classList.add("hidden");
    }

    // Load chat history
    loadChatHistory(currentReceiver, false);

    // Focus on input
    messageInput.focus();
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
          messagesContainer.innerHTML = "";
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

    const messageUsername = document.createElement("div");
    messageUsername.className = "message-username";
    messageUsername.textContent =
      message.sender_id === currentUserID ? "you" : currentReceiverName;

    const messageContent = document.createElement("div");
    messageContent.className = "message-content";
    messageContent.textContent = message.content;

    const messageTime = document.createElement("div");
    messageTime.className = "message-time";
    messageTime.textContent = new Date(message.created_at).toLocaleTimeString(
      [],
      {
        hour: "2-digit",
        minute: "2-digit",
      }
    );

    messageContent.appendChild(messageTime);
    messageElement.appendChild(messageContent);
    messageElement.appendChild(messageUsername);

    return messageElement;
  }

  function displayMessage(message) {
    if (
      message.sender_id === currentReceiver ||
      message.receiver_id === currentReceiver ||
      message.sender_id === currentUserID
    ) {
      const messageElement = createMessageElement(message);
      messagesContainer.appendChild(messageElement);
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }

  function sendMessage() {
    const content = messageInput.value.trim();
    if (!content || !currentReceiver) return;

    const message = {
      type: "message",
      sender_id: currentUserID,
      receiver_id: currentReceiver,
      content: content,
    };

    socket.send(JSON.stringify(message));
    messageInput.value = "";

    // Stop typing indicator
    socket.send(
      JSON.stringify({
        type: "stop_typing",
        receiver_id: currentReceiver,
      })
    );

    // Optimistically show the message
    displayMessage({
      sender_id: currentUserID,
      receiver_id: currentReceiver,
      content: content,
      created_at: new Date().toISOString(),
    });
  }

  function showTypingIndicator() {
    const existingIndicator = document.querySelector(".typing-indicator");
    if (existingIndicator) return;

    const typingDiv = document.createElement("div");
    typingDiv.className = "typing-indicator";
    typingDiv.innerHTML = `
                    <div class="typing-dots">
                        <div class="typing-dot"></div>
                        <div class="typing-dot"></div>
                        <div class="typing-dot"></div>
                    </div>
                    <span style="font-size: 0.75rem; color: var(--gray-500);">Typing...</span>
                `;

    messagesContainer.appendChild(typingDiv);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
  }

  function hideTypingIndicator() {
    const typingIndicator = document.querySelector(".typing-indicator");
    if (typingIndicator) {
      typingIndicator.remove();
    }
  }

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

  // Event Listeners
  searchInput.addEventListener("input", renderFilteredUsers);

  messageInput.addEventListener("input", () => {
    if (!currentReceiver) return;

    socket.send(
      JSON.stringify({
        type: "typing",
        receiver_id: currentReceiver,
      })
    );

    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
      socket.send(
        JSON.stringify({
          type: "stop_typing",
          receiver_id: currentReceiver,
        })
      );
    }, 1000);
  });

  messageInput.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      sendMessage();
    }
  });

  sendButton.addEventListener("click", sendMessage);

  backButton.addEventListener("click", () => {
    usersSidebar.classList.remove("hidden");
    emptyState.style.display = "flex";
    chatHeader.style.display = "none";
    messagesContainer.style.display = "none";
    messageInputContainer.style.display = "none";
    currentReceiver = null;
    currentReceiverName = null;
  });

  messagesContainer.addEventListener(
    "scroll",
    throttle(() => {
      if (messagesContainer.scrollTop < 50 && currentReceiver) {
        loadChatHistory(currentReceiver, true);
      }
    }, 1000)
  );

  // Initialize users list
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

  // Initialize WebSocket
  connectWebSocket();
});
