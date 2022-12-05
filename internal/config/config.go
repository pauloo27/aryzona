package config

type BotConfig struct {
	Prefix         string `json:"prefix"`
	Token          string `json:"token"`
	OwnerID        string `json:"owner_id"`
	DonateMessage  string `json:"donate_message"`
	Env            string `json:"env"`
	Presence       string `json:"presence"`
	GitRepoURL     string `json:"git_repo_url"`
	YoutubeAPIKey  string `json:"youtube_api_key"`
	HTTPServerPort int    `json:"http_server_port"`
	HTTPServerExternalURL string `json:"http_server_external_url"`
}
