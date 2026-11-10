// Command server runs the Kitab Dapur HTTP API.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	kitabhandler "kitab-dapur/api/internal/kitab/handler"
	kitabrepo "kitab-dapur/api/internal/kitab/repository"

	"kitab-dapur/api/internal/kitab"
	"kitab-dapur/api/internal/platform"
	"kitab-dapur/api/internal/recipe"
	recipehandler "kitab-dapur/api/internal/recipe/handler"
	reciperepo "kitab-dapur/api/internal/recipe/repository"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(log); err != nil {
		log.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := platform.LoadConfig(os.Getenv)
	if err != nil {
		return err
	}
	db, err := platform.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Each feature serves its own operations of the contract; together they make the whole API.
	api := struct {
		*platform.Health
		*kitabhandler.Kitabs
		*recipehandler.Recipes
	}{
		platform.NewHealth(db, log),
		kitabhandler.NewKitabs(kitab.NewService(kitabrepo.NewPostgres(db))),
		recipehandler.NewRecipes(recipe.NewService(reciperepo.NewPostgres(db))),
	}
	return platform.Serve(ctx, cfg.HTTPAddr, platform.NewAPIHandler(api, log), log)
}
