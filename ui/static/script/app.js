import { Card } from '../messageComponents/CardComponent/card.js';
import {initSocket} from './socket.js';
import { MessageCarriers } from '../messageComponents/messagesHistoryComponent/message.js';


document.addEventListener('DOMContentLoaded', () => {

    const chatViews={};
    let currentChatReceiver=null;//Track current chat receiver

    //Initialize socket ONCE
   const socket=initSocket();

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
    const { showConnections, renderedView} = Card();
    const {AddMessage}=MessageCarriers();

    //store reference to AddMessage globally
     window.addMessageToChat = AddMessage;

        // Function to set current chat receiver
    window.setCurrentChatReceiver = (receiverId) => {
        currentChatReceiver = receiverId;
        console.log("Current chat receiver set to:", receiverId);
    };
    socket.addEventListener("message", (e) => {
        try {
            const data = JSON.parse(e.data);
            console.log("raw backend data",data);
            const { message, value,currentUser} = data;//Extract senderID
            switch (message) {
                case "connected_client_list":
                    showConnections(value,currentUser);
                    break;
                case "send_private_message":
                    console.log("Received private data",value,"from:",senderID)
                     // Only add message if we're currently chatting with this sender
                    if (currentChatReceiver === data.senderID) {
                        AddMessage(data.value, 'left', senderID);
                    }else {
                    console.log(`Message from ${data.senderID} received, but not the current chat. Chatting with: ${currentChatReceiver}`);
                }
                    break;
                case "message_sent_confirmation":
                     console.log("Message sent confirmation received for message to:", data.receiverID, "value:", data.value); // Access data.receiverID if needed
                    break;
                default:
                    console.log("Unknown message type:", message);
            }

        } catch (error) {
            console.error("WebSocket message handling error:", error);
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