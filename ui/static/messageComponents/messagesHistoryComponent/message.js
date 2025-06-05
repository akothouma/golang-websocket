 import { initSocket } from "../../script/socket.js";

 export const MessageCarriers=()=>{
    const chatContainer=document.createElement('div');
    chatContainer.className='message_history'

    const chat=document.createElement('div');
    chat.style.flex='1';
    chat.style.overflowY='auto';
    chat.style.maxHeight = '400px';
    chat.style.padding = '10px';
    chat.style.border = '1px solid #ccc';
    chat.style.marginBottom = '10px';

     function AddMessage(mess,side='left',senderId=null){
          const msg = document.createElement('p');
        msg.style.fontSize = '12px';
        msg.style.margin = '5px 0';
        msg.style.padding = '8px 12px';
        msg.style.borderRadius = '12px';
        msg.style.maxWidth = '70%';
        msg.style.wordWrap = 'break-word';
        
        if (side === 'right') {
            // Sender's message (right side)
            msg.style.backgroundColor = '#007bff';
            msg.style.color = 'white';
            msg.style.marginLeft = 'auto';
            msg.style.textAlign = 'right';
        } else {
            // Receiver's message (left side)
            msg.style.backgroundColor = '#f1f1f1';
            msg.style.color = 'black';
            msg.style.marginRight = 'auto';
            msg.style.textAlign = 'left';
        }
        
        msg.textContent = mess;
        chat.appendChild(msg);
        
        // Auto-scroll to bottom
        chat.scrollTop = chat.scrollHeight;
        
        console.log(`Message added: ${mess} (side: ${side})`);
        
    }
      const messageInput = document.createElement('input');
    messageInput.placeholder = 'Type here...';
    messageInput.style.position = 'sticky';
    messageInput.style.bottom = '0';
    messageInput.style.width = '100%';
    messageInput.style.padding = '10px';
    messageInput.style.border = '1px solid #ccc';
    messageInput.style.borderRadius = '5px';

    chatContainer.append(chat, messageInput);


//    const socket=initSocket();
messageInput.addEventListener('keydown',(e)=>{
    if (e.key=="Enter"){
        e.preventDefault();
        const message=e.target.value.trim();
        if (!message){
            console.log("Empty message, not sending");
            return;
        } 

         // Use the global socket instead of creating a new one
        const socket=window.globalSocket;
            
         if (socket && socket.readyState==WebSocket.OPEN){
            const receiverId=chatContainer.id;
              if (!receiverId) {
                    console.error("No receiver ID set for chat container");
                    return;
                }
            const request={
               event:"sending_message",
               payload:{
                  messageType:"chat_message",
                  receiverID:receiverId,
                  content:message,
               }
            }
            // console.log("Sending message:", request);
            socket.send(JSON.stringify(request))
            //Add message to sender's view immediately
            AddMessage(message);

            // Clear input after sending
            messageInput.value=''
         }else{
             console.error("WebSocket not connected");
         }
      }
   })
  return {chatContainer,AddMessage} 
  }