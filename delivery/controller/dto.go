package controller

type TranslateDTO struct {

	Content string `json:"content" binding:"required"`
}