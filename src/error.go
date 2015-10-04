package gopredit

import (
	"fmt"
	"html/template"
	"net/http"
)

type Go2MdError struct {
	Title   string
	Comment string
	Detail  string
	No      int
}

func (err Go2MdError) Error() string {
	return fmt.Sprintf("[%s]%s(%d)\n%s", err.Title, err.Comment, err.No, err.Detail)
}

func errorPage(w http.ResponseWriter, title, comment, detail string, no int) {
	err := Go2MdError{
		Title:   title,
		Comment: comment,
		Detail:  detail,
		No:      no,
	}

	tmpl, parseErr := template.ParseFiles("template/error.tmpl")
	if parseErr != nil {
		panic(parseErr)
	}

	exeError := tmpl.ExecuteTemplate(w, "error", err)
	if exeError != nil {
		panic(exeError)
	}
}
