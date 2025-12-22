package handlers

import (
	"net/http"
	"strconv"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/middleware"
	"github.com/geekible-ltd/serviceframework/internal/services"
	"github.com/gin-gonic/gin"
)

type UserMaintenanceHandler struct {
	jwtSecret              string
	userMaintenanceService *services.UserMaintenanceService
}

func NewUserMaintenanceHandler(jwtSecret string, userMaintenanceService *services.UserMaintenanceService) *UserMaintenanceHandler {
	return &UserMaintenanceHandler{jwtSecret: jwtSecret, userMaintenanceService: userMaintenanceService}
}

func (h *UserMaintenanceHandler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/user-maintenance")

	api.POST("/reset-password-request", h.ResetPasswordRequest)
	api.POST("/reset-password", h.ResetPassword)
	api.POST("/verify-email", h.VerifyEmail)

	protected := api.Use(middleware.BearerAuthMiddleware(h.jwtSecret))
	{
		protected.DELETE("/user", h.DeleteUser)
		protected.PUT("/user", h.UpdateUser)
		protected.GET("/users/get-all", h.GetAllUsers)
		protected.GET("/users/get-roles", h.GetUserRoles)
	}
}

func (h *UserMaintenanceHandler) DeleteUser(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	userID, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if userID <= 0 {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("User ID is required"))
		return
	}

	currentUserID, err := strconv.Atoi(tokenDto.Sub)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if userID == currentUserID {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You cannot delete yourself"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantUser) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to delete this user"))
		return
	}

	err = h.userMaintenanceService.DeleteUser(tokenDto.TenantID, uint(userID))
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "User deleted successfully")
}

func (h *UserMaintenanceHandler) UpdateUser(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	var updateUserDTO frameworkdto.UserUpdateRequestDTO
	if err := c.ShouldBindJSON(&updateUserDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	currentUserID, err := strconv.Atoi(tokenDto.Sub)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid User ID format"))
		return
	}

	if updateUserDTO.UserID != uint(currentUserID) || tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to update this user"))
		return
	}

	err = h.userMaintenanceService.UpdateUser(tokenDto.TenantID, updateUserDTO.UserID, updateUserDTO)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "User updated successfully")
}

func (h *UserMaintenanceHandler) ResetPasswordRequest(c *gin.Context) {
	var resetPasswordRequestDTO frameworkdto.ResetPasswordRequestDTO
	if err := c.ShouldBindJSON(&resetPasswordRequestDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.SetResetPasswordToken(resetPasswordRequestDTO.Email)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Reset password request sent successfully")
}

func (h *UserMaintenanceHandler) ResetPassword(c *gin.Context) {
	var resetPasswordDTO frameworkdto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&resetPasswordDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.UpdateUserPassword(resetPasswordDTO.ResetToken, resetPasswordDTO.NewPassword)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Password reset successfully")
}

func (h *UserMaintenanceHandler) VerifyEmail(c *gin.Context) {
	var verifyEmailDTO frameworkdto.VerifyEmailDTO
	if err := c.ShouldBindJSON(&verifyEmailDTO); err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.BadRequest("Invalid request body"))
		return
	}

	err := h.userMaintenanceService.VerifyEmail(verifyEmailDTO.TenantID, verifyEmailDTO.UserID, verifyEmailDTO.Token)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusAccepted, nil, "Email verified successfully")
}

func (h *UserMaintenanceHandler) GetAllUsers(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get all users"))
		return
	}

	users, err := h.userMaintenanceService.GetAllUsersByTenantID(tokenDto.TenantID)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, users, "Users fetched successfully")
}

func (h *UserMaintenanceHandler) GetUserRoles(c *gin.Context) {
	tokenDto, err := frameworkutils.GetTokenDTO(c)
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.UnauthorizedError("Unauthorized"))
		return
	}

	if tokenDto.Role != string(frameworkconstants.UserRoleTenantAdmin) {
		frameworkutils.ErrorResponse(c, frameworkutils.Forbidden("You are not authorized to get this resource"))
		return
	}

	roles, err := h.userMaintenanceService.GetUserRoles()
	if err != nil {
		frameworkutils.ErrorResponse(c, frameworkutils.InternalServerError(err.Error()))
		return
	}

	frameworkutils.SuccessResponse(c, http.StatusOK, roles, "User roles fetched successfully")
}
