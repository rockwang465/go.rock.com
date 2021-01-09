package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-gomail/gomail"
	"go.rock.com/rock-platform/rock/server/conf"
	"golang.org/x/crypto/pbkdf2"
	"math/rand"
	"time"
)

const (
	SALT_LENGTH          = 16
	SALT_RESOURCE_LETTER = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ITER_COUNT           = 4096
	KEY_LENGTH           = 16

	NEW_USER_EMAIL_CONTENT = "用户创建提示"
	NEW_USER_EMAIL_SUBJECT = "系统邮件，请勿回复！\n%s，您好：\n    管理员为您创建了Rock平台的账户：\n        用户名：%s\n        密码：%s \n    请您尽快登录Rock平台并修改初始密码。"
)

// 生成随机salt字符串
func GenerateSalt() string {
	bytesStr := []byte(SALT_RESOURCE_LETTER)
	bytesRandom := []byte{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < SALT_LENGTH; i++ {
		num := rand.Intn(len(SALT_RESOURCE_LETTER))
		bytesRandom = append(bytesRandom, bytesStr[num])
	}
	salt := string(bytesRandom)
	fmt.Println(salt)
	return salt
}

//  加密密码
func EncryptPwd(password, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), ITER_COUNT, KEY_LENGTH, sha256.New)
	pwd := hex.EncodeToString(dk)
	return pwd

}

// 获取过期的秒数
func GetExpireDuration() int {
	config := conf.GetConfig()
	duration := config.Viper.GetDuration("server.tokenExpire") // default 10 minutes
	return int(duration / time.Second)
}

// Send New Password Eamil
func SendNewPwdEmail(userName, destEmail, userPwd string) error {
	config := conf.GetConfig()
	user := config.Viper.GetString("email.user")
	pwd := config.Viper.GetString("email.password")
	port := config.Viper.GetInt("email.smtp.port")
	addr := config.Viper.GetString("email.smtp.addr")
	m := gomail.NewMessage()
	content := fmt.Sprintf(NEW_USER_EMAIL_SUBJECT, userName, userName, userPwd)
	m.SetHeader("From", user)
	m.SetHeader("To", destEmail)
	m.SetHeader("Subject", NEW_USER_EMAIL_CONTENT)
	m.SetBody("text/plain", content)

	d := gomail.NewDialer(addr, port, user, pwd)
	if err := d.DialAndSend(m); err != nil {
		errMsg := fmt.Sprintf("go mail DialAndSend failed , %v\n", err)
		newErr := NewRockError(500, 50000003, errMsg)
		panic(newErr)
	}
	return nil
}
