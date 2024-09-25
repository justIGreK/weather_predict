package handler

import (
	"weather/internal/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	Weather
}

func NewHandler(weather *service.WeatherService) *Handler {
	return &Handler{Weather: weather}
}

func (h *Handler) InitRouters() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/weather", h.GetWeather)
	return r
}
