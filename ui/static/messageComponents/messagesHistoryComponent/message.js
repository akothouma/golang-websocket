 import { initSocket } from "../../script/socket.js";

 export const MessageCarriers=()=>{
    const chatContainer=document.createElement('div');
    chatContainer.className='message_history'

    const chat=document.createElement('div');
    chat.style.flex='1';
    chat.style.overflowY='auto';

     function AddMessage(mess,side='left'){
        const msg=document.createElement('p');
        msg.style.fontSize='12px';
        msg.classList.add(`msg${side}`);
        msg.textContent=mess;
        msg.style.justifyContent='center';
        chat.appendChild(msg)
        
    }
    const messageInput=document.createElement('input');
    messageInput.placeholder='Type here...';
    messageInput.style.position='sticky';
    messageInput.style.bottom='0';
    chatContainer.append(chat,messageInput);


   const socket=initSocket();
   messageInput.addEventListener('keydown',(e)=>{
      if (e.key=="Enter"){
         const message=e.target.value;
         if (socket && socket.readyState==WebSocket.OPEN){

            const request={
               event:"sending_message",
               payload:{
                  messageType:"chat_message",
                  receiverID:chatContainer.id,
                  content:message,
               }
            }
            socket.send(JSON.stringify(request))
         }
         AddMessage(message)
         messageInput.value=''
      }
   })
  return {chatContainer,AddMessage} 
  }