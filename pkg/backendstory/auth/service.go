package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/ActuallyHello/backendstory/pkg/config"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/Nerzal/gocloak/v13"
)

type TokenCtx string
type UserInfoCtx string

const (
	keycloakAuthService             = "KEYCLOAK_SERVICE_CODE"
	guestRole                       = "guest"
	TokenCtxKey         TokenCtx    = "token"
	UserInfoCtxKey      UserInfoCtx = "userInfo"
)

type AuthService interface {
	RegisterUser(ctx context.Context, username, email, password string) error
	DeleteUser(ctx context.Context, email string) error
	Login(ctx context.Context, username, password string) (JWT, error)

	RefreshToken(ctx context.Context, refreshToken string) (JWT, error)

	GetUserByEmail(ctx context.Context, email string) (UserDTO, error)
	GetUsers(ctx context.Context) ([]UserDTO, error)
	GetRoles(ctx context.Context) ([]string, error)
	GetRolesByUser(ctx context.Context, username string) ([]string, error)
	GetTokenUserInfo(ctx context.Context, token string) (TokenUserInfo, error)
}

type keycloakService struct {
	client   *gocloak.GoCloak
	clientID string
	token    *gocloak.JWT
	cfg      *config.KeycloakConfig
}

func NewKeycloakService(ctx context.Context, cfg *config.KeycloakConfig) (*keycloakService, error) {
	client := gocloak.NewClient(cfg.Host)

	// try to get client token
	tokenObj, err := core.Retry(
		"Login keycloak client",
		func() (any, error) {
			return client.LoginClient(ctx, cfg.ClientID, cfg.ClientSecret, cfg.Realm)
		},
		core.SetMaxRetriesOpt(3),
		core.SetMaxDelayOpt(5*time.Second),
	)
	if err != nil {
		slog.Error("Coundn't establish connection with keycloak", "error", err)
		return nil, core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	token, ok := tokenObj.(*gocloak.JWT)
	if !ok {
		slog.Error("Couldn't parse token to JWT format", "error", err)
		return nil, core.NewTechnicalError(nil, keycloakAuthService, "Unexpected token type!")
	}

	// try to get specified client
	clients, err := client.GetClients(ctx, token.AccessToken, cfg.Realm, gocloak.GetClientsParams{
		ClientID: &cfg.ClientID, // Фильтруем по ClientID
	})
	if err != nil {
		slog.Error("Failed to get clients", "error", err)
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Failed to get clients list")
	}
	var clientUUID string
	if len(clients) == 0 {
		slog.Warn("No clients found with ClientID", "clientID", cfg.ClientID)
	} else {
		clientUUID = *clients[0].ID
	}

	return &keycloakService{
		client:   client,
		clientID: clientUUID,
		token:    token,
		cfg:      cfg,
	}, nil
}

func (kc *keycloakService) RegisterUser(ctx context.Context, username, email, password string) error {
	kcRole, err := kc.client.GetClientRole(ctx, kc.token.AccessToken, kc.cfg.Realm, kc.clientID, guestRole)
	if err != nil {
		slog.Error("Cannot find role!", "role", guestRole)
		return core.NewTechnicalError(err, keycloakAuthService, "No such role as 'guest'!")
	}

	var (
		user = gocloak.User{
			Username: gocloak.StringP(username),
			Email:    gocloak.StringP(email),
			Enabled:  gocloak.BoolP(true),
		}
	)

	userID, err := kc.client.CreateUser(ctx, kc.token.AccessToken, kc.cfg.Realm, user)
	if err != nil {
		slog.Error("Error while creating user", "user", user, "error", err)
		return core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	err = kc.client.SetPassword(ctx, kc.token.AccessToken, userID, kc.cfg.Realm, password, false)
	if err != nil {
		slog.Error("Error setting password to user", "user", user, "error", err)
		return core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	err = kc.client.AddClientRolesToUser(ctx, kc.token.AccessToken, kc.cfg.Realm, kc.clientID, userID, []gocloak.Role{*kcRole})
	if err != nil {
		slog.Error("Cannot set role!", "role", kcRole.Name, "userID", userID)
		return core.NewTechnicalError(err, keycloakAuthService, "Cannot set role guest for user")
	}

	slog.Info("Register user", "user", user)
	return nil
}

func (kc *keycloakService) Login(ctx context.Context, username, password string) (JWT, error) {
	token, err := kc.client.Login(ctx, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm, username, password)
	if err != nil {
		slog.Error("Cannot login user", "username", username, "error", err)
		return JWT{}, core.NewAccessError(err, keycloakAuthService, err.Error())
	}
	return JWT{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
	}, nil
}

func (kc *keycloakService) RefreshToken(ctx context.Context, refreshToken string) (JWT, error) {
	token, err := kc.client.RefreshToken(ctx, refreshToken, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm)
	if err != nil {
		slog.Error("Cannot refresh token", "error", err)
		return JWT{}, core.NewAccessError(err, keycloakAuthService, err.Error())
	}
	return JWT{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        token.ExpiresIn,
		RefreshExpiresIn: token.RefreshExpiresIn,
	}, nil
}

func (kc *keycloakService) GetRoles(ctx context.Context) ([]string, error) {
	params := gocloak.GetRoleParams{}
	kcRoles, err := kc.client.GetClientRoles(ctx, kc.token.AccessToken, kc.cfg.Realm, kc.clientID, params)
	if err != nil {
		slog.Error("Cannot get roles", "error", err)
		return nil, core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	roles := make([]string, 0, len(kcRoles))
	for _, kcRole := range kcRoles {
		roles = append(roles, *kcRole.Name)
	}
	return roles, nil
}

func (kc *keycloakService) GetRolesByUser(ctx context.Context, username string) ([]string, error) {
	userDTO, err := kc.GetUserByEmail(ctx, username)
	if err != nil {
		return nil, err
	}

	kcRoles, err := kc.client.GetClientRolesByUserID(ctx, kc.token.AccessToken, kc.cfg.Realm, kc.clientID, userDTO.ID)
	if err != nil {
		slog.Error("Couldn't find user roles", "username", username, "error", err)
		return nil, core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	roles := make([]string, 0, len(kcRoles))
	for _, kcRole := range kcRoles {
		roles = append(roles, *kcRole.Name)
	}

	return roles, nil
}

func (kc *keycloakService) GetTokenUserInfo(ctx context.Context, token string) (TokenUserInfo, error) {
	_, claims, err := kc.client.DecodeAccessToken(ctx, token, kc.cfg.Realm)
	if err != nil {
		slog.Error("Couldn't decode access token", "error", err)
		return TokenUserInfo{}, core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}

	var tokenUserInfo TokenUserInfo

	emailRaw, ok := (*claims)["email"]
	if !ok {
		slog.Error("Couldn't decode email tag", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any resource_access tags for token")
	}
	email, ok := emailRaw.(string)
	if !ok {
		slog.Error("Couldn't decode email struct", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Expected client resources!")
	}
	tokenUserInfo.Email = email

	usernameRaw, ok := (*claims)["preferred_username"]
	if !ok {
		slog.Error("Couldn't decode username tag", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any resource_access tags for token")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		slog.Error("Couldn't decode username struct", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Expected client resources!")
	}
	tokenUserInfo.Username = username

	resource_access, ok := (*claims)["resource_access"]
	if !ok {
		slog.Error("Couldn't decode recource access tag", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any resource_access tags for token")
	}
	tokenClients, ok := resource_access.(map[string]any)
	if !ok {
		slog.Error("Couldn't decode recource access struct", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Expected client resources!")
	}
	clientResourcesRaw, ok := tokenClients[kc.cfg.ClientID]
	if !ok {
		slog.Error("Couldn't get token clients by client id tag", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no client resources!")
	}
	clientResources, ok := clientResourcesRaw.(map[string]any)
	if !ok {
		slog.Error("Couldn't get token clients by client id struct", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no client resources as expected structure!")
	}

	rolesRaw, ok := clientResources["roles"]
	if !ok {
		slog.Error("Couldn't get roles tag", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. There is no any roles for client in token")
	}
	rolesSlice, ok := rolesRaw.([]interface{})
	if !ok {
		slog.Error("Couldn't get roles struct", "error", err)
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Invalid token structure. Roles structed not in array")
	}
	roles := make([]string, len(rolesSlice))
	for i := 0; i < len(rolesSlice); i++ {
		roles[i] = rolesSlice[i].(string)
	}
	tokenUserInfo.Roles = roles

	return tokenUserInfo, nil
}

func (kc *keycloakService) GetUsers(ctx context.Context) ([]UserDTO, error) {
	params := gocloak.GetUsersParams{}
	kcUsers, err := kc.client.GetUsers(ctx, kc.token.AccessToken, kc.cfg.Realm, params)
	if err != nil {
		slog.Error("Couldn't get users", "error", err)
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Couldn't find any users in current realm!")
	}

	users := make([]UserDTO, 0, len(kcUsers))
	for _, kcUser := range kcUsers {
		users = append(users, UserDTO{
			ID:        *kcUser.ID,
			Email:     *kcUser.Email,
			Username:  *kcUser.Username,
			CreatedAt: time.Unix(0, *kcUser.CreatedTimestamp),
		})
	}

	return users, nil
}

func (kc *keycloakService) GetUserByEmail(ctx context.Context, email string) (UserDTO, error) {
	params := gocloak.GetUsersParams{
		Email: &email,
	}
	// always return 1 element
	kcUsers, err := kc.client.GetUsers(ctx, kc.token.AccessToken, kc.cfg.Realm, params)
	if err != nil {
		slog.Error("Couldn't find user by", "email", email, "error", err)
		return UserDTO{}, core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	if len(kcUsers) == 0 {
		slog.Error("User doesn't exist!", "email", email)
		return UserDTO{}, core.NewLogicalError(err, keycloakAuthService, "Couldn't find user by this username")
	}

	userDTO := UserDTO{
		ID:        *kcUsers[0].ID,
		Email:     *kcUsers[0].Email,
		Username:  *kcUsers[0].Username,
		CreatedAt: time.Unix(0, *kcUsers[0].CreatedTimestamp),
	}

	return userDTO, nil
}

func (kc *keycloakService) DeleteUser(ctx context.Context, email string) error {
	userDTO, err := kc.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	err = kc.client.DeleteUser(ctx, kc.token.AccessToken, kc.cfg.Realm, userDTO.ID)
	if err != nil {
		slog.Error("Error while deleteting user", "email", email, "error", err)
		return core.NewTechnicalError(err, keycloakAuthService, err.Error())
	}
	slog.Info("User deleted", "email", email)
	return nil
}

func GetTokenCtx(ctx context.Context) (string, error) {
	tokenCtxKey, ok := ctx.Value(TokenCtxKey).(string)
	if !ok {
		return "", core.NewLogicalError(nil, keycloakAuthService, "Couldn't find token in context")
	}
	return tokenCtxKey, nil
}

func GetUserInfoCtx(ctx context.Context) (TokenUserInfo, error) {
	userInfoCtxKey, ok := ctx.Value(UserInfoCtxKey).(TokenUserInfo)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Couldn't find userInfo in context")
	}
	return userInfoCtxKey, nil
}
