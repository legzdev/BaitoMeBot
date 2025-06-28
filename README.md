# BaitoMeBot - Telegram File Stream Bot

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev)
[![Telegram](https://img.shields.io/badge/Telegram-Bot-2CA5E0?logo=telegram)](https://telegram.org)

BaitoMeBot is a Telegram bot to generate direct and streamable links to your Telegram files.

---

## ⚙ Configuration

#### Bot Commands
```
start - start the bot
txt - toggle txt mode
txtcaption - /txt but uses the first caption line as filename
txtcaptionfull - /txt but uses the full caption as filename
help - show bot help
```


#### Environment variables
Variables without a default value are required.

| Name | Default Value | Description |
| --- | --- | --- |
| `TELEGRAM_API_HASH` | | This is the API hash for your Telegram account, which can be obtained from https://my.telegram.org |
| `TELEGRAM_API_ID` | | This is the API ID for your Telegram account, which can be obtained from https://my.telegram.org |
| `TELEGRAM_BOT_TOKEN` | | This is the bot token for your bot, which can be obtained from [@BotFather](https://t.me/BotFather) |
| `TELEGRAM_CHAT_ID` | | This is the chat_id for the chat where the bot stores these files to make the links works |
| `ALLOWED_USERS` | nil | A list of user IDs separated by comma (,). If this is set, only the users in this list or in the auth channel will be able to use the bot |
| `AUTH_CHANNEL_ID` | 0 | This is the chat_id for the auth channel. If this is set, only the users in this channel (admins and members) or in ALLOWED_USERS will be able to use the bot |
| `SERVER_ADDDESS` | :8080 | This is the address that the server will listen to |
| `SERVER_HOST` | http://<your ip>:8080 | A [FQDN]() or [IP address] of the server where the bot is running |
| `HASH_LENGTH` | 6 | Hash length for generated URLs. The hash length must be greater than 5 and less than or equal to 32 |


An example `.env` file
```ini
TELEGRAM_API_HASH=foobar
TELEGRAM_API_ID=12345
TELEGRAM_BOT_TOKEN=567890:abcdef
TELEGRAM_CHAT_ID=-10012345678
SERVER_HOST=http://example.org
```

---

##  🛠Contributing
Contributions are welcome! Open an issue or submit a pull request.

---

## 📜 License
This project is licensed under the **GNU Affero General Public License v3.0**.
See the full terms in [LICENSE](LICENSE) file.

