export default  MessageCarriers=()=>{
    const chatContainer=document.createElement('div');
    chatContainer.className='message_history'

     function addMessage(mess,side='left'){
        const msg=document.createElement('p');
        msg.className=`msg ${side}`;
        msg.textContent=mess
        chatContainer.appendChild(msg)
    }
    const messageInput=document.createElement('input');
    messageInput.placeholder='Type here...';
    messageInput.style.bottom='0';
    chatContainer.appendChild(messageInput);
  return {chatContainer ,addMessage} 
}