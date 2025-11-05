package controller_test

import (
	"grf/core/tests"
	"grf/domain/auth/dto"
	"net/http"
	"testing"
)

func TestPermissionCRUD(t *testing.T) {
	clearAuthTables(testApp.DB)
	_, err := createTestFixtures(testApp.DB)
	if err != nil {
		t.Fatalf("Falha ao criar fixtures: %v", err)
	}

	adminToken, _ := loginAs(t, "admin", "admin123")
	userToken, _ := loginAs(t, "user", "user123")

	t.Run("GET /permissions (Admin 200)", func(t *testing.T) {
		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet, URL: "/v1/permissions", Token: adminToken,
		})
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Esperado 200, obteve %d: %s", resp.StatusCode, body)
		}
	})

	t.Run("GET /permissions (User 403)", func(t *testing.T) {
		resp, _ := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodGet, URL: "/v1/permissions", Token: userToken,
		})
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Esperado 403, obteve %d", resp.StatusCode)
		}
	})

	t.Run("POST /permissions (Admin 201)", func(t *testing.T) {
		dto := dto.PermissionCreateDTO{
			Module: "test", Action: "create", Description: "Test perm",
		}
		resp, body := tests.MakeRequest(t, testApp.FiberApp, tests.RequestOptions{
			Method: http.MethodPost, URL: "/v1/permissions", Token: adminToken, Body: dto,
		})
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Esperado 201, obteve %d: %s", resp.StatusCode, body)
		}
	})
}
