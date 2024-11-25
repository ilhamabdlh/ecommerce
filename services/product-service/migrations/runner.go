package migrations

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Migration struct {
	Version   int
	Name      string
	Timestamp time.Time
	Up        func(*mongo.Database) error
	Down      func(*mongo.Database) error
}

var migrations = []Migration{
	{
		Version:   1,
		Name:      "init",
		Timestamp: time.Now(),
		Up:        InitMigration,
		Down:      nil,
	},
}

func RunMigrations(db *mongo.Database) error {
	ctx := context.Background()

	// Create migrations collection if not exists
	migrationsCollection := db.Collection("migrations")

	for _, migration := range migrations {
		// Check if migration has been applied
		count, err := migrationsCollection.CountDocuments(ctx, map[string]interface{}{
			"version": migration.Version,
		})
		if err != nil {
			return err
		}

		if count == 0 {
			log.Printf("Running migration %d: %s", migration.Version, migration.Name)

			if err := migration.Up(db); err != nil {
				return err
			}

			// Record migration
			_, err = migrationsCollection.InsertOne(ctx, map[string]interface{}{
				"version":   migration.Version,
				"name":      migration.Name,
				"timestamp": time.Now(),
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
