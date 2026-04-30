package app

import (
	"clean_architecture_go/internal/app/services"
	"clean_architecture_go/internal/app/services/background_job_service"
	jwtservice "clean_architecture_go/internal/app/services/jwt_service"
	passwordservice "clean_architecture_go/internal/app/services/password_service"
	"log"
	"os"
)

type ServiceProvider struct {
	tokenService         services.TokenService
	pwdService           services.PasswordService
	backgroundJobService services.BackgroundJobService
}

var serviceProvider *ServiceProvider

func InitServices() {
	secretKey := os.Getenv("JWT_SECRET")

	serviceProvider = &ServiceProvider{
		tokenService:         jwtservice.NewJWTTokenService(secretKey),
		pwdService:           passwordservice.NewBcryptPasswordService(),
		backgroundJobService: background_job_service.NewJobQueue(5, 100),
	}

	log.Println("✅ Services initialized")
}

func GetTokenService() services.TokenService {
	if serviceProvider == nil {
		panic("Services not initialized. Call InitServices() first")
	}
	return serviceProvider.tokenService
}

func GetPwdService() services.PasswordService {
	if serviceProvider == nil {
		panic("Services not initialized. Call InitServices() first")
	}
	return serviceProvider.pwdService
}

func GetBKService() services.BackgroundJobService {
	if serviceProvider == nil {
		panic("Services not initialized. Call InitServices() first")
	}
	return serviceProvider.backgroundJobService
}

func AddApp() {
	InitServices()
}
