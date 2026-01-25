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
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Невозможно установить соединение с keycloak")
	}
	token, ok := tokenObj.(*gocloak.JWT)
	if !ok {
		return nil, core.NewTechnicalError(nil, keycloakAuthService, "Невозможно получить информацию из токена авторизации")
	}

	// try to get specified client
	clients, err := client.GetClients(ctx, token.AccessToken, cfg.Realm, gocloak.GetClientsParams{
		ClientID: &cfg.ClientID, // Фильтруем по ClientID
	})
	if err != nil {
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Невозможно получить клиентов keycloak")
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
		return core.NewTechnicalError(err, keycloakAuthService, "Роль 'Гость' отсутствует")
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
		return core.NewTechnicalError(err, keycloakAuthService, "Ошибка при создании пользователя в keycloak")
	}

	err = kc.client.SetPassword(ctx, kc.token.AccessToken, userID, kc.cfg.Realm, password, false)
	if err != nil {
		return core.NewTechnicalError(err, keycloakAuthService, "Ошибка при установке пароля для пользователя в keycloak")
	}

	err = kc.client.AddClientRolesToUser(ctx, kc.token.AccessToken, kc.cfg.Realm, kc.clientID, userID, []gocloak.Role{*kcRole})
	if err != nil {
		return core.NewTechnicalError(err, keycloakAuthService, "Невозможно установить роль 'Гость' для пользователя")
	}

	return nil
}

func (kc *keycloakService) Login(ctx context.Context, username, password string) (JWT, error) {
	token, err := kc.client.Login(ctx, kc.cfg.ClientID, kc.cfg.ClientSecret, kc.cfg.Realm, username, password)
	if err != nil {
		return JWT{}, core.NewAccessError(err, keycloakAuthService, "Ошибка при авторизации в keycloak")
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
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Невозможно получить роли keycloak")
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
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Невозможно получить роли keycloak")
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
		return TokenUserInfo{}, core.NewTechnicalError(err, keycloakAuthService, "Ошибка при расшифровке токена авторизации")
	}

	var tokenUserInfo TokenUserInfo

	emailRaw, ok := (*claims)["email"]
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Не найден тэг : email")
	}
	email, ok := emailRaw.(string)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Тэг email некорректный!")
	}
	tokenUserInfo.Email = email

	usernameRaw, ok := (*claims)["preferred_username"]
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Не найден тэг : preferred_username")
	}
	username, ok := usernameRaw.(string)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Тэг preferred_username некорректный!")
	}
	tokenUserInfo.Username = username

	resource_access, ok := (*claims)["resource_access"]
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Не найден тэг : resource_access")
	}
	tokenClients, ok := resource_access.(map[string]any)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Тэг resource_access некорректный!")
	}
	clientResourcesRaw, ok := tokenClients[kc.cfg.ClientID]
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Не найдены клиенты по переданному идентификатору!")
	}
	clientResources, ok := clientResourcesRaw.(map[string]any)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Некорректный формат информации о клиентах!")
	}

	rolesRaw, ok := clientResources["roles"]
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Нет данных по ролям клиента!")
	}
	rolesSlice, ok := rolesRaw.([]interface{})
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Неверный формат токена. Некорректный формат ролей!")
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
		return nil, core.NewTechnicalError(err, keycloakAuthService, "Ошибка при поиске пользователей по заданным параметрам")
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
		return UserDTO{}, core.NewTechnicalError(err, keycloakAuthService, "Ошибка при получении пользователя!")
	}
	if len(kcUsers) == 0 {
		return UserDTO{}, core.NewLogicalError(err, keycloakAuthService, "Пользователя с такими данными не существует!")
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
		return core.NewTechnicalError(err, keycloakAuthService, "Ошибка при удалении пользователя!")
	}
	return nil
}

func GetTokenCtx(ctx context.Context) (string, error) {
	tokenCtxKey, ok := ctx.Value(TokenCtxKey).(string)
	if !ok {
		return "", core.NewLogicalError(nil, keycloakAuthService, "Токен не установлен в контекст!")
	}
	return tokenCtxKey, nil
}

func GetUserInfoCtx(ctx context.Context) (TokenUserInfo, error) {
	userInfoCtxKey, ok := ctx.Value(UserInfoCtxKey).(TokenUserInfo)
	if !ok {
		return TokenUserInfo{}, core.NewLogicalError(nil, keycloakAuthService, "Информация о пользователя отсутствует в токене!")
	}
	return userInfoCtxKey, nil
}
