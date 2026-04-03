package database

import "log"

func Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		active BOOLEAN DEFAULT true,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS records (
		id SERIAL PRIMARY KEY,
		user_id INT REFERENCES users(id),
		amount NUMERIC NOT NULL,
		type TEXT CHECK (type IN ('income', 'expense')),
		category TEXT,
		note TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	log.Println("Tables created")
}