package auth

import (
	"context"
	"time"

	"github.com/ActuallyHello/backendstory/internal/config"
	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/Nerzal/gocloak/v13"
)

const (
	keycloakAuthService = "KEYCLOAK_SERVICE_CODE"
)

type keycloakService struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	cfg    *config.KeycloakConfig
}

func NewKeycloakService(cfg *config.KeycloakConfig) (*keycloakService, error) {
	client := gocloak.NewClient(cfg.Host)

	token, err := client.LoginClient(context.Background(), cfg.ClientID, cfg.ClientSecret, cfg.Realm)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	return &keycloakService{
		client: client,
		token:  token,
		cfg:    cfg,
	}, nil
}

func (kc *keycloakService) RegisterUser(ctx context.Context, username, email, password string) error {
	user := gocloak.User{
		Username: gocloak.StringP(username),
		Email:    gocloak.StringP(email),
		Enabled:  gocloak.BoolP(true),
	}

	userID, err := kc.client.CreateUser(ctx, kc.token.AccessToken, kc.cfg.Realm, user)
	if err != nil {
		return appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	err = kc.client.SetPassword(ctx, kc.token.AccessToken, userID, kc.cfg.Realm, password, false)
	if err != nil {
		return appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	// TODO set roles
	return nil
}

func (kc *keycloakService) Login(ctx context.Context, username, password string) (map[string]any, error) {
	token, err := kc.client.Login(ctx, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm, username, password)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	return map[string]any{
		"token":              token.AccessToken,
		"refresh_token":      token.RefreshToken,
		"expires_in":         token.ExpiresIn,
		"refresh_expires_in": token.RefreshExpiresIn,
	}, nil
}

func (kc *keycloakService) GetRolesFromToken(ctx context.Context, token string) ([]string, error) {
	_, claims, err := kc.client.DecodeAccessToken(ctx, token, kc.cfg.Realm)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	resource_access, ok := (*claims)["resource_access"]
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any resource_access tags for token")
	}
	tokenClients, ok := resource_access.(map[string]any)
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Expected client resources!")
	}
	clientResourcesRaw, ok := tokenClients[kc.cfg.ClientID]
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no client resources!")
	}
	clientResources, ok := clientResourcesRaw.(map[string]any)
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no client resources as expected structure!")
	}

	rolesRaw, ok := clientResources["roles"]
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any roles for client in token")
	}
	rolesSlice, ok := rolesRaw.([]interface{})
	if !ok {
		return nil, appError.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Roles structed not in array")
	}
	roles := make([]string, len(rolesSlice))
	for i := 0; i < len(rolesSlice); i++ {
		roles[i] = rolesSlice[i].(string)
	}

	return roles, nil
}

func (kc *keycloakService) GetUsers(ctx context.Context, token string) ([]dto.UserDTO, error) {
	kcUsers, err := kc.client.GetUsers(ctx, token, kc.cfg.Realm, gocloak.GetUsersParams{})
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, "Couldn't find any users in current realm!")
	}

	users := make([]dto.UserDTO, 0, len(kcUsers))
	for _, kcUser := range kcUsers {
		users = append(users, dto.UserDTO{
			Email:     *kcUser.Email,
			Username:  *kcUser.Username,
			CreatedAt: time.Unix(0, *kcUser.CreatedTimestamp),
			Roles:     *kcUser.RealmRoles,
		})
	}

	return users, nil
}
