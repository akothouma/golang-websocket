import Image from '../imageComponent/image.js'
import  { MessageCarriers}  from '../messagesHistoryComponent/message.js'

const data = [
    { senderId: "user1", messageContent: "Hey, how are you?" },
    { senderId: "user2", messageContent: "I’m good! Just working on the project." },
    { senderId: "user3", messageContent: "Nice! Need any help?" },
    { senderId: "user4", messageContent: "You up for a quick call?" },
    { senderId: "user5", messageContent: "Sure, I’m free now." },
    { senderId: "user6", messageContent: "Let’s do it!" }
];



 export const Card = () => {
    const cardContainer = document.createElement("div");
    cardContainer.display='flex';
    cardContainer.flexDirection='column'
    cardContainer.style.justifyContent='space-between'
    cardContainer.style.gap='15px';
   
    
    data.forEach((oneConnection)=>{
        const onecard = document.createElement("div");
        const displayContainer=document.createElement('div');
        const messageContainer=document.createElement('div');
        onecard.style.width = 'fit-content';
        onecard.style.height = 'fit-content';
        onecard.style.display = 'flex';
        onecard.id = oneConnection.senderId;
        onecard.classList.add('card');
        onecard.style.flexDirection = 'column'
        onecard.style.justifyContent='space-between'
        onecard.style.gap='15px';


        displayContainer.style.display='flex';
        displayContainer.flexDirection='row'
        const imageside = document.createElement('div');
        imageside.appendChild(Image());

        const lastMessageSide = document.createElement('div')
        const lastMessage = document.createElement('p');
        lastMessage.textContent =oneConnection.messageContent
        lastMessageSide.appendChild(lastMessage);

        displayContainer.append(imageside,lastMessage)

        const {chatContainer,AddMessage}=MessageCarriers();
        messageContainer.appendChild(chatContainer)
        messageContainer.style.display = 'none';

        onecard.addEventListener("click",()=>{
           
        })
        
        onecard.append(displayContainer,messageContainer);

        cardContainer.appendChild(onecard)
    })
    return cardContainer

}
