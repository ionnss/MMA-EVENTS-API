// internal/storage/db.go
package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "internal/data/mma_events.db")
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o banco de dados: %w", err)
	}

	// Criação da tabela organizations
	query := `
	CREATE TABLE IF NOT EXISTS organizations (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	name TEXT NOT NULL UNIQUE,
    	url TEXT NOT NULL,
		eventurl TEXT NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela de organizações: %w", err)
	}

	// Criação da tabela events
	createEventsTable := `
    CREATE TABLE IF NOT EXISTS events (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		org TEXT,
        event_name TEXT,
        date_main_card TEXT,
        date_prelims_card TEXT,
        location TEXT,
        city TEXT,
        state TEXT,
        country TEXT
    );`
	if _, err := db.Exec(createEventsTable); err != nil {
		return nil, fmt.Errorf("erro ao criar tabela de eventos: %w", err)
	}

	return db, nil
}

func InsertEvent(db *sql.DB, org, eventName, dateMainCard, datePrelimsCard, location, city, state, country string) error {
	insertQuery := `
    INSERT INTO events (org, event_name, date_main_card, date_prelims_card, location, city, state, country)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.Exec(insertQuery, org, eventName, dateMainCard, datePrelimsCard, location, city, state, country)
	if err != nil {
		return fmt.Errorf("erro ao inserir evento %s: %w", eventName, err)
	}

	return nil
}
