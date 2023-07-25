package utility

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var saltLen = 8

var Counter int64

func init() {
	rand.Seed(time.Now().UnixNano())
}

func HashPassword(password string) (hashed string, salt string, err error) {
	salt = randomSalt()
	hash, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return string(hash), salt, nil
}

func PasswordIsValid(password, salt, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+salt))
	return err == nil
}

func randomSalt() string {
	var salt string
	for i := 0; i < saltLen; i++ {
		char := rand.Int31n(122) + 41
		salt += string(char)
	}

	return salt
}

func GetURLHash(id, url string) (string, error) {

	hash := md5.New()
	if _, err := hash.Write([]byte(id)); err != nil {
		return "", nil
	}
	md5Hash := hash.Sum([]byte(url))
	cntStr := fmt.Sprint(Counter)
	atomic.AddInt64(&Counter, 1) // Increment counter

	base64Str := base64.URLEncoding.EncodeToString(append(md5Hash, []byte(cntStr)...))
	return base64Str[:7], nil
}
