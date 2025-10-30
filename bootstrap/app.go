package bootstrap

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/driver/sqlite" // (ou postgres)
	"gorm.io/gorm"
)

// App armazena dependências compartilhadas
type App struct {
	FiberApp  *fiber.App
	DB        *gorm.DB
	Validator *validator.Validate
}

// NewApp cria e inicializa a instância da aplicação.
func NewApp() *App {
	// 1. Conectar ao Banco de Dados
	// dsn := os.Getenv("DATABASE_URL")
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	db, err := gorm.Open(sqlite.Open("prod.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Falha ao conectar ao banco de dados: %v", err))
	}

	// (AutoMigrate deve ser movido para main.go ou um script de migração)
	// db.AutoMigrate(&models.User{})

	// 2. Inicializar o Validator
	validate := validator.New()

	// 3. Inicializar o Fiber
	app := fiber.New()
	app.Use(logger.New())

	return &App{
		FiberApp:  app,
		DB:        db,
		Validator: validate,
	}
}
