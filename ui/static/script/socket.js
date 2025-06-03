let socket = null;

export const initSocket = () => {
    // Return existing socket if already connected
    if (socket && socket.readyState === WebSocket.OPEN) {
        return socket;
    }
    
    // Close existing socket if it exists but not open
    if (socket) {
        socket.close();
    }
    
    socket = new WebSocket("ws://localhost:8000/ws");
    
    socket.addEventListener('open', () => {
        console.log('WebSocket connected successfully');
    });
    
    socket.addEventListener('close', (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason);
        socket = null; // Reset socket reference
    });
    
    socket.addEventListener('error', (error) => {
        console.error('WebSocket error:', error);
    });
    
    return socket;
};

export const getSocket = () => {
    return socket;
};