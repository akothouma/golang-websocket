// ui/messageComponents/CardComponent/card.js
import { MessageCarriers } from '../messagesHistoryComponent/message.js';

export const Card = () => {
    const renderedView = document.createElement("div");
    renderedView.className = 'chat-wrapper';

    const cardContainer = document.createElement("div");
    cardContainer.className = 'user-list-container';
    
    const messageArea = document.createElement('div');
    messageArea.className = 'message-area-container';
    // No more style.display here, we use CSS classes

    renderedView.append(cardContainer, messageArea);

    function updateUserList(users) {
        users.sort((a, b) => {
            const timeA = a.lastMessageTime ? new Date(a.lastMessageTime).getTime() : 0;
            const timeB = b.lastMessageTime ? new Date(b.lastMessageTime).getTime() : 0;
            
            if (timeA !== timeB) {
                return timeB - timeA;
            }
            return a.username.localeCompare(b.username); 
        });

        cardContainer.innerHTML = ''; 

        users.forEach((user) => {
            if (user.userID === window.myUserID) return;

            const userCard = document.createElement("div");
            userCard.className = 'user-card';
            userCard.dataset.userId = user.userID; 

              // The user's status, including a class for styling and the text itself
            const statusHTML = user.isOnline 
                ? `<span class="user-status online">Online</span>`
                : `<span class="user-status offline">Offline</span>`;

            userCard.innerHTML = `
                <div class="user-card-avatar ${user.isOnline ? 'online' : ''}">
                      <div class="initials">${user.username.charAt(0).toUpperCase()}</div>
                    <div class="status-dot"></div>
                </div>
                <div class="user-card-info">
                    // FIX: Use lowercase 'username'
                    <strong class="username">${user.username}</strong>
                    <p class="last-message">${user.lastMessageContent || 'Click to chat'}</p>
                       ${statusHTML}
                </div>
            `;
            
            userCard.addEventListener("click", () => openChatWith(user));
            
            cardContainer.appendChild(userCard);
        });
    }

    // This function was provided in a previous fix and might have been updated.
    // Ensure it uses the classList toggle approach.
    function openChatWith(user) {
        cardContainer.classList.add('chat-active');
        messageArea.classList.add('chat-active');
        
        document.querySelectorAll('.user-card.active').forEach(c => c.classList.remove('active'));
       
        const cardElement = document.querySelector(`.user-card[data-user-id='${user.userID}']`);
        if (cardElement) {
            cardElement.classList.add('active');
        }

        if (window.handleHistoryResponse) window.handleHistoryResponse = null;
        if (window.appendMessageToActiveChat) window.appendMessageToActiveChat = null;
        window.activeChatUserID = null;
    
        const chatWindow = MessageCarriers(user.userID, user.username);
        messageArea.innerHTML = '';
        messageArea.appendChild(chatWindow);
    }

    return { renderedView, updateUserList };
};