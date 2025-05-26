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
        msg.textContent=mess
        chat.appendChild(msg)
    }
    const messageInput=document.createElement('input');
    messageInput.placeholder='Type here...';
    messageInput.style.position='sticky';
    messageInput.style.bottom='0';
    chatContainer.append(chat,messageInput);

    AddMessage("hello");
    AddMessage("hello from the server side...At least I can sat that I have connected",'right')
  return {chatContainer} 
  }