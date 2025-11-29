package config

type KeycloakConfig struct {
	Host         string `mapstructure:"host"`
	Realm        string `mapstructure:"realm"`
	ClientID     string `mapstructure:"client-id"`
	ClientSecret string `mapstructure:"client-secret"`
	RedirectURI  string `mapstructure:"redirect-url"`
}
