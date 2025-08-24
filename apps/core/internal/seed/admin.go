package seed

import (
	"context"
	"fmt"

	"github.com/dopeCape/kova/internal/models"
)

// Creates the admin user and stores in the db
func (s *Seed) SeedAdmin() {
	fmt.Printf("Creating admin user: %s\n", s.SeedConfig.AdminEmail)
	_, err := s.userService.CreateUser(context.Background(), &models.CreateUserRequest{
		Email:    s.SeedConfig.AdminEmail,
		Username: s.SeedConfig.AdminUserName,
		Password: s.SeedConfig.AdminPassword,
	})

	if err != nil {
		fmt.Printf("Error creating admin user: %s\n", err)
		return
	}
}
