import {Card} from '../messageComponents/CardComponent/card.js';
import { MessageCarriers } from '../messageComponents/messagesHistoryComponent/message.js';

document.addEventListener('DOMContentLoaded',()=>{
    
        const socket=new WebSocket("ws://localhost:8000/ws");
    
        socket.addEventListener("open",()=>{
            const request={
                event:"open",
                payload:{
                       messageType:"get_online_users",
                }
            }
            console.log("server connected succesfully");
            socket.send(JSON.stringify(request));

        })
        const {ShowConnections}=Card()
        socket.addEventListener("message",(e)=>{
            try{
              const data =JSON.parse(e.data);
              const {message,value}=data;

              switch (message){
                case "connected_client_list":
                    ShowConnections(value);
              }



            }catch(error){

            }
        })
    
    
    const root=document.getElementById("message_layout");
    root.style.height='100px';
    root.style.display='flex';
    root.style.flexDirection='column';
    root.style.background='var(--color-white)';
    root.style.padding='var(--card-padding)';
    root.style.borderRadius='var(--card-border-radius)';
    root.style.fontSize='1.4rem';
    root.style.height='max-content';  
    const renderedView=Card();
    root.appendChild(renderedView);

    const {AddMessage}=MessageCarriers();
    
})