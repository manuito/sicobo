package store

import (
	"github.com/manuito/sicobo/application"

	"go.mongodb.org/mongo-driver/mongo"
)

type BookDatabaseSpec struct {
	Name         string
	TotalBooks   int64
	CreateTime   string
	LastLoadTime string
}

type activeBookDatabase struct {
	BookDatabaseSpec BookDatabaseSpec
	Database         *mongo.Database
}

var active activeBookDatabase

func switchToDatabase(spec BookDatabaseSpec, database *mongo.Database) {
	application.Info("Switching to database", spec.Name)
	active = activeBookDatabase{spec, database}
}

// ActiveSpec shortcut to active spec
func ActiveSpec() *BookDatabaseSpec {
	return &active.BookDatabaseSpec
}
