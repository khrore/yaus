package responce

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Responce struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "Ok"
	StatusERROR = "Error"
)

func Ok() Responce {
	return Responce{
		Status: StatusOK,
	}
}

func Error(msg string) Responce {
	return Responce{
		Status: StatusERROR,
		Error:  msg,
	}
}

func ValidationError(errors validator.ValidationErrors) Responce {
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

	return Responce{
		Status: StatusERROR,
		Error:  strings.Join(errMsgs, ", "),
	}
}
