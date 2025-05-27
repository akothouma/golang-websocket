export function initSocket(){
    const socket = new WebSocket("ws://localhost:8000/ws");
    return  socket
}