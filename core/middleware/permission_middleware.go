package middleware

import "github.com/gofiber/fiber/v2"

func MiddlewarePermission(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// user, ok := c.Locals("user").(*User)
		// if !ok {
		//    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Autenticação requerida"})
		// }

		// if role == "admin" && !user.IsAdmin {
		//    return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permissão negada (Admin requerido)"})
		// }

		// if role == "self" {
		//    // Lógica para verificar se o usuário é o dono do recurso (ex: c.Params("id"))
		// }

		return c.Next()
	}
}
