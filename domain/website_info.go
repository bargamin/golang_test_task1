package domain

type WebsiteInfo struct {
	Url               string
	HTMLVersion       string
	Title             string
	HeadingsCounts    map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	IsExistLoginForm  bool
}
