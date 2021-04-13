package utils

import (
	"github.com/dgrijalva/jwt-go"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/log"
	"time"
)

const (
	AdminRole string = "admin"
)

type Claim struct {
	UserId     int64  `json:"user_id"`
	Username   string `json:"username"`
	DroneToken string `json:"drone_token"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("Is_a_jwt_secret_key_from_rock_platform")

// generate token
func GenerateToken(userId int64, username, droneToken, password, role string) (string, error) {
	config := conf.GetConfig()
	logger := log.GetLogger()
	nowTime := time.Now()

	userDefineExpire := config.Viper.GetDuration("server.tokenExpire")
	if userDefineExpire == 0 {
		logger.Warning("Not defined the server token expire time, please check.")
		userDefineExpire = time.Minute * 10 // Set the default token expire time to 10 minutes
	}

	// generate admin(galaxias_api_token) token is 100 years
	var expireTime time.Time
	if username == "galaxias_api_token" {
		expireTime = nowTime.Add(time.Hour * 24 * 365 * 100) // 2121-03-19 11:04:32
		logger.Debugf("username: %v, expire_time: %v", username, expireTime)
	} else {
		expireTime = nowTime.Add(time.Minute + userDefineExpire)
		logger.Debugf("username: %v, expire_time: %v", username, expireTime)
	}

	//expireTime = nowTime.Add(time.Minute + userDefineExpire)
	claim := &Claim{
		UserId:     userId,
		Username:   username,
		DroneToken: droneToken,
		Password:   password,
		Role:       role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    "Rock Wang",
			Subject:   "Login token",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaims.SignedString(jwtKey)
	if username == "galaxias_api_token" {
		logger.Debugf("username: %v, token: %v", username, token)
	}
	return token, err
}

// parse token body
func ParseToken(tokenString string) (*jwt.Token, *Claim, error) {
	var claims = &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
