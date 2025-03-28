package crawler

import "time"

type News struct {
	Title       string
	Description string
	Cover       string
	Link        string
	Time        time.Time
	Content     string
}
