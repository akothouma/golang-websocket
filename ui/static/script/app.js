
/**
 * @file This script is the main entry point and controller for the chat application.
 * It initializes the WebSocket connection, manages the overall state, and dispatches
 * incoming messages from the server to the appropriate UI components.
 */

// Import UI components and the WebSocket initializer.
import { Card } from '../messageComponents/CardComponent/card.js';
import { initSocket } from './socket.js';

// The main logic runs after the entire HTML document has been loaded.
document.addEventListener('DOMContentLoaded', () => {
    
    // --- Global State Management ---
    // These variables are attached to the `window` object to be accessible
    // across different component files, acting as a simple shared state store.
    const root = document.getElementById("message_layout");

    if (root){

    /**
     * @property {Object} window.chatCache - A simple cache to store message history for each chat.
     * The key is the other user's ID, and the value is an array of message objects.
     * This can be used for performance optimizations later.
     */
    window.chatCache = {}; 
    
    /**
     * @property {string} window.myUserID - The unique ID of the currently logged-in user.
     * This is crucial for determining if a message is 'sent' or 'received'. It is set
     * once upon receiving the first `user_list_update` from the server.
     */
    window.myUserID = '';

    // --- Component Initialization ---

    // `Card` is our main UI component factory. It returns the rendered HTML element
    // and a function to update its content (the user list).
    const { renderedView, updateUserList } = Card();
    
    // Find the designated layout container in `index.html` and append our chat system's root element.
    root.appendChild(renderedView);

    // --- WebSocket Initialization & Message Handling ---
    
    // Initialize the WebSocket connection. `initSocket` ensures only one connection is made.
    const socket = initSocket();
    // Store the socket instance globally so other components can use it to send messages.
    window.globalSocket = socket;

    /**
     * The central message listener. It acts as a router, parsing incoming JSON data from the server
     * and deciding which function or component should handle it based on the `data.type` field.
     */
    socket.addEventListener("message", (e) => {
        try {
            const data = JSON.parse(e.data);
            console.log("Received data from backend:", data);

            // Use a switch statement to route incoming messages.
            switch (data.type) {
                // The server sends this when any user connects, disconnects, or sends a message.
                // It contains the complete, up-to-date state of all users.
                case 'user_list_update':
                    if (data.userList) {
                        // On the first `user_list_update`, we identify our own user ID
                        // from the user object that has `isMe: true`.
                        const me = data.userList.find(u => u.isMe);
                        if (me && !window.myUserID) {
                            window.myUserID = me.userID;
                            console.log("My user ID has been set to:", window.myUserID);
                        }
                        // Pass the fresh user list to the Card component to re-render the UI.
                        updateUserList(data.userList);
                    }
                    break;

                // A real-time private message has arrived from the server.
                case 'private_message':
                    handleIncomingPrivateMessage(data);
                    break;
                
                // The server is responding to our `get_message_history` request.
                case 'history_response':
                    // We delegate the handling of this response to the currently active chat window.
                    // The `message.js` component creates `window.handleHistoryResponse` when it's active.
                    if (window.handleHistoryResponse) {
                       window.handleHistoryResponse(data.target, data.messages);
                    }
                    break;

                             // ---- START: NEW TYPING HANDLERS ----
                case 'typing_started':
                    updateTypingIndicator(data.sender, true);
                    break;

                case 'typing_stopped':
                    updateTypingIndicator(data.sender, false);
                    break;
                // ---- END: NEW TYPING HANDLERS ----


                // Catch-all for any message types we don't recognize.
                default:
                    console.warn("Unknown message type:", data.type, "Full data:", data);
            }
        } catch (error) {
            console.error("WebSocket message handling error:", error, "Raw data was:", e.data);
        }
    });

     /**
     * Updates the UI to show or hide the "typing..." indicator for a specific user.
     * @param {string} senderID The ID of the user who is typing.
     * @param {boolean} isTyping True to show the indicator, false to hide it.
     */
    function updateTypingIndicator(senderID, isTyping) {
        const userCard = document.querySelector(`.user-card[data-user-id='${senderID}']`);
        if (!userCard) return;

        const lastMessageElement = userCard.querySelector('.last-message');
        if (!lastMessageElement) return;

        if (isTyping) {
            // Store the original message content before replacing it
            lastMessageElement.dataset.originalText = lastMessageElement.textContent;
            // Show the typing indicator
            lastMessageElement.innerHTML = `<span class="typing-indicator">typing...</span>`;
        } else {
            // Restore the original message content if it was saved
            if (lastMessageElement.dataset.originalText) {
                lastMessageElement.textContent = lastMessageElement.dataset.originalText;
            }
        }
    }

    /**
     * Handles an incoming `private_message`. It determines who the message is from/to,
     * adds it to the cache, and appends it to the active chat window if it's open.
     * @param {Object} data - The WebSocket message payload from the server.
     */
    function handleIncomingPrivateMessage(data) {
        const message = data.messages[0]; // The full message object is nested.
        if (!message) return;

        // Determine the ID of the other person in the conversation, regardless of
        // whether we were the sender or the receiver.
        const otherUserID = message.sender === window.myUserID ? message.receiver : message.sender;

        // Cache the message. (Currently simple, can be expanded for more features).
        if (!window.chatCache[otherUserID]) {
            window.chatCache[otherUserID] = [];
        }
        window.chatCache[otherUserID].push(message);

        // If the chat with this user is currently active, call the function exposed by
        // `message.js` to append the message to the DOM in real-time.
        if (window.appendMessageToActiveChat && window.activeChatUserID === otherUserID) {
            window.appendMessageToActiveChat(message);
        } else {
            // If the chat is not active, we just log it. This is where you would implement
            // a notification badge on the user's card in the list.
            console.log(`New message from ${otherUserID}, but chat is not active.`);
        }
    }
}
});