package implementation

import (
	"errors"
	"final/internal/domain"
	"final/internal/repository"
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
		return nil, errors.New("user data is invalid")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(*inputUser.Password), bcrypt.DefaultCost)
	hashStr := string(hash)
	if err != nil {
		return nil, err
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
		return nil, errors.New("user data is invalid")
	}
	user, err := uc.userRepository.GetUserByEmail(*inputUser.Email)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*inputUser.Password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return user, nil
}
