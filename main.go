package main

import (
	"grf/bootstrap"
	"grf/routes"
	"log"
	// (Importar "github.com/joho/godotenv" para carregar .env em produção)
)

func main() {
	// 1. Carregar .env
	// godotenv.Load()

	// 2. Criar a instância da aplicação (obtém o contêiner de dependências)
	app := bootstrap.NewApp()

	// 3. Registrar todas as rotas da aplicação
	// Passa o contêiner 'app' para o pacote de rotas.
	routes.RegisterRoutes(app)

	// 4. Iniciar o servidor
	// port := os.Getenv("PORT")
	// if port == "" {
	//    port = "3000"
	// }
	// app.FiberApp.Listen(":" + port)

	log.Fatal(app.FiberApp.Listen(":3000"))
}
