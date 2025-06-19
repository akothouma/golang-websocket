
/**
 * @file This script defines the `Card` component, which is responsible for rendering
 * the user list sidebar and managing the transition to the main chat window. It acts
 * as a factory function, creating and returning the necessary DOM elements and update logic.
 */

// Import the component that creates the actual chat window (history, input, etc.).
import { MessageCarriers } from '../messagesHistoryComponent/message.js';

/**
 * The Card component factory. When called, it creates the DOM structure for the entire
 * chat area (user list + message panel) and returns methods to interact with it.
 * @returns {{renderedView: HTMLElement, updateUserList: Function}} An object containing the
 * root DOM element (`renderedView`) to be appended to the page, and the `updateUserList` function.
 */
export const Card = () => {
    
    // --- Create Core Layout Elements ---

    // `renderedView` is the top-level container for the entire chat system.
    const renderedView = document.createElement("div");
    renderedView.className = 'chat-wrapper';

    // `cardContainer` will hold the list of all user cards (the sidebar).
    const cardContainer = document.createElement("div");
    cardContainer.className = 'user-list-container';
    
    // `messageArea` is the container that will hold the active chat window.
    // It's initially hidden by CSS.
    const messageArea = document.createElement('div');
    messageArea.className = 'message-area-container';
    
    // Assemble the layout.
    renderedView.append(cardContainer, messageArea);

    /**
     * Updates the user list UI based on fresh data from the server.
     * This function is called by `app.js` whenever a `user_list_update` event is received.
     * @param {Array<Object>} users - An array of user objects from the backend.
     */
    function updateUserList(users) {
        
        // 1. Sort the array of users according to the required criteria.
        users.sort((a, b) => {
            // First, by last message time (newest first). A null time is treated as 0.
            const timeA = a.lastMessageTime ? new Date(a.lastMessageTime).getTime() : 0;
            const timeB = b.lastMessageTime ? new Date(b.lastMessageTime).getTime() : 0;
            if (timeA !== timeB) {
                return timeB - timeA;
            }
            // If times are the same (e.g., new users), sort alphabetically by username.
            return a.username.localeCompare(b.username); 
        });

        // 2. Clear the existing list to prevent duplicates on re-render.
        cardContainer.innerHTML = ''; 

        // 3. Iterate over the sorted user data to create and append a card for each user.
        users.forEach((user) => {
            // Don't create a card for ourself in the list.
            if (user.userID === window.myUserID) return;

            // Create the main card element.
            const userCard = document.createElement("div");
            userCard.className = 'user-card';
            // Store the user's ID in a data attribute for easy access later.
            userCard.dataset.userId = user.userID; 

            // Create the HTML for the Online/Offline status text based on the `isOnline` flag.
            const statusHTML = user.isOnline 
                ? `<span class="user-status online">Online</span>`
                : `<span class="user-status offline">Offline</span>`;

                  // ---- NEW: Create the unread badge conditionally ----
            const unreadBadgeHTML = user.unreadCount > 0
                ? `<span class="unread-badge">${user.unreadCount}</span>`
                : '';

            // Use a template literal to construct the inner HTML for the card efficiently.
            userCard.innerHTML = `
                <div class="user-card-avatar ${user.isOnline ? 'online' : ''}">
                    <div class="initials">${user.username.charAt(0).toUpperCase()}</div>
                    <div class="status-dot"></div>
                </div>
                <div class="user-card-info">
                    <strong class="username">${user.username}</strong>
                    <p class="last-message">${user.lastMessageContent || 'Click to chat'}</p>
                    ${statusHTML}
                     ${unreadBadgeHTML} <!-- Add the badge here -->
                </div>
            `;
            
            // Add a click listener to each card to open the chat with that user.
            userCard.addEventListener("click", () => openChatWith(user));
            
            // Append the newly created card to the sidebar container.
            cardContainer.appendChild(userCard);
        });
    }

    /**
     * Handles the logic for opening a private chat with a selected user.
     * It manipulates CSS classes to show the message panel and instantiates a new MessageCarriers component.
     * @param {Object} user - The user object associated with the card that was clicked.
     */
   function openChatWith(user) {
        // ---- START: CORRECTED LOGIC ----

        // 1. Find the specific card that was clicked.
        const cardElement = document.querySelector(`.user-card[data-user-id='${user.userID}']`);
        if (!cardElement) return;

        // 2. Visually highlight the currently active chat in the sidebar.
        // First, remove 'active' from any other card.
        document.querySelectorAll('.user-card.active').forEach(c => c.classList.remove('active'));
        // Then, add 'active' to the clicked card.
        cardElement.classList.add('active');

        //Hide the last message preview *only* on the active card.
        const lastMessageElement = cardElement.querySelector(".last-message");
        if (lastMessageElement) {
            lastMessageElement.style.display = "none";
        }

        // 3. Create the new chat window component. This is the most important step.
        // The constructor for MessageCarriers will set up all the necessary global handlers
        // like `window.handleHistoryResponse` BEFORE any data is requested.
        // This resolves the race condition.
        const chatWindow = MessageCarriers(user.userID, user.username);

        // 4. Now that the new component is ready, clear the old one from the DOM
        // and append the new one.
        messageArea.innerHTML = '';
        messageArea.appendChild(chatWindow);

        // 5. Finally, apply the CSS classes to trigger the layout transition.
        cardContainer.classList.add('chat-active');
        messageArea.classList.add('chat-active');

        // ---- END: CORRECTED LOGIC ----
    }
    // Expose the public API of this component.
    return { renderedView, updateUserList };
};