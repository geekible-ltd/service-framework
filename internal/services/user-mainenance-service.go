package services

import (
	"time"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserMaintenanceService struct {
	userRepo *repositories.UserRepository
}

func NewUserMaintenanceService(userRepo *repositories.UserRepository) *UserMaintenanceService {
	return &UserMaintenanceService{userRepo: userRepo}
}

func (s *UserMaintenanceService) DeleteUser(tenantID uint, userID uint) error {
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return err
	}

	s.userRepo.Delete(user)

	return nil
}

func (s *UserMaintenanceService) UpdateUser(tenantID uint, userID uint, userDTO frameworkdto.UserUpdateRequestDTO) error {
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return err
	}

	user.FirstName = userDTO.FirstName
	user.LastName = userDTO.LastName
	user.Email = userDTO.Email
	user.Role = userDTO.Role
	user.IsActive = userDTO.IsActive
	user.IsEmailVerified = userDTO.IsEmailVerified

	s.userRepo.Update(user)

	return nil
}

func (s *UserMaintenanceService) SetResetPasswordToken(email string) error {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return err
	}

	user.ResetPasswordToken = uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)
	user.ResetPasswordTokenExpiresAt = &expiresAt

	s.userRepo.Update(user)

	return nil
}

func (s *UserMaintenanceService) UpdateUserPassword(resetToken string, password string) error {
	user, err := s.userRepo.GetByResetPasswordToken(resetToken)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(passwordHash)

	s.userRepo.Update(user)

	return nil
}

func (s *UserMaintenanceService) VerifyEmail(tenantID uint, userID uint, token string) error {
	user, err := s.userRepo.GetByID(tenantID, userID)
	if err != nil {
		return err
	}

	user.IsEmailVerified = true

	s.userRepo.Update(user)

	return nil
}

func (s *UserMaintenanceService) GetAllUsersByTenantID(tenantID uint) ([]frameworkdto.GetUsersResponseDTO, error) {
	users, err := s.userRepo.GetAll(tenantID)
	if err != nil {
		return nil, err
	}

	usersDTO := make([]frameworkdto.GetUsersResponseDTO, len(users))
	for i, user := range users {
		usersDTO[i] = frameworkdto.GetUsersResponseDTO{
			UserID:    user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
	}

	return usersDTO, nil
}

func (s *UserMaintenanceService) GetUserRoles() ([]frameworkdto.GetUserRoles, error) {
	roles := []frameworkdto.GetUserRoles{
		{Role: string(frameworkconstants.UserRoleTenantAdmin)},
		{Role: string(frameworkconstants.UserRoleTenantUser)},
	}
	return roles, nil
}
