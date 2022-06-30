package news

func GetTHNFeed() (*NewsFeed, error) {
	return ParseFeed("https://feeds.feedburner.com/TheHackersNews?format=xml")
}
