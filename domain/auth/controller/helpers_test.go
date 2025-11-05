package controller_test

import (
	"encoding/json"
	tests2 "grf/core/tests"
	"grf/domain/auth/dto"
	"grf/domain/auth/model"
	"net/http"
	"testing"

	"gorm.io/gorm"
)

var authTables = []string{
	"auth_user_permissions",
	"auth_user_groups",
	"auth_group_permissions",
	"auth_user",
	"auth_group",
	"auth_permission",
}

func clearAuthTables(db *gorm.DB) {
	tests2.ClearTables(db, authTables)
}

type TestFixtures struct {
	AdminUser    *model.User
	NormalUser   *model.User
	PermViewUser *model.Permission
	PermAddUser  *model.Permission
	PermViewGrp  *model.Permission
	PermAddGrp   *model.Permission
	PermViewPerm *model.Permission
	PermAddPerm  *model.Permission
}

func createTestFixtures(db *gorm.DB) (*TestFixtures, error) {
	perms := []*model.Permission{
		{Module: "auth", Action: "view_user"}, {Module: "auth", Action: "add_user"},
		{Module: "auth", Action: "change_user"}, {Module: "auth", Action: "delete_user"},
		{Module: "auth", Action: "view_group"}, {Module: "auth", Action: "add_group"},
		{Module: "auth", Action: "change_group"}, {Module: "auth", Action: "delete_group"},
		{Module: "auth", Action: "view_permission"}, {Module: "auth", Action: "add_permission"},
		{Module: "auth", Action: "change_permission"}, {Module: "auth", Action: "delete_permission"},
	}
	if err := db.Create(&perms).Error; err != nil {
		return nil, err
	}

	adminGroup := model.Group{Name: "Admin"}
	if err := db.Create(&adminGroup).Error; err != nil {
		return nil, err
	}
	err := db.Model(&adminGroup).Association("Permissions").Append(perms)
	if err != nil {
		return nil, err
	}

	adminUser := model.User{Username: "admin", Email: "admin@test.com", IsActive: true, IsSuperuser: true}
	err = adminUser.SetPassword("admin123")
	if err != nil {
		return nil, err
	}
	if err := db.Create(&adminUser).Error; err != nil {
		return nil, err
	}

	normalUser := model.User{Username: "user", Email: "user@test.com", IsActive: true}
	err = normalUser.SetPassword("user123")
	if err != nil {
		return nil, err
	}
	if err := db.Create(&normalUser).Error; err != nil {
		return nil, err
	}

	return &TestFixtures{
		AdminUser:    &adminUser,
		NormalUser:   &normalUser,
		PermViewUser: perms[0],
		PermAddUser:  perms[1],
		PermViewGrp:  perms[4],
		PermAddGrp:   perms[5],
		PermViewPerm: perms[8],
		PermAddPerm:  perms[9],
	}, nil
}

func loginAs(t *testing.T, username, password string) (accessToken, refreshToken string) {
	loginDTO := dto.ObtainTokenDTO{
		Login:    username,
		Password: password,
	}

	resp, body := tests2.MakeRequest(t, testApp.FiberApp, tests2.RequestOptions{
		Method: http.MethodPost,
		URL:    "/v1/auth/token",
		Body:   loginDTO,
	})

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Falha ao logar como %s, status %d: %s", username, resp.StatusCode, body)
	}

	var tokenResp dto.TokenResponseDTO
	err := json.Unmarshal([]byte(body), &tokenResp)
	if err != nil {
		return "", ""
	}
	if tokenResp.AccessToken == "" || tokenResp.RefreshToken == "" {
		t.Fatal("Falha ao logar, tokens vazios")
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken
}
