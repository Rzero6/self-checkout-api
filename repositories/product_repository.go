package repositories

import (
	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/models"
	"github.com/jackc/pgx/v5"
)

func GetProductsCount() (int, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	var total int
	err := config.DB.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM products`,
	).Scan(&total)

	return total, err
}

func GetAllProducts(page int, limit int) ([]models.Product, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	offset := (page - 1) * limit
	if page < 1 {
		page = 1
		offset = 0
	}

	rows, err := config.DB.Query(
		ctx,
		"SELECT id, barcode, name, price FROM products ORDER BY id LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Barcode, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByBarcode(barcode string) (*models.Product, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
	Select id, barcode, name, price
	FROM products
	WHERE barcode = $1
	LIMIT 1
	`

	row := config.DB.QueryRow(ctx, query, barcode)

	var item models.Product
	err := row.Scan(&item.ID, &item.Barcode, &item.Name, &item.Price)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &item, nil
}

func GetProductRandom() (*models.Product, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
	Select id, barcode, name, price
	FROM products
	ORDER BY RANDOM()
	LIMIT 1
	`
	row := config.DB.QueryRow(ctx, query)

	var item models.Product
	err := row.Scan(&item.ID, &item.Barcode, &item.Name, &item.Price)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func CreateProduct(item models.Product) (*models.Product, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		INSERT INTO products (barcode, name, price)
		VALUES ($1, $2, $3)
		RETURNING id, barcode, name, price
	`

	var created models.Product
	err := config.DB.QueryRow(ctx, query,
		item.Barcode, item.Name, item.Price).Scan(
		&created.ID,
		&created.Barcode,
		&created.Name,
		&created.Price,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func GetDonationsProduct() ([]models.Product, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	rows, err := config.DB.Query(
		ctx,
		`
		SELECT id, barcode, name, price 
		FROM products
		WHERE name ILIKE '%Donation%' 
		ORDER BY price
		`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Barcode, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}
