
let socket = null;

export const initSocket = () => {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
        return socket;
    }
    
    socket = new WebSocket("ws://localhost:8000/ws");
    
    socket.addEventListener('open', () => {
        console.log('WebSocket connected. Requesting initial user list...');
        // FIX: Send a simplified, standardized message to get the initial list.
        socket.send(JSON.stringify({ type: 'get_user_list' }));
    });
    
    socket.addEventListener('close', (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        socket = null; 
    });
    
    socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
    });
    
    return socket;
};