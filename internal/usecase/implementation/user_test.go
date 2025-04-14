package implementation

import (
	"final/internal/domain"
	auth "final/internal/midleware"
	"final/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	id := "123e4567-e89b-12d3-a456-426614174000"
	email := "user@example.com"
	role := "employee"
	password := "password"

	t.Run("Successful_Registration", func(t *testing.T) {
		inputUser := &domain.InputUser{Email: &email, Role: &role, Password: &password}
		hash, err := bcrypt.GenerateFromPassword([]byte(*inputUser.Password), bcrypt.DefaultCost)
		hashStr := string(hash)
		if err != nil {
			log.Error(err)
		}

		user := &domain.User{
			ID:       &id,
			Email:    &email,
			Role:     &role,
			Password: &hashStr,
		}

		mockRepo.EXPECT().Register(gomock.Any()).Return(user, nil)

		result, err := useCase.Register(inputUser)
		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})

	t.Run("Failure_Invalid_Credentials", func(t *testing.T) {
		inputUser := &domain.InputUser{Role: &role, Password: &password}

		result, err := useCase.Register(inputUser)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Nil(t, result)
	})
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_repository.NewMockUserRepository(ctrl)
	useCase := NewUserUseCase(mockRepo)

	if err := godotenv.Load("../../../.env"); err != nil {
		log.Error("error loading .env file")
	}
	err := auth.InitAuthFromConfig()
	if err != nil {
		log.Error(err)
		return
	}

	password := "password"
	email := "user@example.com"
	role := "employee"
	t.Run("success", func(t *testing.T) {
		userInput := domain.InputUser{
			Role:     &role,
			Email:    &email,
			Password: &password,
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error(err)
		}
		hashStr := string(hash)
		user := &domain.User{
			Email:    &email,
			Password: &hashStr,
			Role:     &role,
		}

		mockRepo.EXPECT().GetUserByEmail(email).Return(user, nil)

		result, err := useCase.Login(&userInput)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user.Email, result.Email)
		assert.Equal(t, user.Role, result.Role)
	})

	t.Run("Failure_Invalid_Credentials", func(t *testing.T) {
		inputUser := &domain.InputUser{
			Password: &password,
		}

		result, err := useCase.Login(inputUser)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_User_Not_Found", func(t *testing.T) {
		userInput := domain.InputUser{
			Role:     &role,
			Email:    &email,
			Password: &password,
		}

		mockRepo.EXPECT().GetUserByEmail(email).Return(nil, nil)

		result, err := useCase.Login(&userInput)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrFindUser, err)
		assert.Nil(t, result)
	})

	t.Run("Failure_Bcrypt_Compare_Error", func(t *testing.T) {
		userInput := domain.InputUser{
			Role:     &role,
			Email:    &email,
			Password: &password,
		}

		hash, _ := bcrypt.GenerateFromPassword([]byte("wrongpassword"), bcrypt.DefaultCost)
		hashStr := string(hash)

		user := &domain.User{
			Email:    &email,
			Password: &hashStr,
			Role:     &role,
		}

		mockRepo.EXPECT().GetUserByEmail(email).Return(user, nil)

		result, err := useCase.Login(&userInput)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Nil(t, result)
	})
}
