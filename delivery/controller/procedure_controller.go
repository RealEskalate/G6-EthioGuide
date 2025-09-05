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

// @Summary      Create Procedure
// @Description  Create new procedure.
// @Tags         Procedures
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        request body ProcedureCreateRequest true "Procedure Detail"
// @Success      200 {object} domain.Procedure "Procedure Created"
// @Failure      400 {string}  "Invalid request"
// @Failure      401 {string}  "Unauthorized"
// @Failure      500 {string}  "Server error"
// @Router       /procedures [post]
func (ctrl *ProcedureController) CreateProcedure(c *gin.Context) {
	var proc ProcedureCreateRequest
	if err := c.ShouldBindJSON(&proc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	proc.OrganizationID = userID.(string)

	domainProc := toDomainProcedure(&proc)
	err := ctrl.procedureUsecase.CreateProcedure(c.Request.Context(), domainProc)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Procedure created successfully", "user": domainProc})
}

// @Summary      Get Procedure by ID
// @Description  Retrieves a single procedure by its unique ID.
// @Tags         Procedures
// @Produce      json
// @Param        id   path      string  true  "Procedure ID"
// @Success      200  {object}  domain.Procedure
// @Failure      404  {string} Procedure not found
// @Router       /procedures/{id} [get]
func (pc *ProcedureController) GetProcedureByID(ctx *gin.Context) {
	id := ctx.Param("id")
	procedure, err := pc.procedureUsecase.GetProcedureByID(ctx.Request.Context(), id)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, procedure)
}

// @Summary      Update Procedure
// @Description  Updates an existing procedure. Requires admin or organization ownership.
// @Tags         Procedures
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Procedure ID"
// @Param        request body domain.Procedure true "Procedure Update Data"
// @Success      200  {string}  success
// @Failure      400  {string}  Invalid request body
// @Failure      401  {string}  Unauthorized
// @Failure      403  {string}  Permission Denied
// @Failure      404  {string}  Procedure not found
// @Router       /procedures/{id} [patch]
func (pc *ProcedureController) UpdateProcedure(ctx *gin.Context) {
	id := ctx.Param("id")
	var procedure domain.Procedure
	if err := ctx.ShouldBindJSON(&procedure); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := pc.procedureUsecase.UpdateProcedure(ctx.Request.Context(), id, &procedure)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, "Updated Procedure Successfully.")
}

// @Summary      Delete Procedure
// @Description  Deletes an existing procedure. Requires admin or organization ownership.
// @Tags         Procedures
// @Param        Authorization header string true "Bearer token"
// @Param        id path string true "Procedure ID"
// @Success      204  "No Content"
// @Failure      401  {string} Unauthorized
// @Failure      403  {string} Permission Denied
// @Failure      404  {string} Procedure not found
// @Router       /procedures/{id} [delete]
func (pc *ProcedureController) DeleteProcedure(ctx *gin.Context) {
	id := ctx.Param("id")
	err := pc.procedureUsecase.DeleteProcedure(ctx.Request.Context(), id)
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
