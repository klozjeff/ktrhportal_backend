package utilities

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var OTPCHARS = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GoDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading the .env file")
	}

	return os.Getenv(key)
}
func Validate(err error) []string {
	if verr, ok := err.(validator.ValidationErrors); ok {
		return simple(verr)
	}

	return []string{
		err.Error() + ".",
	}
}

func simple(verr validator.ValidationErrors) []string {
	var errs []string

	for _, f := range verr {
		err := f.ActualTag()
		if f.Param() != "" {
			err = fmt.Sprintf("%s=%s", err, f.Param())
		}
		errs = append(errs, fmt.Sprintf("%s is %s", f.Field(), err))
	}

	return errs
}

func GeneratePassword(password string) string {
	secrete, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(secrete)
}

func GenerateOTP(length int) (string, error) {

	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = OTPCHARS[int(b[i])%len(OTPCHARS)]
	}
	return string(b), nil
}

func SetCookie(c *gin.Context, name string, value string, expiration time.Time) {
	cookie := buildCookie(name, value, expiration.Second())
	http.SetCookie(c.Writer, cookie)
}

func ClearCookie(c *gin.Context, name string) {
	cookie := buildCookie(name, "", -1)

	http.SetCookie(c.Writer, cookie)
}

func buildCookie(name string, value string, expires int) *http.Cookie {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   expires,
	}

	return cookie
}

func GenerateAutoIncrementNumber(modelCount int) string {
	today := time.Now()
	day := today.Day()
	month := today.Month()
	monthString := fmt.Sprintf("%02d", month)
	year := today.Year()
	return fmt.Sprintf("%d%s%s%003d", day, monthString, strconv.Itoa(year)[2:], modelCount)
}
