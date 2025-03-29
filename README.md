# Web Forum Project Documentation

## Overview

This project is a web-based forum that enables communication between users through posts and comments. It allows user authentication, post categorization, interaction via likes/dislikes, and filtering functionalities. The back-end is built using Go, with SQLite as the database, and the project is containerized using Docker.

## Features

1. ### User Authentication



- Users can register and log in to the forum.

- Registration requires:

    - A unique email (duplicate emails return an error).

    - A username.

    - A password (bonus: encrypted storage using bcrypt).

- Login session management with cookies (session expiration included).

- Each user can only have one active session at a time.

- Users can upload profile pictures.

2. ### Posts and Comments



- Only registered users can create posts and comments.

- Posts can be categorized (customizable categories).

- Users can attach media (images, videos, GIFs) to posts.

- All users (registered or not) can view posts and comments.

3. ### Likes and Dislikes

- Only registered users can like/dislike posts and comments.

- Like/dislike counts are visible to all users.

4. ### Filtering Mechanism

    Users can filter displayed posts by:

    - Categories (subforums based on topics).

    - Created posts (for logged-in users only).

    - Liked posts (for logged-in users only).

## Technology Stack

- #### Backend: Go (Golang)

- #### Database: SQLite

- #### Authentication: Sessions and cookies

- #### Containerization: Docker

- #### Security: Password hashing with bcrypt (bonus feature)

## Database Schema

### Tables



- users (id, email, username, password, session_uuid, session_expiry, -profile_picture, content_type)

- posts (id, user_id, title, content, category, media, content_type, created_at)

- comments (id, post_id, user_id, content, created_at)

- likes (id, user_id, post_id, comment_id, type)

## Example Queries

- ### Create Table:


    ```
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        email TEXT UNIQUE NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
        session_uuid TEXT,
        session_expiry DATETIME,
        profile_picture BLOB,
        content_type TEXT
    );
    ```

- ### Insert Data:

    ```
    INSERT INTO users (email, username, password) VALUES ('user@example.com', 'username', 'hashed_password');
    ```

- ### Select Data:

    ```
    SELECT * FROM posts WHERE category = 'Technology';
    ```


## API Endpoints

### User Authentication

- POST /register – User registration.

- POST /login – User login.

- GET /logout – End user session.

- POST /upload-profile – Upload user profile picture.

### Posts & Comments

- GET /posts – Retrieve all posts.

- POST /posts – Create a new post (authenticated users only).

- GET /posts/{id} – Retrieve a specific post.

- POST /posts/{id}/comment – Add a comment to a post.

- GET /comments/{id}/replies – Retrieve all replies for a comment.

### Likes & Dislikes

- POST /posts/{id}/like – Like a post.

- POST /posts/{id}/dislike – Dislike a post.

### Filtering

- GET /posts?category=Tech – Filter posts by category.

- GET /user/posts – Get posts created by the logged-in user.

- GET /user/liked – Get posts liked by the logged-in user.

## Error Handling

- 400 Bad Request: Invalid input or missing parameters.

- 401 Unauthorized: User not logged in.

- 403 Forbidden: User does not have permission.

- 404 Not Found: Resource does not exist.

- 500 Internal Server Error: Unexpected server failure.

## Deployment & Docker Setup

### Dockerfile
```
FROM golang:1.18
WORKDIR /app
COPY . .
RUN go build -o forum
CMD ["./forum"]
EXPOSE 8080
```
## Running the Project

**1. Build Docker Image:**
    
```
docker build -t forum-app .
```

**2. Run the Container:**
```
docker run -p 8080:8080 forum-app
```
**3. Testing**

- Unit tests are implemented for critical functions.

- Use httptest for handler testing.

- Run tests with:
    ```
    go test ./...
    ```
## Conclusion

This forum project demonstrates core web development concepts such as authentication, database interactions, and API design. It is built with Go, uses SQLite for storage, and follows best practices in handling errors and security.

