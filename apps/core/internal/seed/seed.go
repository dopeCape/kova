package seed

import (
	"github.com/dopeCape/kova/internal/services"
	"github.com/dopeCape/kova/internal/store"
	"github.com/dopeCape/kova/internal/utils"
)

type SeedConfig struct {
	AdminEmail    string
	AdminPassword string
	AdminUserName string
}

type Seed struct {
	SeedConfig
	userService *services.UserService
}

func NewSeederWithDefaults(store store.Store) *Seed {
	return &Seed{
		SeedConfig: SeedConfig{AdminEmail: DEFAULT_EMAIL,
			AdminPassword: DEFAULT_PASSWORD,
			AdminUserName: DEFAULT_USERNAME},
		userService: services.NewUserService(store),
	}
}

func NewSeederFromEnvs(store store.Store) *Seed {
	adminEmail := utils.GetEnv("ADMIN_EMAIL", DEFAULT_EMAIL)
	adminPassword := utils.GetEnv("ADMIN_PASSWORD", DEFAULT_PASSWORD)
	adminUserName := utils.GetEnv("ADMIN_USERNAME", DEFAULT_USERNAME)
	return &Seed{
		SeedConfig: SeedConfig{AdminEmail: adminEmail,
			AdminPassword: adminPassword,
			AdminUserName: adminUserName},
		userService: services.NewUserService(store),
	}
}
