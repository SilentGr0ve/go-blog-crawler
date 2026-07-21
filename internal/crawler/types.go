package crawler

type Job struct {
	URL   string
	Depth int
}

type Page struct {
	URL   string
	Title string
	Text  string
	Links []string
}
