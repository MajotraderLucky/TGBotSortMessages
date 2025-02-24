# Telegram Message Bot

A Telegram bot built using [gotd](https://github.com/gotd/td) that retrieves unread messages (excluding channels) and saves them in JSON and Markdown formats.

## Features
- **Automatic Telegram Login**: Uses phone number authentication.
- **Fetch Unread Messages**: Retrieves unread messages from private chats and groups.
- **Storage**: Saves messages as:
  - `unread.json` (formatted JSON file)
  - `latest_messages.md` (Markdown file with the last 50 messages)
- **Scheduled Execution**: Runs every minute to check for new messages.

## Installation
### Prerequisites
- Go 1.18+
- Telegram API credentials (API_ID and API_HASH)
- A valid phone number linked to a Telegram account

### Setup
1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/telegram-message-bot.git
   cd telegram-message-bot
   ```
2. Create a `.env` file with your Telegram API credentials:
   ```sh
   API_ID=your_api_id
   API_HASH=your_api_hash
   ```
3. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage
To run the bot, execute:
```sh
go run main.go
```
The bot will log into Telegram and start checking messages every minute.

## File Structure
```
.
├── internal
│   ├── auth            # Handles authentication
│   ├── bot             # Core bot logic
│   ├── config          # Loads environment variables
│   ├── storage         # Saves messages in JSON and Markdown
│   ├── telegram        # Telegram API integration
├── main.go             # Entry point
├── go.mod              # Go modules file
├── go.sum              # Dependencies checksum
└── README.md           # Project documentation
```

## How It Works
1. Loads the Telegram API credentials from the `.env` file.
2. Logs in using the provided phone number.
3. Retrieves unread messages from personal chats and groups.
4. Saves the messages in `messages/unread.json` and `messages/latest_messages.md`.
5. Runs a loop that checks for new messages every minute.

## Contributions
Feel free to fork and submit pull requests! Any improvements are welcome.

## License
This project is licensed under the MIT License.


