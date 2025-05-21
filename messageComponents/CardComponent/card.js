import {profile} from './imageComponent/image.js'

export const card=()=>{
const onecard=document.createElement("div");
onecard.style.width='fit-content';
onecard.style.height='500px';
onecard.style.display='flex';
onecard.style.flexDirection='row'
const imageside=document.createElement('div');
const lastMessageSide=document.createElement('div')
onecard.append(imageside,lastMessageSide);
imageside.appendChild(profile);
const lastMessage=document.createElement('p');
lastMessageSide.appendChild(lastMessage);

return onecard

}