package main

import (
	"context"

	"github.com/dopeCape/kova/internal/config"
	"github.com/dopeCape/kova/internal/seed"
	"github.com/dopeCape/kova/internal/store/postgres/repository"
)

func main() {
	appCfg := config.Load()
	db, store := repository.NewDefaultStore(context.Background(), appCfg)
	defer db.Close()
	seeder := seed.NewSeederWithDefaults(store)
	seeder.SeedAdmin()
}
