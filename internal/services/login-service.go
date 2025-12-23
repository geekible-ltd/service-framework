package services

import (
	"time"

	frameworkconstants "github.com/geekible-ltd/serviceframework/framework-constants"
	frameworkdto "github.com/geekible-ltd/serviceframework/framework-dto"
	frameworkutils "github.com/geekible-ltd/serviceframework/framework-utils"
	"github.com/geekible-ltd/serviceframework/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService struct {
	cfg        *frameworkdto.FrameworkConfig
	userRepo   *repositories.UserRepository
	tenantRepo *repositories.TenantRepository
}

func NewLoginService(cfg *frameworkdto.FrameworkConfig, userRepo *repositories.UserRepository, tenantRepo *repositories.TenantRepository) *LoginService {
	return &LoginService{
		cfg:        cfg,
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
	}
}

func (s *LoginService) Login(loginRequest frameworkdto.LoginDTO, ipAddress string) (frameworkdto.LoginResponseDTO, error) {
	user, err := s.userRepo.GetByEmail(loginRequest.Email)
	if err != nil && err == gorm.ErrRecordNotFound {
		return frameworkdto.LoginResponseDTO{}, frameworkconstants.ErrUserNotFound
	} else if err != nil {
		return frameworkdto.LoginResponseDTO{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginRequest.Password)); err != nil {
		user.FailedLoginAttempts++
		if user.FailedLoginAttempts >= frameworkconstants.MaxFailedLoginAttempts {
			user.IsActive = false
		}
		if err := s.userRepo.Update(user); err != nil {
			return frameworkdto.LoginResponseDTO{}, err
		}
		return frameworkdto.LoginResponseDTO{}, frameworkconstants.ErrInvalidPassword
	}

	tenant, err := s.tenantRepo.GetByID(user.TenantID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return frameworkdto.LoginResponseDTO{}, frameworkconstants.ErrTenantNotFound
	} else if err != nil {
		return frameworkdto.LoginResponseDTO{}, err
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ipAddress
	user.FailedLoginAttempts = 0

	if err := s.userRepo.Update(user); err != nil {
		return frameworkdto.LoginResponseDTO{}, err
	}

	token, err := frameworkutils.GenerateJWT(user.ID, tenant.ID, user.Email, user.FirstName, user.LastName, user.Role, []byte(s.cfg.JWTSecret))
	if err != nil {
		return frameworkdto.LoginResponseDTO{}, err
	}

	return frameworkdto.LoginResponseDTO{
		Token: frameworkdto.BearerToken(token),
	}, nil
}
