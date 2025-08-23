package main

import "github.com/PPRAMANIK62/snippetbox/internal/models"

type templateData struct {
	Snippet *models.Snippet
	Snippets []*models.Snippet
}
