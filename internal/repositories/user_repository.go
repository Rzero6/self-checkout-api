package repositories

import (
	"github.com/Rzero6/self-checkout-api/config"
	"github.com/Rzero6/self-checkout-api/internal/models"
)

func GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := config.DBContext()
	defer cancel()

	query := `
		SELECT id, username, password, email, role
		FROM users
		WHERE username = $1
	`
	row := config.DB.QueryRow(ctx, query, username)
	var user models.User

	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
