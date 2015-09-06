package go2md

import ()

type Slide struct {
	UserKey   string
	Title     string
	SubTitle  string
	SpeakDate string
	Tags      string
	Markdown  string
}

type Html struct {
	UserKey string
	Content string
}
