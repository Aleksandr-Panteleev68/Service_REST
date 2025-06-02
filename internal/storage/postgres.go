package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"pet-project/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewDB(cfg *config.Config) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxConns)
	return &PostgresStorage{DB: db}, nil
}

func (d *PostgresStorage) Close() error {
	return d.DB.Close()
}

func (s *PostgresStorage) CreateUser(user User) (int, error) {
	log.Printf("Creating user: %s %s", user.FirstName, user.LastName)
	query := `
	INSERT INTO users (first_name, last_name, age, is_married, password)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var id int
	err := s.DB.QueryRowContext(context.Background(), query, user.FirstName, user.LastName, user.Age, user.IsMarried, user.Password).Scan(&id)
	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint
		"idx_user_name" (SQLSTATE 23505)`{
		return 0, fmt.Errorf("user with name %s %s already exists", user.FirstName, user.LastName)	
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
		
	}
	log.Printf("Create user with ID: %d", id)
	return id, nil
}

func (s *PostgresStorage) GetUserByID(id int) (User, error) {
	log.Printf("Fetching user with ID: %d", id)
	query := `
	SELECT id, first_name, last_name, full_name, age, is_married, password
	FROM users
	WHERE id = $1
	`
	var user User
	err := s.DB.QueryRowContext(context.Background(), query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.FullName,
		&user.Age,
		&user.IsMarried,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		return User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	log.Printf("Fetched user: %s %s", user.FirstName, user.LastName)
	return user, nil
}
