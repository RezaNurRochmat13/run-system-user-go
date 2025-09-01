package userService

import (
	"encoding/json"
	"errors"
	"fmt"
	"runs-system-user-go/database"
	userModel "runs-system-user-go/module/user/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetPaginatedUsers(page, limit int, status, search string) ([]userModel.User, int, error) {
	db := database.DB
	redisDB := database.RedisDB
	redisTTL := database.RedisCacheTTL

	offset := (page - 1) * limit
	cacheKey := fmt.Sprintf("users:page:%d:limit:%d:status:%s:search:%s", page, limit, status, search)

	// Try fetching from cache
	cachedUsers, err := redisDB.Get(redisDB.Context(), cacheKey).Result()
	if err == nil {
		var users []userModel.User
		if json.Unmarshal([]byte(cachedUsers), &users) == nil {
			return users, len(users), nil
		}
	}

	// Fetch from database
	var users []userModel.User
	query := db.Offset(offset).Limit(limit)

	result := query.Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// Cache results
	data, err := json.Marshal(users)
	if err == nil {
		redisDB.Set(redisDB.Context(), cacheKey, data, redisTTL)
	}

	return users, len(users), nil
}

func GetUserByID(id string) (userModel.User, error) {
	db := database.DB
	var user userModel.User

	if err := db.First(&user, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}

	return user, nil
}


func CreateUser(user *userModel.User) error {
	db := database.DB
	user.ID = uuid.New()

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func UpdateUser(id string, data map[string]interface{}) (userModel.User, error) {
	db := database.DB
	var user userModel.User

	if err := db.First(&user, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errors.New("user not found")
	}

	if err := db.Model(&user).Updates(data).Error; err != nil {
		return user, err
	}

	return user, nil
}

func DeleteUser(id string) error {
	db := database.DB
	var user userModel.User

	if err := db.First(&user, "id = ?", id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}

	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
