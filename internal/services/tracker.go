package services

import (
	"encoding/json"
	"io"
	"log/slog"
	"luganodes/internal/database"
	"net/http"
	"os"
	"strings"
	"time"
)

type Tracker struct {
	Address  string
	NumTx    string
	Interval time.Duration
	Notifier *Notifier
	ApiKey   string
}

func NewTracker(address string, numTx string, interval time.Duration, notifier *Notifier) *Tracker {
	return &Tracker{
		Address:  address,
		NumTx:    numTx,
		Interval: interval,
		Notifier: notifier,
		ApiKey:   os.Getenv("ALCHEMY_API_KEY"),
	}
}

func (t *Tracker) FindOrCreateTx(tx database.Transaction) {
	_, err := database.FindTxByHash(tx.Hash)

	if err == nil {
		return
	}

	slog.Info("New transaction: ",
		slog.String("hash", tx.Hash),
		slog.String("from", tx.From),
		slog.String("to", tx.To),
		slog.Float64("value", tx.Value))

	err = database.CreateTx(tx)

	if err != nil {
		slog.Error("Error creating transaction: ", err)
	}

	t.Notifier.NewTx(tx)
}

func (t *Tracker) FetchTx() {
	requestBody := `{
		"jsonrpc": "2.0",
		"id": 1,
		"method": "alchemy_getAssetTransfers",
		"params": [
			{
				"fromBlock": "0x0",
				"toBlock": "latest",
				"toAddress": "` + t.Address + `",
				"category": ["external"],
				"withMetadata": true,
				"maxCount": "` + t.NumTx + `",
				"order": "desc"
			}
		]
	}`

	res, err := http.Post(
		"https://eth-mainnet.alchemyapi.io/v2/"+t.ApiKey,
		"application/json",
		strings.NewReader(requestBody),
	)

	if err != nil {
		slog.Error("Error fetching transactions: ", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		slog.Error("Error fetching transactions: ", res.Status)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response body: ", err)
		return
	}

	var response struct {
		Result struct {
			Transfers []database.Transaction `json:"transfers"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Error unmarshalling response: ", err)
		return
	}

	for _, tx := range response.Result.Transfers {
		t.FindOrCreateTx(tx)
	}
}

func (t *Tracker) Start() {
	for {
		t.FetchTx()
		time.Sleep(t.Interval)
	}
}
