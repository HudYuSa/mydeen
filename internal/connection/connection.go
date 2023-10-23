package connection

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/HudYuSa/mydeen/internal/config"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(config *config.Config) {
	var err error

	// log gorm query
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second, // Slow SQL threshold
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,        // Disable color
	})

	DB, err = gorm.Open(postgres.Open(config.DSN), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}

	fmt.Println("? Connected Successfully to the Database")
}
