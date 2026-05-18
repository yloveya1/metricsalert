package handler

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/yloveya1/metricsalert/internal/service/controller"
)

//go:embed templates/*
var templatesFS embed.FS

type Handler struct {
	metricCtrl  controller.IMetricController
	metricsTmpl *template.Template
	key         string
}

func New(metric controller.IMetricController, key string) *Handler {
	tmplFS, _ := fs.Sub(templatesFS, "templates")

	return &Handler{
		metricCtrl: metric,
		metricsTmpl: template.Must(
			template.
				New("metrics.html").
				Funcs(templFunc).
				ParseFS(tmplFS, "metrics.html"),
		),
		key: key,
	}
}
