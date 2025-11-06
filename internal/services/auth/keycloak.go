package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/ActuallyHello/backendstory/internal/config"
	"github.com/ActuallyHello/backendstory/internal/core"
	appError "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/Nerzal/gocloak/v13"
)

const (
	keycloakAuthService = "KEYCLOAK_SERVICE_CODE"
	guestRole           = "guest"
)

type keycloakService struct {
	client *gocloak.GoCloak
	token  *gocloak.JWT
	cfg    *config.KeycloakConfig
}

func NewKeycloakService(cfg *config.KeycloakConfig) (*keycloakService, error) {
	client := gocloak.NewClient(cfg.Host)

	tokenObj, err := core.Retry(
		"Login keycloak client",
		func() (any, error) {
			return client.LoginClient(context.Background(), cfg.ClientID, cfg.ClientSecret, cfg.Realm)
		},
		core.SetMaxRetriesOpt(3),
		core.SetMaxDelayOpt(5*time.Second),
	)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	token, ok := tokenObj.(*gocloak.JWT)
	if !ok {
		return nil, appError.NewTechnicalError(nil, keycloakAuthService, "Unexpected token type!")
	}

	return &keycloakService{
		client: client,
		token:  token,
		cfg:    cfg,
	}, nil
}

func (kc *keycloakService) RegisterUser(ctx context.Context, username, email, password string) (dto.JWT, error) {
	user := gocloak.User{
		Username: gocloak.StringP(username),
		Email:    gocloak.StringP(email),
		Enabled:  gocloak.BoolP(true),
	}

	userID, err := kc.client.CreateUser(ctx, kc.token.AccessToken, kc.cfg.Realm, user)
	if err != nil {
		return dto.JWT{}, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	err = kc.client.SetPassword(ctx, kc.token.AccessToken, userID, kc.cfg.Realm, password, false)
	if err != nil {
		return dto.JWT{}, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	kcRole, err := kc.client.GetRealmRole(ctx, kc.token.AccessToken, kc.cfg.Realm, guestRole)
	if err != nil {
		slog.Error("Cannot find role!", "role", guestRole, "userID", userID)
		return dto.JWT{}, appError.NewTechnicalError(err, keycloakAuthService, "No such role as 'guest'!")
	}

	err = kc.client.AddRealmRoleToUser(ctx, kc.token.AccessToken, kc.cfg.Realm, userID, []gocloak.Role{*kcRole})
	if err != nil {
		slog.Error("Cannot set role!", "role", kcRole.Name, "userID", userID)
		return dto.JWT{}, appError.NewTechnicalError(err, keycloakAuthService, "Cannot set role guest for user")
	}

	// GET ACCESS TOKEN FOR USER AND CHECK REGISTRATION PROCESS
	jwt, err := kc.Login(ctx, username, password)
	if err != nil {
		return dto.JWT{}, err
	}
	return jwt, nil
}

func (kc *keycloakService) Login(ctx context.Context, username, password string) (dto.JWT, error) {
	token, err := kc.client.Login(ctx, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm, username, password)
	if err != nil {
		return dto.JWT{}, appError.NewAccessError(err, keycloakAuthService, err.Error())
	}
	return dto.JWT{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
	}, nil
}

func (kc *keycloakService) RefreshToken(ctx context.Context, refreshToken string) (dto.JWT, error) {
	token, err := kc.client.RefreshToken(ctx, refreshToken, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm)
	if err != nil {
		return dto.JWT{}, appError.NewAccessError(err, keycloakAuthService, err.Error())
	}
	return dto.JWT{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
	}, nil
}

func (kc *keycloakService) GetRoles(ctx context.Context) ([]string, error) {
	params := gocloak.GetRoleParams{}
	kcRoles, err := kc.client.GetRealmRoles(ctx, kc.token.AccessToken, kc.cfg.Realm, params)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, "No roles for current realm!")
	}
	roles := make([]string, 0, len(kcRoles))
	for _, kcRole := range kcRoles {
		roles = append(roles, *kcRole.Name)
	}
	return roles, nil
}

func (kc *keycloakService) GetRolesByUser(ctx context.Context, username string) ([]string, error) {
	params := gocloak.GetUsersParams{
		Username: &username,
	}
	// always return 1 element
	kcUsers, err := kc.client.GetUsers(ctx, kc.token.AccessToken, kc.cfg.Realm, params)
	if err != nil {
		slog.Error("Couldn't find user by", "username", username, "err", err)
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	if len(kcUsers) == 0 {
		slog.Error("User doesn't exist!", "username", username)
		return nil, appError.NewLogicalError(err, keycloakAuthService, "Couldn't find user by this username")
	}

	kcUser := kcUsers[0]
	kcRoles, err := kc.client.GetRealmRolesByUserID(ctx, kc.token.AccessToken, kc.cfg.Realm, *kcUser.ID)
	if err != nil {
		slog.Error("Couldn't find user roles", "username", username, "err", err)
		return nil, appError.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	roles := make([]string, 0, len(kcRoles))
	for _, kcRole := range kcRoles {
		roles = append(roles, *kcRole.Name)
	}

	return roles, nil
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

func (kc *keycloakService) GetUsers(ctx context.Context) ([]dto.UserDTO, error) {
	params := gocloak.GetUsersParams{}
	kcUsers, err := kc.client.GetUsers(ctx, kc.token.AccessToken, kc.cfg.Realm, params)
	if err != nil {
		return nil, appError.NewTechnicalError(err, keycloakAuthService, "Couldn't find any users in current realm!")
	}

	users := make([]dto.UserDTO, 0, len(kcUsers))
	for _, kcUser := range kcUsers {
		users = append(users, dto.UserDTO{
			Email:     *kcUser.Email,
			Username:  *kcUser.Username,
			CreatedAt: time.Unix(0, *kcUser.CreatedTimestamp),
		})
	}

	return users, nil
}
