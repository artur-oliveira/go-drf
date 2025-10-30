package routes

import (
	"grf/bootstrap"
	"grf/domain/auth"
)

// RegisterRoutes registra todos os grupos de rotas da aplicação.
func RegisterRoutes(app *bootstrap.App) {
	// 1. Criar o grupo principal da API
	apiV1 := app.FiberApp.Group("/v1")

	// 2. Delegar o registro para módulos específicos
	// Passamos o grupo 'api' e as dependências compartilhadas.
	auth.RegisterUserRoutes(apiV1, app.DB, app.Validator)

	// 3. (Quando adicionar produtos)
	// RegisterProductRoutes(api, app.DB, app.Validator)

	// 4. (Quando adicionar pedidos)
	// RegisterOrderRoutes(api, app.DB, app.Validator)
}
