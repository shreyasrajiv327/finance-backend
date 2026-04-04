package repository

import (
	"finance-backend/internal/models"
	"database/sql"
)

type RecordRepository struct {
	DB *sql.DB
}

func (r *RecordRepository) Create(record *models.Record) error {
	query := `
		INSERT INTO records (user_id, amount, type, category, date, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id, created_at
	`

	return r.DB.QueryRow(
		query,
		record.UserID,
		record.Amount,
		record.Type,
		record.Category,
		record.Date,
		record.Notes,
	).Scan(&record.ID, &record.CreatedAt)
}