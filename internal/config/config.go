package config

type BotConfig struct {
	Prefix        string `yaml:"prefix"`
	Token         string `yaml:"token"`
	OwnerID       string `yaml:"owner_id"`
	DonateMessage string `yaml:"donate_message"`
	Env           string `yaml:"env"`
	Presence      string `yaml:"presence"`
	GitRepoURL    string `yaml:"git_repo_url"`
	YoutubeAPIKey string `yaml:"youtube_api_key"`
}
