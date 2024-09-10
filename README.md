# Ethereum Transaction Tracker

## How to run

You can easily run the project with docker-compose. Just run the following command:

```bash
docker-compose up -d
```

After that, you can check the logs for the transactions using

```bash
docker logs -f bot
```

The telegram webhook will be available at the port specified in the environment variable `WEBHOOK_PORT` (default is 8000).

**In order to activate the webhook, you need to send a request to telegram api to update the webhook.**

```bash
curl -F "url=<YOUR_URL_HERE>/webhook" https://api.telegram.org/bot<YOUR_BOT_TOKEN>/setWebhook
```

After that, you can start sending messages to the bot.

The bot is available at `@EthTBot` on telegram.

That's it! Enjoy!
