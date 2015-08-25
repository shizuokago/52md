package go2md

import ()

type User struct {
	Name      string
	Job       string
	Email     string
	Url       string
	TwitterId string
}

type Slide struct {
	Title     string
	SubTitle  string
	SpeakDate string
	Tags      string
}

type Markdown struct {
	Content string
}

type Html struct {
	Content string
}
