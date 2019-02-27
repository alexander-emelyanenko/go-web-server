package controllers

import (
	"github.com/alexander-emelyanenko/go-web-server/models"
	"github.com/alexander-emelyanenko/go-web-server/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

type Galleries struct {
	New *views.View
	gs  models.GalleryService
}
