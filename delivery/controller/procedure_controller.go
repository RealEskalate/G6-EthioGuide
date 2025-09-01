package controller

import (
	"EthioGuide/domain"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProcedureController struct {
	procedureUsecase domain.IProcedureUseCase
}

func NewProcedureController(procedureUsecase domain.IProcedureUseCase) *ProcedureController {
	return &ProcedureController{
		procedureUsecase: procedureUsecase,
	}
}

func (pc *ProcedureController) GetProcedureByID(ctx *gin.Context) {
	id := ctx.Param("id")
	procedure, err := pc.procedureUsecase.GetProcedureByID(context.Background(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, FromDomainProcedureToDTO(procedure))
}

func (pc *ProcedureController) UpdateProcedure(ctx *gin.Context) {
	id := ctx.Param("id")
	var dto ProcedureDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := pc.procedureUsecase.UpdateProcedure(context.Background(), id, dto.FromDTOToDomainProcedure())
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Updated Procedure Successfully.")
}