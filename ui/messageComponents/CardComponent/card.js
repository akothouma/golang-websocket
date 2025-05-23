import {Image} from './imageComponent/image.js'
import MessageHistory from './messageHistory/message.js'

const data = [
    { senderId: "user1", messageContent: "Hey, how are you?" },
    { senderId: "user2", messageContent: "I’m good! Just working on the project." },
    { senderId: "user 5", messageContent: "Nice! Need any help?" },
    { senderId: "user3", messageContent: "Anyone up for a quick call?" },
    { senderId: "user4", messageContent: "Sure, I’m free now." },
    { senderId: "user7", messageContent: "Let’s do it!" }
];

 export const Card = () => {
    const cardContainer = document.createElement("div")
    data.forEach((onemessage) => {
        const onecard = document.createElement("div");
        onecard.style.width = 'fit-content';
        onecard.style.height = '500px';
        onecard.style.display = 'flex';
        onecard.id = `${onemessage.senderId}`
        onecard.style.flexDirection = 'row'

        const imageside = document.createElement('div');
        imageside.appendChild(Image());

        const lastMessageSide = document.createElement('div')
        const lastMessage = document.createElement('p');
        lastMessage.textContent = onemessage.messageContent
        lastMessageSide.appendChild(lastMessage);

        const { chatContainer, addMessage } = MessageHistory();
        chatHistory.style.display = 'none';
        addMessage("Hello")//bound to be dynamic
        addMessage("Hello from server","right")
        onecard.append(imageside, lastMessageSide, chatContainer);

        cardContainer.appendChild(onecard)

    })
    return cardContainer

}