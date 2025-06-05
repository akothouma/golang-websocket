import { Card } from '../messageComponents/CardComponent/card.js';
import { initSocket } from './socket.js';
import { MessageCarriers } from '../messageComponents/messagesHistoryComponent/message.js';


document.addEventListener('DOMContentLoaded', () => {
    let activeAddMessageFunction = null;

    // Holds the AddMessage for the current chat
    // Don't get AddMessage from MessageCarriers() here directly for the global var initially,
    // as it won't be the one for the displayed chat window

    // ---- NEW ----
    window.setGlobalAddMessageFunction = (newAddMessageFunc) => {
        activeAddMessageFunction = newAddMessageFunc;
        console.log("Global AddMessage function updated.");
    };

    const chatViews = {};
    let currentChatReceiver = null;//Track current chat receiver

    //Initialize socket ONCE
    const socket = initSocket();

    // Store socket reference globally so other components can use it
    window.globalSocket = socket;

    socket.addEventListener("open", () => {
        const request = {
            event: "open",
            payload: {
                messageType: "get_online_users",
            }
        }
        console.log("server connected succesfully");
        socket.send(JSON.stringify(request));
        console.log("request sent succesfully");

    })
    const { showConnections, renderedView } = Card();
    const { AddMessage } = MessageCarriers();

    //store reference to AddMessage globally
    window.addMessageToChat = AddMessage;

    // Function to set current chat receiver
    window.setCurrentChatReceiver = (receiverId) => {
        currentChatReceiver = receiverId;
        console.log("Current chat receiver set to:", receiverId);
        if (!receiverId) { // When going back, no chat is active
            activeAddMessageFunction = null;
            console.log("Global AddMessage function cleared (no active chat).");
        }
    };
    socket.addEventListener("message", (e) => {
        try {
            const data = JSON.parse(e.data);
            console.log("raw backend data", data);
            // IMPORTANT: 'currentUser' is only relevant for "connected_client_list" from data
            // 'senderID' is relevant for "send_private_message" from data
            // 'receiverID' is relevant for "message_sent_confirmation" from data
            //Destructure only whats common 
            const { message, value } = data;//Extract senderID
            switch (message) {
                case "connected_client_list":
                    // 'currentUser' comes from data in this specific message type
                    showConnections(value, data.currentUser);//Accesing it directly
                    break;
                case "send_private_message":
                    // Access data.senderID directly for this message type
                    console.log("Received private data", data.value, "from sender:", data.senderID)
                    // Only add message if we're currently chatting with this sender
                    if (currentChatReceiver === data.senderID) {
                        if (activeAddMessageFunction) {
                            activeAddMessageFunction(data.value, 'right', data.senderID);
                        } else {
                            console.error("No active AddMessage function to display received message!");
                        }
                    } else {
                        console.log(`Message from ${data.senderID} received, but not the current chat. Chatting with: ${currentChatReceiver}`);
                    }
                    break;
                case "message_sent_confirmation":
                    // Access data.receiverID directly for this message type
                    console.log("Message sent confirmation received for message to:", data.receiverID, "value:", data.value);
                    break;

                case "message_history":
                    if (activeAddMessageFunction) {
                        value.forEach((onemessage) => {
                                if (onemessage.sender== currentChatReceiver) {
                                    activeAddMessageFunction(onemessage.message,'right')
                                }
                                else {
                                    activeAddMessageFunction(onemessage.message)
                                }
                            })
                    } else {
                        console.log("couldn't add message history")
                    }
            break;
                default:
            console.log("Unknown message type:", message, "Full data:", data);
        }

        } catch (error) {
        console.error("WebSocket message handling error:", error, "Raw data was:", e.data);
        // return;
    }
})


const root = document.getElementById("message_layout");
root.style.height = '100px';
root.style.display = 'flex';
root.style.flexDirection = 'column';
root.style.background = 'var(--color-white)';
root.style.padding = 'var(--card-padding)';
root.style.borderRadius = 'var(--card-border-radius)';
root.style.fontSize = '1.4rem';
root.style.height = 'max-content';

root.appendChild(renderedView);

    // const { AddMessage } = MessageCarriers();

})