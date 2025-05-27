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

    function showConnections(data) {
       
        cardContainer.innerHTML = '';

        data.forEach((oneConnection) => {
            const onecard = document.createElement("div");
            const displayContainer = document.createElement('div');

            onecard.style.width = 'fit-content';
            onecard.style.height = 'fit-content';
            onecard.style.display = 'flex';
            onecard.id = oneConnection.senderId;
            onecard.classList.add('card');
            onecard.style.flexDirection = 'column';
            onecard.style.alignContent = 'space-between';
            onecard.style.gap = '15px';
            onecard.style.cursor = 'pointer';

            displayContainer.style.display = 'flex';
            displayContainer.style.flexDirection = 'row';

            const imageside = document.createElement('div');
            imageside.appendChild(Image());

            const lastMessageSide = document.createElement('div');
            const lastMessage = document.createElement('p');
            lastMessage.textContent = oneConnection.messageContent;
            lastMessage.style.fontSize = '12px';
            lastMessageSide.appendChild(lastMessage);

            displayContainer.append(imageside, lastMessageSide);

            onecard.appendChild(displayContainer);
            onecard.addEventListener("click", () => {
                showPrivateMessages(oneConnection.senderId, cardContainer);
            });

            cardContainer.appendChild(onecard);
        });
    }

    return {
        renderedView,
        showConnections
    };
}

function showPrivateMessages(senderId, cardsView) {
    const renderedView = cardsView.parentNode;
    const messageContainer = document.createElement('div');
    messageContainer.classList.add('message_container');
    messageContainer.id = senderId;

    const backButton = document.createElement('button');
    backButton.textContent = "← Back";
    backButton.addEventListener("click", () => {
        renderedView.innerHTML = '';
        renderedView.appendChild(cardsView);
    });

    const { chatContainer } = MessageCarriers();
    chatContainer.id=senderId;
    messageContainer.append(backButton, chatContainer);
    renderedView.innerHTML = '';
    renderedView.appendChild(messageContainer);
}
