package implementation

import (
	"final/internal/domain"
	"final/internal/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (uc *UserUseCase) Register(inputUser *domain.InputUser) (*domain.User, error) {
	if !checkUserData(inputUser) {
		return nil, domain.ErrInvalidCredentials
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*inputUser.Password), bcrypt.DefaultCost)
	hashStr := string(hash)
	if err != nil {
		log.Error(err)
		return nil, domain.ErrRegisterUser
	}

	user := &domain.User{
		Email: inputUser.Email,
		Role:  inputUser.Role,
	}
	user.Password = &hashStr
	return uc.userRepository.Register(user)
}

func checkUserData(user *domain.InputUser) bool {
	if user.Email == nil || user.Password == nil {
		return false
	}
	return true
}

func (uc *UserUseCase) Login(inputUser *domain.InputUser) (*domain.User, error) {
	if !checkUserData(inputUser) {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := uc.userRepository.GetUserByEmail(*inputUser.Email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, domain.ErrFindUser
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*inputUser.Password)); err != nil {
		log.Error(err)
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}
