package schemas

import "github.com/acatalepsy17/pigeon/models"

func ConvertUsers(users []models.User) []UserDataSchema {
	convertedUsers := []UserDataSchema{}
	for i := range users {
		user := UserDataSchema{}.Init(users[i])
		convertedUsers = append(convertedUsers, user)
	}
	return convertedUsers
}
