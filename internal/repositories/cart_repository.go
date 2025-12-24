package repositories

import (
	"errors"
	"fmt"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/internal/models"
	"github.com/Rzero6/self-checkout-api/internal/utils"
	"github.com/jackc/pgx/v5"
)

func CheckCartExist(sessionID string) (int64, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT id
		FROM carts
		WHERE session_id = $1 AND status = $2
	`

	var cartID int64
	err := config.DB.QueryRow(
		ctx,
		query,
		sessionID,
		utils.CartStatusActive,
	).Scan(&cartID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, nil
		}
		return -1, err
	}

	return cartID, nil
}

func CreateCart(cart models.Cart) (*models.Cart, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		INSERT INTO carts (session_id, status)
		VALUES ($1, $2)
		RETURNING id, session_id, status
	`
	var result models.Cart
	err := config.DB.QueryRow(ctx, query, cart.SessionID, cart.Status).Scan(&result.ID, &result.SessionID, &result.Status)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetActiveCartBySession(sessionID string) (*models.Cart, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT id, session_id, status
		FROM carts
		WHERE session_id = $1 AND status = $2
	`

	row := config.DB.QueryRow(ctx, query, sessionID, utils.CartStatusActive)

	var cart models.Cart
	err := row.Scan(&cart.ID, &cart.SessionID, &cart.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func PatchCartStatus(cartID int64, status string) (*models.Cart, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		UPDATE carts
		SET status = $2
		WHERE id = $1
		RETURNING id, session_id, status
	`
	var cart models.Cart
	err := config.DB.QueryRow(
		ctx, query, cartID, status,
	).Scan(&cart.ID, &cart.SessionID, &cart.Status)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func DeleteAllProductFromCart(cartID int64) error {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		DELETE FROM carts
		WHERE id = $1 AND status = $2
	`
	_, err := config.DB.Exec(ctx, query, cartID, utils.CartStatusActive)
	if err != nil {
		return fmt.Errorf("failed to delete cart: %w", err)
	}

	return nil
}
