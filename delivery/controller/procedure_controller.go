package controller

import (
	"EthioGuide/domain"
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
	procedure, err := pc.procedureUsecase.GetProcedureByID(ctx, id)
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
	err := pc.procedureUsecase.UpdateProcedure(ctx, id, dto.FromDTOToDomainProcedure())
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Updated Procedure Successfully.")
}

func (pc *ProcedureController) DeleteProcedure(ctx *gin.Context) {
	id := ctx.Param("id")
	err := pc.procedureUsecase.DeleteProcedure(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, "Deleted Procedure Successfully.")
}
