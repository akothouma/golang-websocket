import { Card } from '../messageComponents/CardComponent/card.js';
import {initSocket} from './socket.js';
import { MessageCarriers } from '../messageComponents/messagesHistoryComponent/message.js';


document.addEventListener('DOMContentLoaded', () => {
   const socket=initSocket();

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
    const {AddMessage}=MessageCarriers();
    socket.addEventListener("message", (e) => {
        try {
            const data = JSON.parse(e.data);
            console.log("raw backend data",data);
            const { message, value,currentUser} = data;
            switch (message) {
                case "connected_client_list":
                    showConnections(value,currentUser);
                    break;
                case "send_private_message":
                    console.log("received private data",value)
                    AddMessage(value,'right')
                break;
            }

        } catch (error) {
            console.error("WebSocket message handling error:", error);
            return;
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