package controller

import (
	"EthioGuide/domain"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User Role not found in token"})
		return
	}

	proc.OrganizationID = userID.(string)

	domainProc := toDomainProcedure(&proc)
	err := ctrl.procedureUsecase.CreateProcedure(c.Request.Context(), domainProc, userID.(string), userRole.(domain.Role))
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

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User Role not found in token"})
		return
	}

	err := pc.procedureUsecase.UpdateProcedure(ctx.Request.Context(), id, &procedure, userID.(string), userRole.(domain.Role))
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

	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userRole, exists := ctx.Get("userRole")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User Role not found in token"})
		return
	}

	err := pc.procedureUsecase.DeleteProcedure(ctx.Request.Context(), id, userID.(string), userRole.(domain.Role))
	if err != nil {
		HandleError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// @Summary      Search and Filter Procedures
// @Description  Search and filter procedures with pagination, sorting, and various filters.
// @Tags         Procedures
// @Accept       json
// @Produce      json
// @Param        page              query     int     false  "Page number (default 1)"
// @Param        limit             query     int     false  "Results per page (default 10)"
// @Param        logic             query     string  false  "Global logic: AND or OR (default AND)"
// @Param        name              query     string  false  "Search by procedure name"
// @Param        organizationID    query     string  false  "Filter by organization ID"
// @Param        groupID           query     string  false  "Filter by group ID"
// @Param        minFee            query     number  false  "Minimum fee"
// @Param        maxFee            query     number  false  "Maximum fee"
// @Param        minProcessingDays query     int     false  "Minimum processing days"
// @Param        maxProcessingDays query     int     false  "Maximum processing days"
// @Param        startDate         query     string  false  "Start date (RFC3339 format)"
// @Param        endDate           query     string  false  "End date (RFC3339 format)"
// @Param        sortBy            query     string  false  "Sort by field (e.g. createdAt, fee, processingTime)"
// @Param        sortOrder         query     string  false  "Sort order: ASC or DESC (default DESC)"
// @Success      200  {object}  PaginatedProcedureResponse
// @Failure      400  {object}  map[string]string "Invalid parameter"
// @Failure      500  {object}  map[string]string "Server error"
// @Router       /procedures/search [get]

func (pc *ProcedureController) SearchAndFilter(c *gin.Context) {
	options := domain.ProcedureSearchFilterOptions{
		GlobalLogic: domain.GlobalLogicAND, // default
		SortOrder:   domain.SortDesc,       // default newest first
	}

	// --- Pagination ---
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'page' parameter"})
		return
	}
	options.Page = page

	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'limit' parameter"})
		return
	}
	options.Limit = limit

	// --- Global logic ---
	if strings.ToUpper(c.Query("logic")) == string(domain.GlobalLogicOR) {
		options.GlobalLogic = domain.GlobalLogicOR
	}

	// --- Search ---
	if name := c.Query("name"); name != "" {
		options.Name = &name
	}

	// --- Filters ---
	if orgID, ok := c.GetQuery("organizationID"); ok {
		options.OrganizationID = &orgID
	}
	if groupID, ok := c.GetQuery("groupID"); ok {
		options.GroupID = &groupID
	}

	// --- Fee range ---
	if minFeeStr := c.Query("minFee"); minFeeStr != "" {
		if f, err := strconv.ParseFloat(minFeeStr, 64); err == nil {
			options.MinFee = &f
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'minFee' parameter"})
			return
		}
	}
	if maxFeeStr := c.Query("maxFee"); maxFeeStr != "" {
		if f, err := strconv.ParseFloat(maxFeeStr, 64); err == nil {
			options.MaxFee = &f
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'maxFee' parameter"})
			return
		}
	}

	// --- Processing time range ---
	if minDaysStr := c.Query("minProcessingDays"); minDaysStr != "" {
		if d, err := strconv.Atoi(minDaysStr); err == nil {
			options.MinProcessingDays = &d
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'minProcessingDays' parameter"})
			return
		}
	}
	if maxDaysStr := c.Query("maxProcessingDays"); maxDaysStr != "" {
		if d, err := strconv.Atoi(maxDaysStr); err == nil {
			options.MaxProcessingDays = &d
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'maxProcessingDays' parameter"})
			return
		}
	}

	// --- Date range ---
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if t, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			options.StartDate = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'startDate' format. Use RFC3339"})
			return
		}
	}
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if t, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			options.EndDate = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid 'endDate' format. Use RFC3339"})
			return
		}
	}

	// --- Sorting ---
	options.SortBy = c.Query("sortBy") // e.g. "createdAt", "fee", "processingTime"
	if strings.ToUpper(c.Query("sortOrder")) == string(domain.SortAsc) {
		options.SortOrder = domain.SortAsc
	}

	// --- Call usecase ---
	procedures, total, err := pc.procedureUsecase.SearchAndFilter(c.Request.Context(), options)
	if err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPaginatedProcedureResponse(procedures, total, options.Page, options.Limit))
}
