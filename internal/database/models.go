package database

import "time"

type Transaction struct {
	UniqueID    string    `json:"uniqueId" gorm:"primary_key"`
	Hash        string    `json:"hash" gorm:"unique"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Value       float64   `json:"value"`
	BlockNumber string    `json:"blockNum"`
	CreatedAt   time.Time `gorm:"default:current_timestamp"`
}

func FindTxByHash(hash string) (Transaction, error) {
	var tx Transaction
	err := db.Where("hash = ?", hash).First(&tx).Error
	return tx, err
}

func CreateTx(tx Transaction) error {
	err := db.Create(&tx).Error
	return err
}

func GetLastTx(numTx int) []Transaction {
	var txs []Transaction
	db.Order("created_at desc").Limit(numTx).Find(&txs)
	return txs
}

type Chats struct {
	ChatID string `json:"chatId" gorm:"primary_key"`
}

func GetChatIDs() []string {
	var chats []Chats
	db.Find(&chats)
	chatIDs := make([]string, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ChatID
	}
	return chatIDs
}

func AddChatID(chatID string) error {
	chat := Chats{ChatID: chatID}
	err := db.Create(&chat).Error

	return err
}
