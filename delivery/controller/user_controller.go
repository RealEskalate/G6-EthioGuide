package controller

import (
	"EthioGuide/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(uc domain.IUserUsecase) *UserController {
	return &UserController{userUsecase: uc}
}

func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	password := req.Password
	userDetail := &domain.UserDetail{
		Username: req.Username,
		SubscriptionPlan: domain.SubscriptionNone,
		IsBanned: false,
		IsVerified: false,
	}
	account := &domain.Account{
		Name: req.Name,
		Email: req.Email,
		PasswordHash: password,
		Role: domain.RoleUser,

		UserDetail: userDetail,
	}
	// user := &domain.Account{
	// 	Username: req.Username,
	// 	Email:    req.Email,
	// 	Password: &password,
	// 	Role:     domain.RoleUser,
	// }
	// fmt.Println(account)
	err := ctrl.userUsecase.Register(c.Request.Context(), account)
	// fmt.Println(account)
	if err != nil {
		switch err {
		case domain.ErrEmailExists:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case domain.ErrPasswordTooShort, domain.ErrInvalidEmailFormat:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An internal error occurred"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": toUserResponse(account)})
}