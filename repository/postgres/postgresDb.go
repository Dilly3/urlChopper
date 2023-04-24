package postgres

import (
	"fmt"
	"os"

	"github.com/dilly3/urlshortner/internal"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	DB *gorm.DB
}

func NewPostgreRepository(user string, dbName string, password string, port string) (*PostgresRepository, error) {

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", user, password, dbName, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&internal.Redirect{})
	if err != nil {
		fmt.Printf("failed to migrate: %v", err)
		os.Exit(1)
	}
	return &PostgresRepository{DB: db}, nil
}

func (p *PostgresRepository) Find(code string) (*internal.Redirect, error) {
	redirect := &internal.Redirect{}

	if err := p.DB.Model(&internal.Redirect{}).Where("code = ?", code).First(redirect).Error; err != nil {
		return nil, errors.Wrap(err, "repository.postgresDb.Find")
	}

	return redirect, nil
}

func (p *PostgresRepository) Store(redirect *internal.Redirect) error {
	if err := p.DB.Create(redirect).Error; err != nil {
		return errors.Wrap(err, "repository.postgresDb.Store")
	}
	return nil
}
