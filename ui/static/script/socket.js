    let socket;
    export const initSocket=()=>{
        socket = new WebSocket("ws://localhost:8000/ws");
        return socket   
    }
  document.addEventListener("DOMContentLoaded",()=>{
      initSocket();
  })
