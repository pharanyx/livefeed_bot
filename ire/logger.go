package ire

import (
	"log"

	"github.com/adayoung/ada-bot/utils/storage"
)

func init() {
	storage.OnReady(initDB)
}

var initDBComplete bool = false

func initDB() {
	sqlTable := `
		CREATE TABLE IF NOT EXISTS "ire_gamefeed" (
			"id" integer NOT NULL PRIMARY KEY,
			"caption" varchar(256) NOT NULL,
			"description" varchar(256) NOT NULL,
			"type" varchar(5) NOT NULL,
			"date" timestamp NOT NULL
		);
	`
	if _, err := storage.DB.Exec(sqlTable); err == nil {
		initDBComplete = true
	} else {
		log.Printf("error: %v", err) // We won't store events, that's what!
	}
}

func logEvent(e Event) {
	if !initDBComplete {
		return // We're not ready to save events
	}

	eventLog := `INSERT INTO ire_gamefeed (
		id, caption, description, type, date
		) VALUES (?, ?, ?, ?, ?)
	`
	eventLog = storage.DB.Rebind(eventLog)
	if _, err := storage.DB.Exec(eventLog, e.ID, e.Caption,
		e.Description, e.Type, e.Date); err != nil {
		log.Printf("error: %v", err)
	}
}
