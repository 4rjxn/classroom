package models

type Config struct {
	ClientId     string `toml:"client_id" json:"client_id"`
	ClientSecret string `toml:"client_secret" json:"client_secret"`
}
