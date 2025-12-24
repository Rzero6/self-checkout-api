package repositories

import (
	"errors"
	"fmt"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/internal/models"
	"github.com/Rzero6/self-checkout-api/internal/utils"
	"github.com/jackc/pgx/v5"
)

func CheckExistingCartDetail(cartID, productID int64) (*models.CartDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT cd.id, p.id AS product_id, p.name AS product_name, p.price, quantity, p.price * quantity AS subtotal
		FROM cart_details cd
		JOIN products p on p.id = cd.product_id
		WHERE cd.cart_id = $1 AND p.id = $2
	`

	var cartDetail models.CartDetail
	err := config.DB.QueryRow(
		ctx, query, cartID, productID,
	).Scan(
		&cartDetail.ID, &cartDetail.ProductID, &cartDetail.ProductName, &cartDetail.Price, &cartDetail.Quantity, &cartDetail.Subtotal,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &cartDetail, nil
}

func CreateCartDetail(cartID int64, item models.Product, quantity int) (*models.CartDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		WITH inserted AS (
			INSERT INTO cart_details (cart_id, product_id, quantity)
			VALUES ($1, $2, $3)
			RETURNING id, cart_id, product_id, quantity
		)
		SELECT
			i.id,
			i.cart_id,
			p.id AS product_id,
			p.name AS product_name,
			p.price,
			i.quantity,
			p.price * i.quantity AS subtotal
		FROM inserted i
		JOIN products p ON p.id = i.product_id;
	`

	var created models.CartDetail
	err := config.DB.QueryRow(ctx, query, cartID, item.ID, quantity).Scan(
		&created.ID,
		&created.CartID,
		&created.ProductID,
		&created.ProductName,
		&created.Price,
		&created.Quantity,
		&created.Subtotal,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func GetCartDetailsBySessionID(sessionID string) ([]models.CartDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT cd.id, cd.cart_id, cd.product_id, p.name AS product_name, p.price, cd.quantity, p.price * cd.quantity AS subtotal
		FROM cart_details cd
		JOIN carts c ON c.id = cd.cart_id
		JOIN products p on cd.product_id = p.id
		WHERE c.session_id = $1 AND c.status = $2
	`
	rows, err := config.DB.Query(ctx, query, sessionID, utils.CartStatusActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	cartProducts := []models.CartDetail{}
	for rows.Next() {
		var ci models.CartDetail
		rows.Scan(&ci.ID, &ci.CartID, &ci.ProductID, &ci.ProductName, &ci.Price, &ci.Quantity, &ci.Subtotal)
		cartProducts = append(cartProducts, ci)
	}
	return cartProducts, nil
}

func DeleteProductFromCart(id int64) (*models.CartDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		WITH deleted AS (
			DELETE FROM cart_details
			WHERE id = $1
			RETURNING id, cart_id, product_id, quantity
		)
		SELECT
			d.id,
			p.name AS product_name,
			p.price,
			d.quantity,
			p.price * d.quantity AS subtotal
		FROM deleted d
		JOIN products p ON p.id = d.product_id;
	`
	var deleted models.CartDetail
	err := config.DB.QueryRow(ctx, query, id).
		Scan(&deleted.ID, &deleted.ProductName, &deleted.Price, &deleted.Quantity, &deleted.Subtotal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to delete product: %w", err)
	}

	return &deleted, nil
}

func UpdateProductInCart(id int64, quantity int) (*models.CartDetail, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		WITH updated AS (
			UPDATE cart_details
			SET quantity = $1
			WHERE id = $2
			RETURNING id, cart_id, product_id, quantity
		)
		SELECT
			u.id,
			u.cart_id,
			p.id AS product_id,
			p.name AS product_name,
			p.price,
			u.quantity,
			p.price * u.quantity AS subtotal
		FROM updated u
		JOIN products p ON p.id = u.product_id;
	`
	var updated models.CartDetail
	err := config.DB.QueryRow(ctx, query, quantity, id).
		Scan(&updated.ID, &updated.CartID, &updated.ProductID, &updated.ProductName, &updated.Price, &updated.Quantity, &updated.Subtotal)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return &updated, nil
}
