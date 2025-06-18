document.addEventListener("DOMContentLoaded", function () {
  // Add notification styles
  const notificationStyles = document.createElement('style');
  notificationStyles.textContent = `
  .notification-badge {
    position: absolute;
    top: -2px;
    right: -2px;
    background: #ff4444;
    color: white;
    border-radius: 50%;
    width: 18px;
    height: 18px;
    font-size: 11px;
    font-weight: bold;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 2px solid white;
    box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    animation: pulse 2s infinite;
  }
  
  .notification-dot {
    position: absolute;
    top: 2px;
    right: 2px;
    background: #ff4444;
    border-radius: 50%;
    width: 8px;
    height: 8px;
    border: 1px solid white;
    animation: pulse 2s infinite;
  }
  
  @keyframes pulse {
    0% { transform: scale(1); }
    50% { transform: scale(1.1); }
    100% { transform: scale(1); }
  }
  
  .user-avatar {
    position: relative;
  }
`;
  document.head.appendChild(notificationStyles);

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
  const messageInputContainer = document.getElementById("messageInputContainer");
  const searchInput = document.getElementById("searchInput");
  const backButton = document.getElementById("backButton");
  const usersSidebar = document.getElementById("usersSidebar");

  let socket;
  let currentReceiver = null;
  let currentReceiverName = null;
  let typingTimeout;
  let messageOffset = 0;
  let isLoadingMessages = false;
  let hasMoreMessages = true;
  let allUsers = [];
  let loadingIndicator = null;
  let unreadMessages = new Map(); // Track unread messages per user

  // Improved throttle function with leading and trailing execution
  function throttle(func, limit) {
    let inThrottle;
    let lastFunc;
    let lastRan;
    return function() {
      const context = this;
      const args = arguments;
      if (!inThrottle) {
        func.apply(context, args);
        lastRan = Date.now();
        inThrottle = true;
      } else {
        clearTimeout(lastFunc);
        lastFunc = setTimeout(function() {
          if ((Date.now() - lastRan) >= limit) {
            func.apply(context, args);
            lastRan = Date.now();
          }
        }, limit - (Date.now() - lastRan));
      }
    }
  }

  // Debounce function for more sensitive operations
  function debounce(func, wait, immediate) {
    let timeout;
    return function() {
      const context = this;
      const args = arguments;
      const later = function() {
        timeout = null;
        if (!immediate) func.apply(context, args);
      };
      const callNow = immediate && !timeout;
      clearTimeout(timeout);
      timeout = setTimeout(later, wait);
      if (callNow) func.apply(context, args);
    };
  }

  // Create and manage loading indicator
  function createLoadingIndicator() {
    if (loadingIndicator) return loadingIndicator;
    
    loadingIndicator = document.createElement("div");
    loadingIndicator.className = "loading-indicator";
    loadingIndicator.innerHTML = `
      <div style="display: flex; justify-content: center; padding: 10px; color: #666;">
        <div style="display: flex; align-items: center; gap: 8px;">
          <div style="width: 16px; height: 16px; border: 2px solid #f3f3f3; border-top: 2px solid #3498db; border-radius: 50%; animation: spin 1s linear infinite;"></div>
          <span style="font-size: 12px;">Loading messages...</span>
        </div>
      </div>
      <style>
        @keyframes spin {
          0% { transform: rotate(0deg); }
          100% { transform: rotate(360deg); }
        }
      </style>
    `;
    return loadingIndicator;
  }

  function showLoadingIndicator() {
    if (!loadingIndicator && hasMoreMessages) {
      const indicator = createLoadingIndicator();
      messagesContainer.prepend(indicator);
    }
  }

  function hideLoadingIndicator() {
    if (loadingIndicator && loadingIndicator.parentNode) {
      loadingIndicator.parentNode.removeChild(loadingIndicator);
      loadingIndicator = null;
    }
  }

  // Improved scroll handler with better detection
  const handleScroll = throttle(() => {
    if (!currentReceiver || isLoadingMessages || !hasMoreMessages) return;

    const scrollTop = messagesContainer.scrollTop;
    const scrollThreshold = 100; // Load more when within 100px of top

    // Check if user scrolled near the top
    if (scrollTop <= scrollThreshold) {
      loadMoreMessages();
    }
  }, 300); // Throttle to 300ms

  // Enhanced message loading function
  function loadMoreMessages() {
    if (isLoadingMessages || !hasMoreMessages || !currentReceiver) return;

    isLoadingMessages = true;
    showLoadingIndicator();

    // Store current scroll position and height for restoration
    const scrollHeightBefore = messagesContainer.scrollHeight;
    const scrollTopBefore = messagesContainer.scrollTop;

    fetch(`/messages?with=${currentReceiver}&offset=${messageOffset}&limit=10`)
      .then((response) => {
        if (!response.ok) throw new Error("Failed to fetch messages");
        return response.json();
      })
      .then((messages) => {
        hideLoadingIndicator();

        if (!Array.isArray(messages)) {
          console.error("Messages response is not an array:", messages);
          return;
        }

        // If we got fewer than 10 messages, we've reached the end
        if (messages.length < 10) {
          hasMoreMessages = false;
        }

        if (messages.length === 0) {
          return;
        }

        // Add messages to the top of the container
        const fragment = document.createDocumentFragment();
        messages.reverse().forEach((message) => {
          const messageElement = createMessageElement(message);
          fragment.appendChild(messageElement);
        });

        // Insert all messages at once for better performance
        messagesContainer.insertBefore(fragment, messagesContainer.firstChild);

        // Restore scroll position to maintain user's view
        const scrollHeightAfter = messagesContainer.scrollHeight;
        const heightDifference = scrollHeightAfter - scrollHeightBefore;
        messagesContainer.scrollTop = scrollTopBefore + heightDifference;

        messageOffset += messages.length;
      })
      .catch((error) => {
        console.error("Error loading more messages:", error);
        hideLoadingIndicator();
        // Show error message to user
        showErrorMessage("Failed to load messages. Please try again.");
      })
      .finally(() => {
        isLoadingMessages = false;
      });
  }

  function showErrorMessage(message) {
    const errorDiv = document.createElement("div");
    errorDiv.className = "error-message";
    errorDiv.style.cssText = `
      background: #fee; 
      color: #c33; 
      padding: 8px 12px; 
      margin: 8px; 
      border-radius: 4px; 
      text-align: center; 
      font-size: 12px;
      border: 1px solid #fcc;
    `;
    errorDiv.textContent = message;
    messagesContainer.prepend(errorDiv);
    
    // Remove error message after 3 seconds
    setTimeout(() => {
      if (errorDiv.parentNode) {
        errorDiv.parentNode.removeChild(errorDiv);
      }
    }, 3000);
  }

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
    hasMoreMessages = true; // Reset for new chat

    // Clear unread messages for this user
    clearUnreadMessages(userId);

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

    // Load initial chat history
    loadInitialChatHistory(currentReceiver);

    // Focus on input
    messageInput.focus();
  }

  function loadInitialChatHistory(otherUserID) {
    if (isLoadingMessages) return;

    isLoadingMessages = true;
    messagesContainer.innerHTML = "";

    fetch(`/messages?with=${otherUserID}&offset=0&limit=10`)
      .then((response) => {
        if (!response.ok) throw new Error("Failed to fetch messages");
        return response.json();
      })
      .then((messages) => {
        if (!Array.isArray(messages)) {
          console.error("Chat history is not an array:", messages);
          return;
        }

        if (messages.length < 10) {
          hasMoreMessages = false;
        }

        messages.reverse().forEach((message) => {
          const messageElement = createMessageElement(message);
          messagesContainer.appendChild(messageElement);
        });

        // Scroll to bottom for initial load
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
        messageOffset = messages.length;
      })
      .catch((error) => {
        console.error("Error loading chat history:", error);
        showErrorMessage("Failed to load chat history.");
      })
      .finally(() => {
        isLoadingMessages = false;
      });
  }

  // Modified createMessageElement function to include sender in message content
  function createMessageElement(message) {
    const messageElement = document.createElement("div");
    const isSentByCurrentUser = message.sender_id === currentUserID;
    messageElement.className = `message ${isSentByCurrentUser ? "sent" : "received"}`;

    // Get sender name
    const senderName = isSentByCurrentUser ? currentUserName : currentReceiverName;
    
    // Create the main message content container
    const messageContent = document.createElement("div");
    messageContent.className = "message-content";

    // Create sender name element (part of message content now)
    const senderElement = document.createElement("div");
    senderElement.className = "message-sender";
    senderElement.textContent = senderName;
    senderElement.style.cssText = `
      font-size: 0.75rem;
      font-weight: 600;
      margin-bottom: 4px;
      color: ${isSentByCurrentUser ? 'rgba(255, 255, 255, 0.8)' : 'rgba(0, 0, 0, 0.7)'};
      opacity: 0.9;
    `;

    // Create message text element
    const messageText = document.createElement("div");
    messageText.className = "message-text";
    messageText.textContent = message.content;
    messageText.style.cssText = `
      margin-bottom: 4px;
      line-height: 1.4;
      word-wrap: break-word;
    `;

    // Create timestamp element
    const messageTime = document.createElement("div");
    messageTime.className = "message-time";
    messageTime.textContent = new Date(message.created_at).toLocaleTimeString([], {
      hour: "2-digit",
      minute: "2-digit",
    });
    messageTime.style.cssText = `
      font-size: 0.65rem;
      opacity: 0.7;
      text-align: ${isSentByCurrentUser ? 'right' : 'left'};
      margin-top: 2px;
      color: ${isSentByCurrentUser ? 'rgba(255, 255, 255, 0.6)' : 'rgba(0, 0, 0, 0.5)'};
    `;

    // Assemble the message content
    messageContent.appendChild(senderElement);
    messageContent.appendChild(messageText);
    messageContent.appendChild(messageTime);

    // Add the content to the message element
    messageElement.appendChild(messageContent);

    // Add some additional styling to the message element
    messageElement.style.cssText = `
      margin-bottom: 12px;
      max-width: 70%;
      ${isSentByCurrentUser ? 'margin-left: auto;' : 'margin-right: auto;'}
    `;

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
    
    // Add notification logic for messages not from current chat
    if (message.sender_id !== currentUserID && message.sender_id !== currentReceiver) {
      addUnreadMessage(message.sender_id);
      updateNotificationBadge(message.sender_id);
    }
  }

  function addUnreadMessage(userId) {
    const currentCount = unreadMessages.get(userId) || 0;
    unreadMessages.set(userId, currentCount + 1);
  }

  function clearUnreadMessages(userId) {
    unreadMessages.delete(userId);
    removeNotificationBadge(userId);
  }

  function updateNotificationBadge(userId) {
    const userItem = document.querySelector(`[data-user-id="${userId}"]`);
    if (!userItem) return;
    
    const userAvatar = userItem.querySelector('.user-avatar');
    const existingBadge = userAvatar.querySelector('.notification-badge, .notification-dot');
    
    if (existingBadge) {
      existingBadge.remove();
    }
    
    const unreadCount = unreadMessages.get(userId) || 0;
    if (unreadCount > 0) {
      const badge = document.createElement('div');
      
      if (unreadCount > 9) {
        badge.className = 'notification-badge';
        badge.textContent = '9+';
      } else if (unreadCount > 1) {
        badge.className = 'notification-badge';
        badge.textContent = unreadCount.toString();
      } else {
        badge.className = 'notification-dot';
      }
      
      userAvatar.appendChild(badge);
    }
  }

  function removeNotificationBadge(userId) {
    const userItem = document.querySelector(`[data-user-id="${userId}"]`);
    if (!userItem) return;
    
    const badge = userItem.querySelector('.notification-badge, .notification-dot');
    if (badge) {
      badge.remove();
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

    socket.send(
      JSON.stringify({
        type: "stop_typing",
        receiver_id: currentReceiver,
      })
    );

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
      <div style="display: flex; align-items: center; gap: 8px; padding: 8px 12px; background: #f0f0f0; border-radius: 18px; margin: 8px 0; max-width: 70%;">
        <div style="font-size: 0.75rem; font-weight: 600; color: rgba(0, 0, 0, 0.7);">${currentReceiverName}</div>
        <div class="typing-dots" style="display: flex; gap: 2px;">
          <div class="typing-dot" style="width: 4px; height: 4px; background: #999; border-radius: 50%; animation: typing 1.4s infinite;"></div>
          <div class="typing-dot" style="width: 4px; height: 4px; background: #999; border-radius: 50%; animation: typing 1.4s infinite 0.2s;"></div>
          <div class="typing-dot" style="width: 4px; height: 4px; background: #999; border-radius: 50%; animation: typing 1.4s infinite 0.4s;"></div>
        </div>
        <span style="font-size: 0.65rem; color: rgba(0, 0, 0, 0.5);">typing...</span>
      </div>
      <style>
        @keyframes typing {
          0%, 60%, 100% { transform: translateY(0); opacity: 0.4; }
          30% { transform: translateY(-8px); opacity: 1; }
        }
      </style>
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

  // Event Listeners
  searchInput.addEventListener("input", debounce(renderFilteredUsers, 300));

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
    messageOffset = 0;
    hasMoreMessages = true;
    hideLoadingIndicator();
  });

  // Improved scroll event listener
  messagesContainer.addEventListener("scroll", handleScroll);

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