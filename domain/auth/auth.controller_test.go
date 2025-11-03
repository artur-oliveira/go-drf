package auth_test

import (
	"encoding/json"
	"grf/domain/auth"
	"grf/tests"
	"net/http"
	"testing"
)

func TestAuthEndpoints(t *testing.T) {
	clearAuthTables(testApp.DB)
	_, err := createTestFixtures(testApp.DB)
	if err != nil {
		t.Fatalf("Falha ao criar fixtures: %v", err)
	}

	t.Run("Login com Falha (Senha Errada)", func(t *testing.T) {
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost,
			URL:    "/v1/auth/token",
			Body:   auth.ObtainTokenDTO{Login: "admin", Password: "senhaerrada"},
		})
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Esperado 401, obteve %d", resp.StatusCode)
		}
	})

	t.Run("Get /me (Sem token)", func(t *testing.T) {
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet,
			URL:    "/v1/auth/me",
		})
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("Esperado 401, obteve %d", resp.StatusCode)
		}
	})

	t.Run("Login, Refresh, GetMe, ChangePassword (Fluxo Completo)", func(t *testing.T) {
		accessToken, refreshToken := loginAs(t, "user", "user123")

		respMe, bodyMe := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet,
			URL:    "/v1/auth/me",
			Token:  accessToken,
		})
		if respMe.StatusCode != http.StatusOK {
			t.Fatalf("/me falhou: %s", bodyMe)
		}
		var userResp auth.UserResponseDTO
		err := json.Unmarshal([]byte(bodyMe), &userResp)
		if err != nil {
			return
		}
		if userResp.Username != "user" {
			t.Errorf("Esperado 'user', obteve '%s'", userResp.Username)
		}

		respRefresh, bodyRefresh := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost,
			URL:    "/v1/auth/refresh",
			Body:   auth.RefreshTokenDTO{Refresh: refreshToken},
		})
		if respRefresh.StatusCode != http.StatusOK {
			t.Fatalf("Falha no Refresh: %s", bodyRefresh)
		}
		var newTokens auth.TokenResponseDTO
		err = json.Unmarshal([]byte(bodyRefresh), &newTokens)
		if err != nil {
			return
		}
		if newTokens.AccessToken == "" {
			t.Fatal("Refresh token não gerou um novo access token")
		}

		respChangeFail, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost,
			URL:    "/v1/auth/change-password",
			Token:  newTokens.AccessToken,
			Body:   auth.ChangePasswordDTO{OldPassword: "senhaerrada", NewPassword: "novasenha123"},
		})
		if respChangeFail.StatusCode != http.StatusBadRequest {
			t.Errorf("ChangePassword (senha errada): Esperado 400, obteve %d", respChangeFail.StatusCode)
		}

		respChangeOK, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost,
			URL:    "/v1/auth/change-password",
			Token:  newTokens.AccessToken,
			Body:   auth.ChangePasswordDTO{OldPassword: "user123", NewPassword: "novasenha123"},
		})
		if respChangeOK.StatusCode != http.StatusNoContent {
			t.Errorf("ChangePassword (sucesso): Esperado 204, obteve %d", respChangeOK.StatusCode)
		}

		respLoginFail, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost,
			URL:    "/v1/auth/token",
			Body:   auth.ObtainTokenDTO{Login: "user", Password: "user123"},
		})
		if respLoginFail.StatusCode != http.StatusUnauthorized {
			t.Errorf("Login (senha antiga): Esperado 401, obteve %d", respLoginFail.StatusCode)
		}

		loginAs(t, "user", "novasenha123")
	})
}
