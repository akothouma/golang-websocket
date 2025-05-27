import Image from '../imageComponent/image.js'
import { MessageCarriers } from '../messagesHistoryComponent/message.js'

// const data = [
//     { senderId: "user1", messageContent: "Hey, how are you?" },
//     { senderId: "user2", messageContent: "I’m good! Just working on the project." },
//     { senderId: "user3", messageContent: "Nice! Need any help?" },
//     { senderId: "user4", messageContent: "You up for a quick call?" },
//     { senderId: "user5", messageContent: "Sure, I’m free now." },
//     { senderId: "user6", messageContent: "Let’s do it!" }
// ];

export const Card = () => {
    const renderedView = document.createElement("div");
    renderedView.style.height = 'fit-content';
    renderedView.style.overflow = 'hidden';

    const cardsView =ShowConnections(data);
    renderedView.appendChild(cardsView);
    return {renderedView,
        ShowConnections};
}

function ShowConnections(data) {
    const cardContainer = document.createElement("div");
    cardContainer.display = 'flex';
    cardContainer.flexDirection = 'column'
    cardContainer.style.alignContent = 'space-between'
    cardContainer.style.gap = '15px';

    //const activecard = null;
    data.forEach((oneConnection) => {
        const onecard = document.createElement("div");
        const displayContainer = document.createElement('div');
        onecard.style.width = 'fit-content';
        onecard.style.height = 'fit-content';
        onecard.style.display = 'flex';
        onecard.id = oneConnection.senderId;
        onecard.classList.add('card');
        onecard.style.flexDirection = 'column'
        onecard.style.alignContent = 'space-between'
        onecard.style.gap = '15px';
        onecard.style.cursor = 'pointer'


        displayContainer.style.display = 'flex';
        displayContainer.flexDirection = 'row'
        const imageside = document.createElement('div');
        imageside.appendChild(Image());

        const lastMessageSide = document.createElement('div')
        const lastMessage = document.createElement('p');
        lastMessage.textContent = oneConnection.messageContent
        lastMessage.style.fontSize='12px';
        lastMessageSide.appendChild(lastMessage);

        displayContainer.append(imageside,lastMessage)

        onecard.appendChild(displayContainer);
        onecard.addEventListener("click", () => {
            showPrivateMessages(oneConnection.senderId, cardContainer);
        });
        cardContainer.appendChild(onecard);

    });
    return cardContainer;
}

function showPrivateMessages(senderId, cardsView) {
    const renderedView = cardsView.parentNode;
    const messageContainer = document.createElement('div');
    messageContainer.classList.add('message_container');
    messageContainer.id = senderId;
    const backButton = document.createElement('button');
    backButton.textContent="← Back";
    backButton.addEventListener("click", () => {
        renderedView.innerHTML = '';
        renderedView.appendChild(cardsView)
    });

    const { chatContainer } = MessageCarriers();
    messageContainer.append(backButton, chatContainer)
    renderedView.innerHTML = '';
    renderedView.appendChild(messageContainer)

}
