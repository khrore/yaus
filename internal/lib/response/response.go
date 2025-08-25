package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusERROR = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusERROR,
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errMsgs []string
	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required filed", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid url", err.Field()))
		case "alias":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid alias", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid url", err.Field()))
		}
	}

	return Response{
		Status: StatusERROR,
		Error:  strings.Join(errMsgs, ", "),
	}
}
