package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DanKRT-Star/task-manager/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")
	DB = db
}

func InitSchema() error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	return DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Task{},
		&model.ProjectMember{},
		&model.Epic{},
		&model.Milestone{},
	)
}

func ResetDatabase() error {
	if DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	if err := DB.Exec("DROP TABLE IF EXISTS milestones, epics, project_members, tasks, projects, users CASCADE").Error; err != nil {
		return err
	}

	return InitSchema()
}

func extractDatabaseName(dsn string) (string, error) {
	trimmed := strings.TrimSpace(dsn)
	if trimmed == "" {
		return "", fmt.Errorf("empty database dsn")
	}

	for _, field := range strings.Fields(trimmed) {
		if strings.HasPrefix(field, "dbname=") {
			return strings.TrimPrefix(field, "dbname="), nil
		}
	}

	return "", fmt.Errorf("dbname not found in dsn")
}
