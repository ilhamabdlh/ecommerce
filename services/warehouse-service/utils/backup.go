package utils

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type BackupService struct {
	db         *mongo.Database
	backupPath string
}

func NewBackupService(db *mongo.Database, backupPath string) *BackupService {
	return &BackupService{
		db:         db,
		backupPath: backupPath,
	}
}

func (b *BackupService) PerformBackup(ctx context.Context) error {
	timestamp := time.Now().Format("20060102150405")
	backupFile := fmt.Sprintf("%s/backup_%s", b.backupPath, timestamp)

	// Create backup command
	cmd := mongo.Pipeline{
		{{"$out", backupFile}},
	}

	// Execute backup for each collection
	collections, err := b.db.ListCollectionNames(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	for _, collection := range collections {
		_, err := b.db.Collection(collection).Aggregate(ctx, cmd)
		if err != nil {
			return fmt.Errorf("failed to backup collection %s: %v", collection, err)
		}
	}

	Logger.Infof("Backup completed successfully to %s", backupFile)
	return nil
}

func (b *BackupService) RestoreFromBackup(ctx context.Context, backupFile string) error {
	// Restore command
	cmd := mongo.Pipeline{
		{{"$merge", backupFile}},
	}

	collections, err := b.db.ListCollectionNames(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to list collections: %v", err)
	}

	for _, collection := range collections {
		_, err := b.db.Collection(collection).Aggregate(ctx, cmd)
		if err != nil {
			return fmt.Errorf("failed to restore collection %s: %v", collection, err)
		}
	}

	Logger.Infof("Restore completed successfully from %s", backupFile)
	return nil
}
