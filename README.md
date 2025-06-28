# Crypto Sentiment Analyzer

This project analyzes cryptocurrency sentiment from Reddit posts and serves the results over a small Go web application with a simple frontend.

## Prerequisites

- **Go** 1.23 or newer (matching the version in `go.mod`).
- Reddit API credentials: set `REDDIT_CLIENT_ID` and `REDDIT_CLIENT_SECRET` in a `.env` file or in your environment. You can create these credentials at <https://www.reddit.com/prefs/apps>.

## Running the Server

1. Install dependencies:
   ```bash
   go mod download
   ```
2. Start the application:
   ```bash
   go run main.go
   ```
   The server listens on port `8080` by default.

3. (Optional) Test Reddit API connectivity:
   ```bash
   go run main.go test
   ```
   This performs a quick connection check before running the server.

## API Endpoints

- `GET /api/analyze` – analyze sentiment for a specific cryptocurrency using the `coin` query parameter.
- `GET /api/trending` – fetch sentiment for several popular coins.
- `GET /api/health` – basic health check reporting Reddit service status.

## Frontend Overview

The frontend lives in `templates/index.html` with supporting assets in `static/`. It offers an input field to specify a coin and displays the sentiment analysis result from the backend using a small JavaScript snippet.

