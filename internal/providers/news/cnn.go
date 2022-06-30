package news

func GetCNNTopStoriesFeed() (*NewsFeed, error) {
	return ParseFeed("http://rss.cnn.com/rss/cnn_topstories.rss")
}

func GetCNNTechFeed() (*NewsFeed, error) {
	return ParseFeed("http://rss.cnn.com/rss/cnn_tech.rss")
}

func GetCNNWorldFeed() (*NewsFeed, error) {
	return ParseFeed("http://rss.cnn.com/rss/cnn_world.rss")
}
