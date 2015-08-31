package go2md

import ()

type User struct {
	UserKey   string
	Name      string
	Job       string
	Email     string
	Url       string
	TwitterId string
}

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
