package storage

import (
	"context"
	"fmt"
	"time"

	"pet-project/internal/config"
	"pet-project/internal/domain"
	"pet-project/internal/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool   *pgxpool.Pool
	logger *logger.Logger
}

func NewDB(ctx context.Context, cfg *config.Config, logger *logger.Logger) (*PostgresStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.MaxConns = int32(cfg.Database.MaxConns)
	config.MinConns = int32(cfg.Database.MaxConns / 2)
	config.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Database connection established", "host", cfg.Database.Host, "port", cfg.Database.Port)
	return &PostgresStorage{pool: pool, logger: logger}, nil
}

func (s *PostgresStorage) Close() {
	s.pool.Close()
	s.logger.Info("Database connection closed")
}

func (s *PostgresStorage) CreateUser(ctx context.Context, user domain.User) (int64, error) {
	s.logger.Info("Creating user", "first_name", user.FirstName, "last_name", user.LastName)
	query := `
	INSERT INTO users (first_name, last_name, age, is_married, password)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	var id int64
	err := s.pool.QueryRow(ctx, query, user.FirstName, user.LastName, user.Age, user.IsMarried, user.Password).Scan(&id)
	if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" && pgErr.ConstraintName == "idx_users_name" {
				return 0, fmt.Errorf("user with name %s %s already exists", user.FirstName, user.LastName)
			}
			s.logger.Error(err, "Failed to create user", "first_name", user.FirstName, "last_name", user.LastName)
			return 0, fmt.Errorf("failed to create user: %w", err)
	}

	s.logger.Info("User created", "id", id)
	return id, nil
}

func (s *PostgresStorage) GetUserByID(ctx context.Context, id int64) (domain.User, error) {
	s.logger.Info("Fetching user", "id", id)
	query := `
	SELECT id, first_name, last_name, age, is_married, password
	FROM users
	WHERE id = $1
	`
	var user domain.User
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.IsMarried,
		&user.Password,
	)
	if err == pgx.ErrNoRows {
		return domain.User{}, fmt.Errorf("user not found")
	}
	if err != nil {
		s.logger.Error(err, "Failed to get user", "id", id)
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	s.logger.Info("User fetched", "id", id, "full_name", user.FullName)
	return user, nil
}

func (s *PostgresStorage) CreateProduct(ctx context.Context, product domain.Product) (int64, error) {
	s.logger.Info("Creating product", "description", product.Description)
	query := `
	INSERT INTO products (description, tags, quantity, price)
	VALUES ($1, $2, $3, $4)
	RETURNIG id`

	var id int64
	err := s.pool.QueryRow(ctx, query, product.Description, product.Tags, product.Quantity, product.Price).Scan(&id)
	if err != nil {
		s.logger.Error(err, "Failed to create product", "description", product.Description)
		return 0, fmt.Errorf("failed to create product^ %w", err)
	}

	s.logger.Info("Product created", "id", id)
	return id, nil
}

func (s *PostgresStorage) GetProductByID(ctx context.Context, id int64) (domain.Product, error) {
	s.logger.Info("Fetching product", "id", id)
	query := `
	SELECT id, description, tags, quantity, price
	FROM products
	WHERE id = $1
	`

	var product domain.Product
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Description,
		&product.Tags,
		&product.Quantity,
		&product.Price,
	)
	if err == pgx.ErrNoRows {
		return domain.Product{}, fmt.Errorf("product not found")
	}
	if err != nil {
		s.logger.Error(err, "Failed to  get product", "id", id)
		return domain.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	s.logger.Info("Product fetched", "id", id, "description", product.Description)
	return product, nil
}

func (s *PostgresStorage) UpdateProductQuantity(ctx context.Context, id int64, quantity int) error {
	s.logger.Info("Update product quantity", "id", id, "quantity", quantity)
	query := `
	UPDATE products
	SET quantity = $1
	WHERE id = $2
	`

	result, err := s.pool.Exec(ctx, query, quantity, id)
	if err != nil {
		s.logger.Error(err, "Failed to update product quantity", "id", id)
		return fmt.Errorf("failed to update product quantity: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}

	s.logger.Info("Product quantity updated", "id", id)
	return nil
}

func (s *PostgresStorage) CreateOrder(ctx context.Context, userID int64, orderProducts []domain.OrderProduct, totalPrice float64) (int64, error) {
	s.logger.Info("Creating order", "user_id", userID, "total_price", totalPrice)
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error(err, "Failed to start transaction")
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var orderID int64
	query := `
	INSERT INTO orders (user_id, created_at, total_price)
	VALUES ($1, $2, $3)
	RETURNING id
`

	err = tx.QueryRow(ctx, query, userID, time.Now(), totalPrice).Scan(&orderID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return 0, fmt.Errorf("user with id %d not found", userID)
		}
		s.logger.Error(err, "Failed to create order", "user_id", userID)
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	for _, op := range orderProducts {
		var availableQty int
		var price float64
		query := `
	SELECT quantity, price
	FROM products
	WHERE id = $1
	FOR UPDATE
	`

		err := tx.QueryRow(ctx, query, op.ProductID).Scan(&availableQty, &price)
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("product with id %d not found", op.ProductID)
		}
		if err != nil {
			s.logger.Error(err, "Failed to check product", "product_id", op.ProductID)
			return 0, fmt.Errorf("failed to check product %d: %w", op.ProductID, err)
		}
		if availableQty < op.Quantity {
			return 0, fmt.Errorf("not enough quantity for product %d: available %d, requested %d", op.ProductID, availableQty, op.Quantity)
		}

		query = `
		INSERT INTO order_product (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
	`
		_, err = tx.Exec(ctx, query, orderID, op.ProductID, op.Quantity, price)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
				return 0, fmt.Errorf("product %d already exists in order %d", op.ProductID, orderID)
			}
			s.logger.Error(err, "Failed to add product to order", "product_id", op.ProductID)
			return 0, fmt.Errorf("failed to add product %d to order: %w", op.ProductID, err)
		}

		query = `
		UPDATE products
		SET quantity = quantity - $1
		WHERE id = $2
		`
		_, err = tx.Exec(ctx, query, op.Quantity, op.ProductID)
		if err != nil {
			s.logger.Error(err, "Failed to update product quantity", "product_id", op.ProductID)
			return 0, fmt.Errorf("failed to update quantity for product %d: %w", op.ProductID, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error(err, "Failed to commit transaction")
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info("Order created", "order_id", orderID)
	return orderID, nil
}

func (s *PostgresStorage) GetOrderByID(ctx context.Context, id int64) (domain.Order, error) {
	s.logger.Info("Fetching order", "id", id)
	query := `
	SELECT id, user_id, created_at, tota_price
	FROM orders
	WHERE id = $1
	`
	var order domain.Order
	err := s.pool.QueryRow(ctx, query, id).Scan(&order.ID, &order.UserID, &order.CreatedAt, &order.TotalPrice)
	if err == pgx.ErrNoRows {
		return domain.Order{}, fmt.Errorf("order not found")
	}
	if err != nil {
		s.logger.Error(err, "Failed to get order", "id", id)
		return domain.Order{}, fmt.Errorf("failed to get order: %w", err)
	}

	query = `
	SELECT order_id, product_id, quantity, price
	FROM order_products
	WHERE order_id = $1
	`

	rows, err := s.pool.Query(ctx, query, id)
	if err != nil {
		s.logger.Error(err, "Failed to get order products", "order_id", id)
		return domain.Order{}, fmt.Errorf("failed to get order products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var op domain.OrderProduct
		if err := rows.Scan(&op.OrderID, &op.ProductID, &op.Quantity, &op.Price); err != nil {
			s.logger.Error(err, "Failed to scan order product", "order_id", id)
			return domain.Order{}, fmt.Errorf("failed to scan order product: %w", err)
		}
		order.OrderProduct = append(order.OrderProduct, op)
	}
	if err := rows.Err(); err != nil {
		s.logger.Error(err, "Failed to iterate order products", "order_id", id)
		return domain.Order{}, fmt.Errorf("failed to iterate rows: %w", err)
	}
	s.logger.Info("Order fetched", "id", id)
	return order, nil
}
