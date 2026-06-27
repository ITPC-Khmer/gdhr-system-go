package database

import (
	"fmt"
	"log"

	"backend/config"
	"backend/models"
	"backend/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global GORM database handle.
var DB *gorm.DB

// Connect opens the MySQL connection, runs migrations and seeds an admin user.
func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	logLevel := logger.Silent
	if cfg.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
	log.Println("✓ connected to MySQL database")

	migrate()
	seed()
}

func migrate() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Institute{},
		&models.Staff{},
		&models.Rank{},
		&models.Position{},
		&models.Holiday{},
		&models.LeaveType{},
		&models.LeaveRole{},
		&models.Leave{},
		&models.LeaveApproval{},
		&models.LeaveRequesterFile{},
		&models.LeaveYear{},
		&models.StaffInstituteRole{},
	); err != nil {
		log.Fatalf("auto-migration failed: %v", err)
	}
	log.Println("✓ database migrated")
}

// seed creates a default admin account on first run.
func seed() {
	var count int64
	DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	hashed, _ := utils.HashPassword("admin123")
	admin := models.User{
		Name:     "Administrator",
		Email:    "admin@example.com",
		Password: hashed,
		Role:     "admin",
		Active:   true,
	}
	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("seed admin failed: %v", err)
		return
	}
	log.Println("✓ seeded default admin -> admin@example.com / admin123")
}
