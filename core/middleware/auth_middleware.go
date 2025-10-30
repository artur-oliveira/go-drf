package middleware

import "github.com/gofiber/fiber/v2"

func MiddlewareAuthentication(c *fiber.Ctx) error {
	// Lógica de validação do Token (ex: 'Authorization: Bearer <token>')
	// ...

	// Se válido:
	// user, err := GetUserFromToken(tokenString)
	// if err != nil {
	//    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token inválido"})
	// }

	// Anexa o usuário para uso posterior (em permissões ou handlers)
	// c.Locals("user", user)
	return c.Next()
}
