package API

type gamePlatform struct {
	PlatformName string
}

type AddGame struct {
	GenreName string
	GameName  string
}

type UpdateReleaseYear struct {
	Year          int
	GameName      string
	PublisherName string
	PlatformName  string
}
