// ui/messageComponents/messagesHistoryComponent/message.js

// A utility function to prevent an event from firing too often.
function throttle(func, limit) {
    let inThrottle;
    return function() {
        const args = arguments;
        const context = this;
        if (!inThrottle) {
            func.apply(context, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

export const MessageCarriers = (receiverId, receiverUsername) => {
    // ---- Component State ----
    let isFetchingHistory = false;
    const oldestMessageTimestamp = () => {
        // Get the timestamp from the first message element's dataset
        const firstMessage = chatHistory.firstChild;
        return firstMessage ? firstMessage.dataset.timestamp : null;
    };
    
    // Set the global active chat ID so app.js knows where to direct incoming messages
    window.activeChatUserID = receiverId;

    // ---- DOM Elements ----
    const chatContainer = document.createElement('div');
    chatContainer.className = 'chat-container';
    chatContainer.id = `chat-with-${receiverId}`;
    
    // Header
    const header = document.createElement('header');
    header.className = 'chat-header';

    
    const backButton = document.createElement('button');
    backButton.className = 'back-button';
    // Using an icon from the library you already have
  backButton.innerHTML = `<i class="uil uil-arrow-left"></i>`;

const headerTitle = document.createElement('h3');
    headerTitle.textContent = receiverUsername;

    header.append(backButton, headerTitle);

    // Message History Area
    const chatHistory = document.createElement('div');
    chatHistory.className = 'chat-history';

    // Message Input Form
    const messageForm = document.createElement('form');
    messageForm.className = 'message-form';
    const messageInput = document.createElement('input');
    messageInput.type = 'text';
    messageInput.placeholder = 'Type a message...';
    messageInput.autocomplete = 'off';
    const sendButton = document.createElement('button');
    sendButton.type = 'submit';
    sendButton.textContent = 'Send';
    messageForm.append(messageInput, sendButton);

    chatContainer.append(header, chatHistory, messageForm);
    
    // ---- Logic and Event Handlers ----
    function requestHistory(timestamp = null) {
        if (isFetchingHistory) return;
        isFetchingHistory = true;
        
        console.log(`Requesting history for ${receiverId} before ${timestamp || 'the beginning'}`);

        const socket = window.globalSocket;
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({
                type: 'get_message_history',
                target: receiverId,
                lastMessageTime: timestamp // Backend will handle null time
            }));
        }
    }

     // ADD THIS NEW EVENT LISTENER
    backButton.addEventListener('click', () => {
        const userListContainer = document.querySelector('.user-list-container');
        const messageAreaContainer = document.querySelector('.message-area-container');

        if (userListContainer) userListContainer.classList.remove('chat-active');
        if (messageAreaContainer) messageAreaContainer.classList.remove('chat-active');
        
        document.querySelectorAll('.user-card.active').forEach(c => c.classList.remove('active'));
        
        // Cleanup global state
        window.activeChatUserID = null;
    });
    
    // Function to add a single message to the DOM.
    function renderMessage(msg, prepend = false) {
        const msgWrapper = document.createElement('div');
        msgWrapper.className = `message-wrapper ${msg.sender === window.myUserID ? 'sent' : 'received'}`;
        // Store timestamp for pagination
        msgWrapper.dataset.timestamp = msg.timestamp;

        const date = new Date(msg.timestamp);
        const formattedDate = `${date.toLocaleDateString()} ${date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}`;

        msgWrapper.innerHTML = `
            <div class="message-bubble">
                <div class="message-info">
                    <strong>${msg.sender === window.myUserID ? 'You' : receiverUsername}</strong>
                    <span class="timestamp">${formattedDate}</span>
                </div>
                <p class="message-content">${msg.content || msg.message}</p> 
            </div>
        `;
        
        if (prepend) {
            chatHistory.insertBefore(msgWrapper, chatHistory.firstChild);
        } else {
            chatHistory.appendChild(msgWrapper);
        }
    }
    
    // Expose a function to app.js to handle incoming history chunks
    window.handleHistoryResponse = (target, messages) => {
        if (target !== receiverId) return; // History for a different chat
        
        const oldScrollHeight = chatHistory.scrollHeight;

        messages.forEach(msg => renderMessage(msg, true)); // Prepend older messages

        // Maintain scroll position after prepending
        const newScrollHeight = chatHistory.scrollHeight;
        chatHistory.scrollTop = newScrollHeight - oldScrollHeight;

        isFetchingHistory = false;
        
        // If it's the initial load, scroll to bottom
        if (!oldestMessageTimestamp()) {
             chatHistory.scrollTop = chatHistory.scrollHeight;
        }
    };
    
    // Expose a function for app.js to append live messages
    window.appendMessageToActiveChat = (msg) => {
        renderMessage(msg, false);
        chatHistory.scrollTop = chatHistory.scrollHeight; // Scroll to bottom on new message
    };
    
    // Infinite scroll listener
    chatHistory.addEventListener('scroll', throttle(() => {
        if (chatHistory.scrollTop === 0 && !isFetchingHistory) {
            const lastTs = oldestMessageTimestamp();
            if (lastTs) {
                requestHistory(lastTs);
            }
        }
    }, 1000)); // Throttle to run at most once per second
    
    // Message sending listener
    messageForm.addEventListener('submit', (e) => {
        e.preventDefault();
        const content = messageInput.value.trim();
        if (content) {
            const socket = window.globalSocket;
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({
                    type: 'private_message',
                    target: receiverId,
                    content: content,
                }));
                messageInput.value = '';
            }
        }
    });

    // Initial fetch
    if(window.chatCache[receiverId]) {
         // If we have a cache, render it first for a snappy UI
        window.chatCache[receiverId].forEach(msg => renderMessage(msg));
        chatHistory.scrollTop = chatHistory.scrollHeight;
        // Then maybe fetch newer messages? For now, we'll just fetch all on open.
    }
    
    requestHistory();

    return chatContainer;
};

// Remove the window functions when the component is "unmounted" (the user clicks back or closes the chat)
// You would need to add this cleanup logic in card.js when the chat window is replaced.
// For example:
// if (window.handleHistoryResponse) window.handleHistoryResponse = null;
// if (window.appendMessageToActiveChat) window.appendMessageToActiveChat = null;