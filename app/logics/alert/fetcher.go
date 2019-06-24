package alert

type Fetcher interface {
	Fetch(pair string) (*Response, error)
}
