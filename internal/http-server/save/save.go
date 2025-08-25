package save

import (
	"errors"
	"log/slog"
	"net/http"
	"yaus/internal/lib/logger/slogext"
	"yaus/internal/lib/response"
	"yaus/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias" validate:"required,alias"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v3@v3.5.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(writer http.ResponseWriter, httpRequest *http.Request) {
		const fnName = "url.save.New"

		log := log.With(
			slog.String("function", fnName),
			slog.String("requiest_id", middleware.GetReqID(httpRequest.Context())),
		)

		var request Request

		err := render.DecodeJSON(httpRequest.Body, &request)
		if err != nil {
			msg := "failed to decode request body"
			log.Error(msg, slogext.Err(err))

			render.JSON(writer, httpRequest, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", request))

		if err := validator.New().Struct(request); err != nil {
			validateError := err.(validator.ValidationErrors)

			log.Error("invalid request", slogext.Err(err))

			render.JSON(writer, httpRequest, response.ValidationError(validateError))

			return
		}
		id, err := urlSaver.SaveURL(request.URL, request.Alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("URL already exists", slog.String("url", request.URL))

			render.JSON(writer, httpRequest, response.Error("URL already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add URL", slogext.Err(err))

			render.JSON(writer, httpRequest, response.Error("failed to add URL"))

			return
		}

		log.Info("URL added", slog.Int64("id", id))

		render.JSON(writer, httpRequest, Response{
			Response: response.Ok(),
			Alias:    request.Alias,
		})
	}
}
