import Image from '../imageComponent/image.js';
import { MessageCarriers } from '../messagesHistoryComponent/message.js';

export const Card = () => {
    const renderedView = document.createElement("div");
    renderedView.style.height = 'fit-content';
    renderedView.style.overflow = 'hidden';
    renderedView.style.display='flex';
    renderedView.style.flexDirection='row';
    renderedView.style.gap='5px';
    
    const cardContainer = document.createElement("div");
    const messageContainer = document.createElement('div');
    // messageContainer.width='fit-content'
    messageContainer.classList.add('message_container');
    messageContainer.style.display="none";
    cardContainer.style.display = 'flex';
    cardContainer.style.flexDirection = 'column';
    cardContainer.style.alignContent = 'space-between';
    cardContainer.style.gap = '15px';
    renderedView.append(cardContainer,messageContainer);

    function showConnections(data, myID) {
        console.log("From card:", data, myID);
        cardContainer.innerHTML = '';
        
        // Filter out current user from the list
        const filteredData = data.filter(connection => connection.UserID !== myID);
        
        filteredData.forEach((oneConnection) => {
            const onecard = document.createElement("div");
            const displayContainer = document.createElement('div');
            
            onecard.style.width = '100%';
            onecard.style.height = 'fit-content';
            onecard.style.display = 'flex';
            onecard.id = `${oneConnection.Username}`;
            onecard.classList.add('card');
            onecard.style.flexDirection = 'column';
            onecard.style.alignContent = 'space-between';
            onecard.style.gap = '15px';
            onecard.style.cursor = 'pointer';
            onecard.style.padding = '10px';
            onecard.style.border = '1px solid #ccc';
            onecard.style.borderRadius = '5px';
            onecard.style.marginBottom = '5px';
            
            displayContainer.style.display = 'flex';
            displayContainer.style.flexDirection = 'row';
            displayContainer.style.alignItems = 'center';
            displayContainer.style.gap = '10px';
            
            const imageside = document.createElement('div');
            imageside.style.display="flex"
            imageside.style.flexDirection='column'
            const username = document.createElement('p');
            username.style.fontSize = '8px';
            username.textContent = oneConnection.Username;
            username.style.margin = '0';
            imageside.appendChild(Image(),username);
            
            const userInfo = document.createElement('div');
            userInfo.style.width='fit-content'
            
            const lastMessageSide = document.createElement('div');
            const lastMessage = document.createElement('p');
            lastMessage.style.width='inherit'
            lastMessage.textContent = oneConnection.messageContent || 'No messages yet';
            lastMessage.style.fontSize = '8px';
            lastMessage.style.color = '#666';
            lastMessage.style.margin = '5px 0 0 0';
            
            userInfo.appendChild(username);
            lastMessageSide.appendChild(lastMessage);
            userInfo.appendChild(lastMessageSide);
            
            displayContainer.append(imageside, userInfo);
            onecard.appendChild(displayContainer);
            
            onecard.addEventListener("click", () => {
                console.log("Opening chat with:", oneConnection.Username, "ID:", oneConnection.UserID);
                cardContainer.style.width='20%'
                messageContainer.style.display="block"
                messageContainer.style.width='80%'
                lastMessage.style.display="none"
                userInfo.style.display="none"
                
                showPrivateMessages(oneConnection.UserID, cardContainer,lastMessage,userInfo);
            });
            
            cardContainer.appendChild(onecard);
        });
    }

    return {
        renderedView,
        showConnections
    };
};
function showPrivateMessages(receiverId, cardsView,lastMessageElement,metadata) {

    
    const request={
        event:"frontend request",
        payload:{
            "messageType":"get_message_history",
            "receiverID":receiverId,
        }
    }
    const socket=window.globalSocket;
    if (socket && socket.readyState == WebSocket.OPEN){
        socket.send(JSON.stringify(request))
    }else{
        console.log("failed to send request for message history")
    }
    console.log("Opening private messages with receiver ID:", receiverId);
    
    const messageContainer = document.querySelector('.message_container');
    messageContainer.innerHTML=''
    messageContainer.id = `${receiverId}`;
    
    const backButton = document.createElement('button');
    backButton.textContent = "← Back";
    backButton.style.padding = '10px 15px';
    backButton.style.marginBottom = '10px';
    backButton.style.backgroundColor = '#007bff';
    backButton.style.color = 'white';
    backButton.style.border = 'none';
    backButton.style.borderRadius = '5px';
    backButton.style.cursor = 'pointer';
    
    backButton.addEventListener("click", () => {
        messageContainer.innerHTML=''
        messageContainer.style.display="none"
        cardsView.style.width='100%'
        lastMessageElement.style.display="block"
        metadata.style.display="block"
        // Clear current chat receiver when going back
        if (window.setCurrentChatReceiver) {
            window.setCurrentChatReceiver(null);
        }
    });

    const { chatContainer,AddMessage:AddMessageForThisChat } = MessageCarriers();
    chatContainer.id = receiverId;

      // ---- NEW ----
    // Update the global reference to point to the AddMessage of THIS new chat window
    if (window.setGlobalAddMessageFunction) {
        window.setGlobalAddMessageFunction(AddMessageForThisChat);
    }
    
    // Set current chat receiver
    if (window.setCurrentChatReceiver) {
        window.setCurrentChatReceiver(receiverId);
    }
    
    messageContainer.append(backButton, chatContainer);
 
}