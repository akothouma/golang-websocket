A single-page web application (SPA) forum with user authentication, post/comment functionality, and real-time private messaging. Built with Golang, SQLite, Vanilla JavaScript, WebSockets, and HTML/CSS 

🚀 Features
🛂 Registration & Login

    Register using nickname, email, password, age, gender, first/last name.

    Login with nickname or email + password.

    Secure authentication using bcrypt and session cookies.

    Logout from any page.

📝 Posts & Comments

    Users can create and view posts with categories.

    Comments visible on clicking a post.

    Posts displayed in a live feed-style layout.

💬 Private Messaging (Chat)

    Real-time messaging using WebSockets.

    Sidebar showing online/offline users, sorted by last message.

    Load 10 messages at a time, with infinite scroll (debounce/throttle implemented).

    Messages include timestamp and sender nickname.

🌐 Tech Stack

    Backend: Golang, Gorilla WebSocket, bcrypt, UUID, SQLite

    Frontend: HTML, CSS, Vanilla JS

    Database: SQLite (local file)


⚙️ Setup Instructions
1. Clone the Repo
```
git clone https://learn.zone01kisumu.ke/git/lakoth/real-time-forum
cd real-time-forum
```
2. Initialize Go Modules
```
go mod tidy
```
3. Build and Run the Server
```
go run ./cmd/web/
```
Server runs on http://localhost:8080

📌 Key Dependencies

    github.com/gorilla/websocket

    github.com/mattn/go-sqlite3

    golang.org/x/crypto/bcrypt

    github.com/google/uuid
    

🔐 Authentication Fields
Field	Required	Type
Nickname	✅	string
Email	✅	string
Password	✅	string
First Name	✅	string
Last Name	✅	string
Age	✅	integer
Gender	✅	string
📡 WebSocket Overview

    WebSocket endpoint: /ws

    Messages structured as JSON:

{
  "to": "userID",
  "message": "Hello there!"
}

    Server broadcasts to appropriate client(s) with timestamp and sender info.

🧠 Concepts Learned

    Go routines and channels

    WebSockets in Go and JS

    SPA routing with vanilla JS

    DOM manipulation and event throttling

    Authentication and secure cookie/session handling

    SQL queries and database design


📄 License

MIT License © 2025