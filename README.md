# Godoku

Godoku is a simple web-based Sudoku application built with Go and HTMX. This project aims to demonstrate dynamic web application capabilities using Go on the backend and HTMX for front-end interactivity.

## Technologies Used

- **Go**: Backend language for game logic and server-side rendering.
- **HTMX**: For dynamic and interactive front-end components without the need for a full front-end framework.
- **HTML/CSS**: Basic styling and layout for the Sudoku board.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/AbsoluteZero000/godoku.git
   cd godoku
   ```

2. **Install dependencies:**

   Ensure you have Go installed on your machine. Then, run the project:

   ```bash
   go run cmd/main.go 
   ```

3. **Access the application:**

   Open your web browser and navigate to `http://localhost:8080` to start playing Sudoku.

## Usage

- **Starting a Game**: The application will display a Sudoku board when you first access it.
- **Inputting Numbers**: Click on any cell to input a number (1-9).
- **Validation**: The backend automatically validates inputs, ensuring the Sudoku rules are respected.

