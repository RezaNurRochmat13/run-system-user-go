package authService

import (
	"errors"
	"runs-system-user-go/database"
	authModel "runs-system-user-go/module/auth/model"
	userModel "runs-system-user-go/module/user/model"
	"runs-system-user-go/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(userRegister *authModel.UserRegister) error {
	db := database.DB

	var user userModel.User

	hash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 12)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	user.Name = userRegister.Name
	user.Email = userRegister.Email
	user.PasswordHash = string(hash)

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func LoginUser(userLogin *authModel.UserLogin) (authModel.AuthResponse, error) {
	db := database.DB
	var user userModel.User

	if err := db.First(&user, "email = ?", userLogin.Email).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return authModel.AuthResponse{}, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userLogin.Password)); err != nil {
		return authModel.AuthResponse{}, err
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		return authModel.AuthResponse{}, err
	}

	return authModel.AuthResponse{Token: token}, nil
}
