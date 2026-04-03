package repository

import (
	"database/sql"
	"finance-backend/internal/models"
)

type UserRepository struct{
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error{
	query := `
	INSERT INTO users (name, email, password, role)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, active
	`
	return r.DB.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.IsActive)
}


func ( r*UserRepository) GetByEmail(email string) (*models.User, error){
	query := `
	SELECT id, name, email, password, role, active, created_At
	FROM users
	WHERE email = $1
	`

	var user models.User

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
	)
	if err != nil{
		return nil,err
	}
    
	return &user, nil
}