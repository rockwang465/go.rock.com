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
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("Is_a_jwt_secret_key_from_rock_platform")

// generate token
func GenerateToken(userId int64, username, password, role string) (string, error) {
	config := conf.GetConfig()
	logger := log.GetLogger()
	nowTime := time.Now()

	userDefineExpire := config.Viper.GetDuration("server.tokenExpire")
	if userDefineExpire == 0 {
		logger.Warning("Not defined the server token expire time, please check.")
		userDefineExpire = time.Minute * 10 // Set the default token expire time to 10 minutes
	}

	// generate admin token is 100 years
	//var expireTime time.Time
	//if username == "admin" {
	//	fmt.Println("admin account")
	//	expireTime = nowTime.Add(time.Hour * 24 * 365 * 100)  // 2120-12-01
	//	fmt.Printf("expireTime: %v\n", expireTime)
	//}

	expireTime := nowTime.Add(time.Minute + userDefineExpire)
	claim := &Claim{
		UserId:   userId,
		Username: username,
		Password: password,
		Role:     role,
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

// parse token body
func ParseToken(tokenString string) (*jwt.Token, *Claim, error) {
	var claims = &Claim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
