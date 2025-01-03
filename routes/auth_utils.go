package routes

import (
	"log"
	"time"

	"github.com/acatalepsy17/pigeon/config"
	"github.com/acatalepsy17/pigeon/models"
	"github.com/acatalepsy17/pigeon/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pborman/uuid"
	"gorm.io/gorm"
)

var cfg = config.GetConfig()
var SECRETKEY = []byte(cfg.JWTSecretKey)

type AccessTokenPayload struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

type RefreshTokenPayload struct {
	Data string `json:"data"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userId uuid.UUID, username string) string {
	expirationTime := time.Now().Add(time.Duration(cfg.AccessTokenExpireMinutes) * time.Minute)
	payload := AccessTokenPayload{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(SECRETKEY)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Access token: ", err)
	}
	return tokenString
}

func GenerateRefreshToken() string {
	expirationTime := time.Now().Add(time.Duration(cfg.RefreshTokenExpireMinutes) * time.Minute)
	payload := RefreshTokenPayload{
		Data: utils.GetRandomString(10),
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, payload)
	// Create the JWT string
	tokenString, err := token.SignedString(SECRETKEY)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		log.Fatal("Error Generating Refresh token: ", err)
	}
	return tokenString
}

func DecodeAccessToken(token string, db *gorm.DB) (*models.User, *string) {
	claims := &AccessTokenPayload{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
	})
	tokenErr := "Auth Token is Invalid or Expired!"
	if err != nil {
		return nil, &tokenErr
	}
	if !tkn.Valid {
		return nil, &tokenErr
	}
	user := models.User{Access: &token}
	// Fetch User model object
	result := db.Where(user).Joins("CityObj").Joins("CityObj.RegionObj").Joins("CityObj.CountryObj").Joins("AvatarObj").Take(&user, claims.UserId)
	if result.Error != nil {
		return nil, &tokenErr
	}
	return &user, nil
}

func DecodeRefreshToken(token string) bool {
	claims := &RefreshTokenPayload{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("JWT Error: ", "Invalid Signature")
		} else {
			log.Println("JWT Error: ", err)
		}
		return false
	}
	if !tkn.Valid {
		log.Println("Invalid Refresh Token")
		return false
	}
	return true
}
