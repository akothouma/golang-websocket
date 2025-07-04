
/**
 * @file This script defines the `MessageCarriers` component, which is responsible for rendering
 * a complete, interactive chat window for a single conversation. This includes the header,
 * the message history area with infinite scroll, and the message input form.
 */

/**
 * A utility function to limit how often a function can be called. This prevents spamming
 * the server with requests, especially for events like scrolling.
 * @param {Function} func The function to throttle.
 * @param {number} limit The cooldown period in milliseconds.
 * @returns {Function} A new, throttled version of the function.
 */
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

/**
 * The MessageCarriers component factory. It's instantiated each time a user is clicked.
 * It encapsulates all the state and logic for one specific conversation.
 * @param {string} receiverId The unique ID of the user we are chatting with.
 * @param {string} receiverUsername The display name of the user we are chatting with.
 * @returns {HTMLElement} The root DOM element for the chat window.
*/
export const MessageCarriers = (receiverId, receiverUsername) => {
    
    // ---- Component-Scoped State ----
    const socket = window.globalSocket;

    let typingTimer;            // This will hold our timeout
    let isTypingSent = false;   // A flag to prevent sending multiple "start_typing" events
    const TYPING_TIMEOUT_MS = 500; // 2 seconds
    
    /** A flag to prevent multiple history fetch requests from being sent simultaneously. */
    let isFetchingHistory = false;
    
    /** 
     * A helper function to get the timestamp of the oldest message currently displayed.
     * This timestamp acts as a cursor for fetching the next page of history.
     * @returns {string|null} The ISO string timestamp or null if no messages are present.
     */
    const oldestMessageTimestamp = () => {
        const firstMessage = chatHistory.firstChild;
        return firstMessage ? firstMessage.dataset.timestamp : null;
    };
    
    // ---- Global State Registration ---
    
    // Set global variables to let other parts of the app (like app.js) know
    // which chat is currently active.
    window.activeChatUserID = receiverId;

    
    const chatContainer = document.createElement('div');
    chatContainer.className = 'chat-container';
    chatContainer.id = `chat-with-${receiverId}`; // Unique ID for this chat window
    
    const header = document.createElement('header');
    header.className = 'chat-header';

    const backButton = document.createElement('button');
    backButton.className = 'back-button';
    backButton.innerHTML = `<i class="uil uil-arrow-left"></i>`;

    const headerTitle = document.createElement('h3');
    headerTitle.textContent = receiverUsername;

    header.append(backButton, headerTitle);
    
    const chatHistory = document.createElement('div');
    chatHistory.className = 'chat-history';

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
    
    /** Sends a request to the server to fetch a page of message history. */
    function requestHistory(timestamp = null) {
        if (isFetchingHistory) return;
        isFetchingHistory = true;
        
        console.log(`Requesting history for ${receiverId} before ${timestamp || 'the beginning'}`);

        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({
                type: 'get_message_history',
                target: receiverId,
                lastMessageTime: timestamp // The backend handles null time correctly.
            }));
        }
    }

    /** The back button click handler. It resets the UI to show the user list again. */
    backButton.addEventListener('click', () => {
        socket.send(JSON.stringify({ type: 'get_user_list' }))
        const userListContainer = document.querySelector('.user-list-container');
        const messageAreaContainer = document.querySelector('.message-area-container');
         const lastMessage=document.querySelector(".last-message")
        lastMessage.style.display="block"
        
        // Toggle CSS classes to reverse the layout transition.
        if (userListContainer) userListContainer.classList.remove('chat-active');
        if (messageAreaContainer) messageAreaContainer.classList.remove('chat-active');
        
        // Un-highlight any active user card.
        document.querySelectorAll('.user-card.active').forEach(c => c.classList.remove('active'));
        
        // Clean up global state so incoming messages are no longer directed here.
        window.activeChatUserID = null;

    });
    
    /** Renders a single message object into the DOM as a message bubble. */
    function renderMessage(msg, prepend = false) {
        const msgWrapper = document.createElement('div');
        const isSentByMe = msg.sender === window.myUserID;
        msgWrapper.className = `message-wrapper ${isSentByMe ? 'sent' : 'received'}`;
        // Store the message timestamp directly on the element for easy access by our pagination logic.
        msgWrapper.dataset.timestamp = msg.timestamp;

        const date = new Date(msg.timestamp);
        const formattedDate = `${date.toLocaleDateString()} ${date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'})}`;

        msgWrapper.innerHTML = `
            <div class="message-bubble">
                <div class="message-info">
                    <strong>${isSentByMe ? 'You' : receiverUsername}</strong>
                    <span class="timestamp">${formattedDate}</span>
                </div>
                <p class="message-content">${msg.content || msg.message}</p> 
            </div>
        `;
        
        // If `prepend` is true, add the message to the top (for history). Otherwise, append to the bottom (for live messages).
        if (prepend) {
            chatHistory.insertBefore(msgWrapper, chatHistory.firstChild);
        } else {
            chatHistory.appendChild(msgWrapper);
        }
    }
    
    /**
     * This function is exposed globally to be called by `app.js` when history data arrives from the server.
     * It handles rendering the chunk of historical messages.
     * @param {string} target The user ID this history is for, to ensure we're updating the correct chat.
     * @param {Array<Object>} messages An array of message objects.
     */
    window.handleHistoryResponse = (target, messages) => {
        if (target !== receiverId) return; // Ignore if it's not for this chat.
         socket.send(JSON.stringify({
            type: 'mark_messages_as_read',
            target: receiverId // The 'target' is the user whose messages we are now reading.
        }));
        // Keep the user's scroll position stable while adding new content at the top.
        const oldScrollHeight = chatHistory.scrollHeight;
        messages.forEach(msg => renderMessage(msg, true)); // Prepend each historical message.
        const newScrollHeight = chatHistory.scrollHeight;
        chatHistory.scrollTop = newScrollHeight - oldScrollHeight;

        isFetchingHistory = false;
        
        // On the very first load, scroll all the way to the bottom to see the newest messages.
        if (messages.length > 0 && chatHistory.scrollTop === 0) {
            chatHistory.scrollTop = chatHistory.scrollHeight;
        }
    };
    
    /** Exposed globally so `app.js` can push live messages into this active chat window. */
    window.appendMessageToActiveChat = (msg) => {
        renderMessage(msg, false); // Append live messages to the end.
        chatHistory.scrollTop = chatHistory.scrollHeight; // Auto-scroll to the bottom.
        
       socket.send(JSON.stringify({
          type: 'mark_messages_as_read',
          target: receiverId // The 'target' is the user whose messages we are now reading.
      }));
    };

        /**
     * This function is exposed globally to be called by app.js. It creates and
     * removes the typing indicator bubble within this specific chat window.
     * @param {string} type - 'typing_started' or 'typing_stopped'.
     * @param {string} senderID - The ID of the user who is typing.
     */

        window.handleTypingIndicator = (type, senderID) => {
        // Only handle events for the person we are currently chatting with.
        if (senderID !== receiverId) {
            return;
        }

        const indicatorId = `typing-indicator-bubble-${senderID}`;
        const existingIndicator = document.getElementById(indicatorId);

        if (type === 'typing_started') {
            // If the indicator is already showing, do nothing.
            if (existingIndicator) return;

            // Create the bubble wrapper
            const indicatorBubble = document.createElement('div');
            indicatorBubble.id = indicatorId;
            indicatorBubble.className = 'message-wrapper received'; // Style it like a received message

            // Create the content with animated dots
            indicatorBubble.innerHTML = `
                <div class="message-bubble">
                    <div class="typing-dots">
                        <span></span>
                        <span></span>
                        <span></span>
                    </div>
                </div>
            `;
            
            // Add it to the chat and scroll down
            chatHistory.appendChild(indicatorBubble);
            chatHistory.scrollTop = chatHistory.scrollHeight;

        } else if (type === 'typing_stopped') {
            // If the indicator exists, remove it.
            if (existingIndicator) {
                existingIndicator.remove();
            }
        }
    };
    
    /** The throttled scroll event listener for implementing infinite scroll. */
    chatHistory.addEventListener('scroll', throttle(() => {
        // If the user has scrolled to the very top and we are not already fetching...
        if (chatHistory.scrollTop === 0 && !isFetchingHistory) {
            const lastTs = oldestMessageTimestamp();
            if (lastTs) { // ...and if there are messages to paginate from...
                requestHistory(lastTs); //...fetch the next page.
            }
        }
    }, 10)); // Limit to one request per second.

      /** Function to send the stop_typing event and reset the state */
    function sendStopTyping() {
        clearTimeout(typingTimer);
        if (isTypingSent) {
            socket.send(JSON.stringify({ type: 'stop_typing', target: receiverId }));
            isTypingSent = false;
        }
    }

     // Listen for input in the message field
    messageInput.addEventListener('input', () => {
        // If a "start_typing" event hasn't been sent yet, send it.
        if (!isTypingSent) {
            socket.send(JSON.stringify({ type: 'start_typing', target: receiverId }));
            isTypingSent = true;
        }
        
        // Clear any existing timer to reset the timeout
        clearTimeout(typingTimer);
        
        // Set a new timer. If the user doesn't type for TYPING_TIMEOUT_MS,
        // we'll assume they stopped typing.
        typingTimer = setTimeout(sendStopTyping, TYPING_TIMEOUT_MS);
    });
    
    /** The submit event handler for the message input form. */
    messageForm.addEventListener('submit', (e) => {
        e.preventDefault();
         sendStopTyping(); // Immediately signal that typing has stopped.

          // Also immediately remove the indicator from our own view for instant feedback
        const indicator = document.getElementById(`typing-indicator-bubble-${receiverId}`);
        if (indicator) {
            indicator.remove();
        }
        const content = messageInput.value.trim();
        if (content) { // Don't send empty messages.
            const socket = window.globalSocket;
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({
                    type: 'private_message',
                    target: receiverId,
                    content: content,
                }));
                messageInput.value = ''; // Clear the input field.
                // socket.send(JSON.stringify({ type: 'get_user_list' }));
            }
        }
    });

       // Also modify the back button to stop typing when leaving the chat
    backButton.addEventListener('click', () => {
        sendStopTyping(); // Ensure indicator is cleared when navigating away.
        
        // ... (rest of the existing backButton logic) ...
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify({ type: 'get_user_list' }));
        }
        // ...
    });


    // --- Initial Data Fetch ---
    // The component is now fully set up, so we make the initial request for chat history.
    requestHistory();

    // Finally, return the main container element for this component.
    return chatContainer;
};