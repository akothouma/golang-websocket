import {profile} from './imageComponent/image.js'
import MessageHistory from './messageHistory/message.js'

export const Card=()=>{
const onecard=document.createElement("div");
onecard.style.width='fit-content';
onecard.style.height='500px';
onecard.style.display='flex';
onecard.style.flexDirection='row'
const imageside=document.createElement('div');
imageside.appendChild(profile);
const lastMessageSide=document.createElement('div')
const lastMessage=document.createElement('p');
lastMessageSide.appendChild(lastMessage);
onecard.append(imageside,lastMessageSide);

const {chatHistory,addMessage}=MessageHistory();
chatHistory.style.display='none';
addMessage("Hello")//bound to be dynamic
onecard.appendChild(chatHistory)

return onecard

}