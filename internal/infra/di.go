package infra

import (
	"clean_architecture_go/internal/domain"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDatabase() error {
	var initError error
	once.Do(func() {
		initError = initializeDB()
	})
	return initError
}

func initializeDB() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		getEnvOrDefault("DB_HOST", "87.248.131.253"),
		getEnvOrDefault("DB_USER", "shahiapp_user"),
		getEnvOrDefault("DB_PASSWORD", "shahiapp_password"),
		getEnvOrDefault("DB_NAME", "clean_architecture_go"),
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_SSLMODE", "disable"),
		getEnvOrDefault("DB_TIMEZONE", "UTC"),
	)

	var err error
	dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	log.Println("✅ Database initialized successfully")
	return nil
}

func AddInfra() {
	if err := InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	domain.AutoMigrate(dbInstance)

}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}