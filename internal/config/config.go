package config

type BotConfig struct {
	Prefix                string `json:"prefix" env:"PREFIX"`
	Token                 string `json:"token" env:"TOKEN,secret"`
	OwnerID               string `json:"owner_id" env:"OWNER_ID"`
	DonateMessage         string `json:"donate_message" env:"DONATE_MESSAGE"`
	Env                   string `json:"env" env:"ENV"`
	Presence              string `json:"presence" env:"PRESENCE"`
	GitRepoURL            string `json:"git_repo_url" env:"GIT_REPO_URL"`
	YoutubeAPIKey         string `json:"youtube_api_key" env:"YOUTUBE_API_KEY,secret"`
	HTTPServerPort        int    `json:"http_server_port" env:"HTTP_SERVER_PORT"`
	HTTPServerExternalURL string `json:"http_server_external_url" env:"HTTP_SERVER_EXTERNAL_URL"`

	DB `json:"db"`
}

type DB struct {
	Host     string `json:"host" env:"DB_HOST"`
	Port     int    `json:"port" env:"DB_PORT"`
	User     string `json:"user" env:"DB_USER"`
	Password string `json:"password" env:"DB_PASSWORD,secret"`
	Database string `json:"database" env:"DB_DATABASE"`
	SSL      bool   `json:"ssl" env:"DB_SSL"`
}
