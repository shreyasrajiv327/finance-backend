package repository

import (
	"finance-backend/internal/models"
	"database/sql"
	"fmt"
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

func (r *RecordRepository) GetRecordsByUserID(userID int) ([]models.Record, error) {
    query := `
        SELECT id, user_id, amount, type, category, notes, date, created_at
        FROM records
        WHERE user_id = $1
        ORDER BY date DESC
    `

    rows, err := r.DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var records []models.Record

    for rows.Next() {
        var record models.Record
        err := rows.Scan(
            &record.ID,
            &record.UserID,
            &record.Amount,
            &record.Type,
            &record.Category,
            &record.Notes,
            &record.Date,
            &record.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        records = append(records, record)
    }

    return records, nil
}

func (r *RecordRepository) GetRecordByID(id int) (*models.Record, error) {
    query := `
        SELECT id, user_id, amount, type, category, notes, date, created_at
        FROM records
        WHERE id = $1
    `

    var record models.Record
    err := r.DB.QueryRow(query, id).Scan(
        &record.ID,
        &record.UserID,
        &record.Amount,
        &record.Type,
        &record.Category,
        &record.Notes,
        &record.Date,
        &record.CreatedAt,
    )

    if err != nil {
        return nil, err
    }

    return &record, nil
}


func (r *RecordRepository) DeleteRecord(id int) error {
    _, err := r.DB.Exec("DELETE FROM records WHERE id = $1", id)
    return err
}


func (r *RecordRepository) UpdateRecord(record *models.Record) error {
	query := `
		UPDATE records
		SET amount = $1, type = $2, category = $3, notes = $4, date = $5
		WHERE id = $6
	`

	_, err := r.DB.Exec(query,
		record.Amount,
		record.Type,
		record.Category,
		record.Notes,
		record.Date,
		record.ID,
	)

	return err
}

func (r *RecordRepository) GetFilteredRecords(
	userID int,
	recordType string,
	category string,
	startDate string,
	endDate string,
	limit int,
	offset int,
) ([]models.Record, error) {

	query := `
		SELECT id, user_id, amount, type, category, notes, date, created_at
		FROM records
		WHERE user_id = $1
	`

	args := []interface{}{userID}
	argIndex := 2

	if recordType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, recordType)
		argIndex++
	}

	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	if startDate != "" {
		query += fmt.Sprintf(" AND date >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}

	if endDate != "" {
		query += fmt.Sprintf(" AND date <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY date DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Record

	for rows.Next() {
		var record models.Record
		err := rows.Scan(
			&record.ID,
			&record.UserID,
			&record.Amount,
			&record.Type,
			&record.Category,
			&record.Notes,
			&record.Date,
			&record.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}