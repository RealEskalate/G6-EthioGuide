package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)


type GeminiController struct{
	geminiUseCase domain.IGeminiUseCase
}

func NewGeminiController(geminiUseCase domain.IGeminiUseCase) *GeminiController {
	return &GeminiController{
		geminiUseCase: geminiUseCase,
	}
}



func (gc *GeminiController) Translate( c *gin.Context){

	var request TranslateDTO
	if err := c.ShouldBindJSON(&request); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid request body"})
		return 
	}
	preferredLang := c.GetHeader("lang")
	if preferredLang == "" {
        preferredLang = "en" 
    }

	

	translated, err := gc.geminiUseCase.TranslateContent(c.Request.Context(), request.Content, preferredLang)

	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to translate content"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": translated})


}



