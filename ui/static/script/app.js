// ui/static/script/app.js

import { Card } from '../messageComponents/CardComponent/card.js';
import { initSocket } from './socket.js';

document.addEventListener('DOMContentLoaded', () => {
    window.chatCache = {}; 
    window.myUserID = ''; // Will be set when we get the first user list.

    const { renderedView, updateUserList } = Card();
    
    const root = document.getElementById("message_layout");
    root.appendChild(renderedView);

    const socket = initSocket();
    window.globalSocket = socket;

    socket.addEventListener("message", (e) => {
        try {
            const data = JSON.parse(e.data);
            console.log("Received data from backend:", data);

            switch (data.type) {
                case 'user_list_update':
                    if (data.userList) {
                        // FIX: Reliably find my user ID and store it globally
                        const me = data.userList.find(u => u.isMe);
                        if (me && !window.myUserID) {
                            window.myUserID = me.userID;
                            console.log("My user ID has been set to:", window.myUserID);
                        }
                        updateUserList(data.userList);
                    }
                    break;
                // ... rest of the switch case is correct ...
                case 'private_message':
                    handleIncomingPrivateMessage(data);
                    break;
                case 'history_response':
                    if (window.handleHistoryResponse) {
                       window.handleHistoryResponse(data.target, data.messages);
                    }
                    break;
                default:
                    console.warn("Unknown message type:", data.type, "Full data:", data);
            }
        } catch (error) {
            console.error("WebSocket message handling error:", error, "Raw data was:", e.data);
        }
    });

    // ... handleIncomingPrivateMessage is correct, no changes needed here ...
    function handleIncomingPrivateMessage(data) {
        // This function will now work correctly because window.myUserID is set.
        const message = data.messages[0];
        if (!message) return;
        const otherUserID = message.sender === window.myUserID ? message.receiver : message.sender;

        if (!window.chatCache[otherUserID]) {
            window.chatCache[otherUserID] = [];
        }
        window.chatCache[otherUserID].push(message);

        if (window.appendMessageToActiveChat && window.activeChatUserID === otherUserID) {
            window.appendMessageToActiveChat(message);
        } else {
            console.log(`New message from ${otherUserID}, but chat is not active.`);
        }
    }
});