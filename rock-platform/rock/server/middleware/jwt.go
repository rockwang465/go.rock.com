package middleWare

import (
	"github.com/dgrijalva/jwt-go"
	"go.rock.com/rock-platform/rock/server/conf"
	"go.rock.com/rock-platform/rock/server/log"
	"time"
)

type Claim struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var jwtKey = []byte("Is_a_jwt_secret_key_from_rock_platform")

// generate token
func GenerateToken(userId int64, username, password string) (string, error) {
	config := conf.GetConfig()
	logger := log.GetLogger()
	nowTime := time.Now()

	userDefineExpire := config.Viper.GetDuration("server.tokenExpire")
	if userDefineExpire == 0 {
		logger.Warning("Not defined the server token expire time, please check.")
		userDefineExpire = time.Minute * 10 // Set the default token expire time to 10 minutes
	}

	expireTime := nowTime.Add(time.Minute + userDefineExpire)
	claim := &Claim{
		UserId:   userId,
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    "Rock Wang",
			Subject:   "Login token",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaims.SignedString(jwtKey)
	return token, err
}
