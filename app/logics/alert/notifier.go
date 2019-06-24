package alert

type Notifier interface {
	Notify(pair string, direction Direction, differencePercentage string, price string, rawDataStr string) error
}
