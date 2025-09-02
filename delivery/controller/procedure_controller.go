package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProcedureController struct {
	procedureUsecase domain.IProcedureUsecase
}

func NewProcedureController(pu domain.IProcedureUsecase) *ProcedureController {
	return &ProcedureController{
		procedureUsecase: pu,
	}
}

func (ctrl *ProcedureController) CreateProcedure(c *gin.Context) {
	var proc ProcedureCreateRequest
	if err := c.ShouldBindJSON(&proc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	domainProc := toDomainProcedure(&proc)
	err := ctrl.procedureUsecase.CreateProcedure(c.Request.Context(), domainProc)
	if err != nil {
		HandleError(c, err)
		return 
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Procedure created successfully", "user": domainProc})
} 