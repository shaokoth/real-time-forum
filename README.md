# Real-Time-Forum
* This project is a fully functional forum web application built using stack that includes SQLite, Golang, JavaScript, HTML and CSS. The app allows users to register, login, create posts, comment on them and send private messages. It also supports real-time communication through websockets.

## Features
1. User Registration and Login

* Users can register by filling out a form with the following details:
- Nickname
- Age
- Gender 
- FirstName
- LastName
- Email
- Password
* Users can login with either their nickname or email combined with their password.

2. Posts and Comments

* Users can create posts under specific categories
* Users can comment on posts to engage in discussions

3. Private Messages 

* Users can send private messages to other users in real time using websockets
* The application includes: 
   - A list of online users organized by recent messages or alphabetically for new users.
   - A real-time chat section for direct communication
   - Chat history, where users can view the last 10 messages, with the ability to scroll and load more.
* Messages display:
    - The date the message was sent.
    - The username of the sender.

## Project Architecture
1. Frontend(JavaScript, HTML, CSS)
    - Javascript handles all frontend events and WebSocket communication with the backend
    - HTML contains the structure of the page and organizes all necessary elements.
    - CSS  is used to style the elements and create a user-friendly interface
* The application follows the Single Page Application (SPA) architecture.

2. Backend(Golang)
    - Golang is used for handling the backend logic, including managing database operations and Websockets connections.
    - The backend also handles user authentication and message storage.

3. Database(SQLite)
   - SQLite stores all necessary data, including user accounts, posts, comments and private messages

## Installation
### Prerequisites
* Go(Golang) - latest  version
* SQlite 
### Usage

```bash
$ git clone https://learn.zone01kisumu.ke/git/hilaromondi/real-time-forum.git
 
 cd real-time-forum
 ```
 
 ```bash
 $ go run cmd/main.go
 ```
* Open  your browser and go to http://localhost:8080

## Technologies Used
 * Frontend
   - HTML
   - CSS
   - JavaScript(for frontend events and Websockets)

 * Backend 
   - Golang (for server-side logic and WebSocket handling)

 * Database
   - SQLite
 
 * WebSocket
   - For real-time messaging and notifications

## Collaboration
  * Incase of any issue feel free to open an issue.

## Collaborators
  * Hilary Omondi [hilaromondi](https://learn.zone01kisumu.ke/git/hilaromondi)
  
  * Shadrack Okoth [shaokoth](https://learn.zone01kisumu.ke/git/shaokoth)