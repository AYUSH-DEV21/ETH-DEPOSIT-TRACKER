package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"luganodes/internal/database"
	"net/http"
	"os"
)

type Notifier struct {
	Token string
}

func NewNotifier() *Notifier {
	return &Notifier{
		Token: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
}

func (n *Notifier) NewTx(tx database.Transaction) {
	chatIds := database.GetChatIDs()
	for _, chatID := range chatIds {
		message := formatTransactionMessage(tx)
		err := sendTelegramMessage(n.Token, chatID, message)
		if err != nil {
			log.Printf("Error sending Telegram message: %v", err)
		}
	}
}

func formatTransactionMessage(tx database.Transaction) string {
	return fmt.Sprintf("New transaction:\n"+
		"From: %s\n"+
		"To: %s\n"+
		"Value: %f\n"+
		"Block Number: %s\n"+
		"Hash: %s",
		tx.From, tx.To, tx.Value, tx.BlockNumber, tx.Hash)
}

func sendTelegramMessage(token, chatID, message string) error {
	url := "https://api.telegram.org/bot" + token + "/sendMessage"
	payload := map[string]string{
		"chat_id": chatID,
		"text":    message,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Telegram API error: %s", resp.Status)
	}
	return nil
}

func (n *Notifier) HandleUpdate(update []byte) {
	var updateData struct {
		Message struct {
			Chat struct {
				ID int64 `json:"id"`
			} `json:"chat"`
			Text string `json:"text"`
		} `json:"message"`
	}

	err := json.Unmarshal(update, &updateData)
	if err != nil {
		slog.Error("Error parsing update: %v", err)
		return
	}

	switch updateData.Message.Text {
	case "/start":
		chatID := fmt.Sprint(updateData.Message.Chat.ID)

		err := database.AddChatID(chatID)
		if err != nil {
			slog.Info("Chat ID already exists")
			return
		}

		welcomeMessage := "Welcome to the Eth Transaction Tracker Bot! You will now receive notifications for new transactions."
		err = sendTelegramMessage(n.Token, chatID, welcomeMessage)

		if err != nil {
			slog.Error("Error sending welcome message: %v", err)
		}

	case "/txs":
		txs := database.GetLastTx(5)

		if len(txs) == 0 {
			slog.Info("No transactions found")
			return
		}

		chatID := fmt.Sprint(updateData.Message.Chat.ID)

		sendTelegramMessage(n.Token, chatID, "Last 5 transactions:")

		for _, tx := range txs {
			message := formatTransactionMessage(tx)
			err := sendTelegramMessage(n.Token, chatID, message)
			if err != nil {
				slog.Error("Error sending transaction message", err)
			}
		}

	}
}

func (n *Notifier) InitWebhook() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading request body: %v", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		n.HandleUpdate(body)

		w.WriteHeader(http.StatusOK)
	})
	port := os.Getenv("WEBHOOK_PORT")
	slog.Info("Starting webhook server on port",
		slog.String("port", port),
	)
	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error("Error starting webhook server: %v", err)
	}
}
