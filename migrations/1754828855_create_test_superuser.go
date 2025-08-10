package migrations

import (
	"errors"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"

	"github.com/joho/godotenv"
)

func init() {
	m.Register(func(app core.App) error {
		godotenv.Load()

		password := os.Getenv("SUPERUSER_PASSWORD")
		if password == "" {
			return errors.New("SUPERUSER_PASSWORD is not set")
		}

		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		record := core.NewRecord(superusers)

		record.Set("email", "test@uley.kz")
		record.Set("password", password)

		return app.Save(record)
	}, func(app core.App) error {
		// add down queries...
		app.DB().Delete("users", dbx.NewExp("email = {:email}", dbx.Params{"email": "test@uley.kz"}))

		return nil
	})
}
