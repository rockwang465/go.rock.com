package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"math/rand"
	"time"
)

const (
	SALT_LENGTH          = 16
	SALT_RESOURCE_LETTER = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ITER_COUNT           = 4096
	KEY_LENGTH           = 16
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
