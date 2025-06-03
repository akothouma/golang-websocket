import Image from '../imageComponent/image.js';
import { MessageCarriers } from '../messagesHistoryComponent/message.js';

export const Card = () => {
    const renderedView = document.createElement("div");
    renderedView.style.height = 'fit-content';
    renderedView.style.overflow = 'hidden';
    
    const cardContainer = document.createElement("div");
    cardContainer.style.display = 'flex';
    cardContainer.style.flexDirection = 'column';
    cardContainer.style.alignContent = 'space-between';
    cardContainer.style.gap = '15px';
    renderedView.appendChild(cardContainer);

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
            imageside.appendChild(Image());
            
            const userInfo = document.createElement('div');
            const username = document.createElement('h4');
            username.textContent = oneConnection.Username;
            username.style.margin = '0';
            
            const lastMessageSide = document.createElement('div');
            const lastMessage = document.createElement('p');
            lastMessage.textContent = oneConnection.messageContent || 'No messages yet';
            lastMessage.style.fontSize = '12px';
            lastMessage.style.color = '#666';
            lastMessage.style.margin = '5px 0 0 0';
            
            userInfo.appendChild(username);
            lastMessageSide.appendChild(lastMessage);
            userInfo.appendChild(lastMessageSide);
            
            displayContainer.append(imageside, userInfo);
            onecard.appendChild(displayContainer);
            
            onecard.addEventListener("click", () => {
                console.log("Opening chat with:", oneConnection.Username, "ID:", oneConnection.UserID);
                showPrivateMessages(oneConnection.UserID, cardContainer);
            });
            
            cardContainer.appendChild(onecard);
        });
    }

    return {
        renderedView,
        showConnections
    };
};

function showPrivateMessages(receiverId, cardsView) {
    console.log("Opening private messages with receiver ID:", receiverId);
    
    const renderedView = cardsView.parentNode;
    const messageContainer = document.createElement('div');
    messageContainer.classList.add('message_container');
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
        // Clear current chat receiver when going back
        if (window.setCurrentChatReceiver) {
            window.setCurrentChatReceiver(null);
        }
        renderedView.innerHTML = '';
        renderedView.appendChild(cardsView);
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
    renderedView.innerHTML = '';
    renderedView.appendChild(messageContainer);
}