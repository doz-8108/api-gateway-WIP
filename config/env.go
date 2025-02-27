package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	JWT_SIGN_SECRET         string
	MYSQL_DSN               string
	USER_ID_ENCODE_INIT_NUM string
	USER_ID_ENCODE_ALPHB    string
	EMAIL_HOST              string
	REDIS_ADDR              string
	REDIS_PASSWORD          string
	REDIS_DB                int
	MAILTRAP_API_ENDPOINT   string
	MAILTRAP_API_TOKEN      string
	MAILTRAP_EMAIL_HOST     string
	MAILTRAP_TEMPLATE_UUID  string
	TLS_CERT_FILE_PATH      string
	TLS_KEY_FILE_PATH       string
}

func LoadEnv() EnvVars {
	godotenv.Load()

	userIdInitNum := os.Getenv("USER_ID_INIT_NUM")
	userIdAlphabet := os.Getenv("USER_ID_ALPHB")
	jwtSignSecret := os.Getenv("JWT_SIGN_SECRET")

	mysqlDsn := os.Getenv("MYSQL_DSN")

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	mailtrapApiEndpoint := os.Getenv("MAILTRAP_API_ENDPOINT")
	mailtrapApiToken := os.Getenv("MAILTRAP_API_TOKEN")
	mailtrapEmailHost := os.Getenv("MAILTRAP_EMAIL_HOST")
	mailtrapTemplateUuid := os.Getenv("MAILTRAP_TEMPLATE_UUID")

	tlsCertFilePath := os.Getenv("TLS_CERT_FILE_PATH")
	tlsKeyFilePath := os.Getenv("TLS_KEY_FILE_PATH")

	if err != nil {
		fmt.Println("Failed to load environment variables")
		panic(err)
	}

	return EnvVars{
		MYSQL_DSN:               mysqlDsn,
		USER_ID_ENCODE_INIT_NUM: userIdInitNum,
		USER_ID_ENCODE_ALPHB:    userIdAlphabet,
		REDIS_ADDR:              redisAddr,
		REDIS_PASSWORD:          redisPassword,
		REDIS_DB:                redisDb,
		MAILTRAP_API_ENDPOINT:   mailtrapApiEndpoint,
		MAILTRAP_API_TOKEN:      mailtrapApiToken,
		MAILTRAP_EMAIL_HOST:     mailtrapEmailHost,
		MAILTRAP_TEMPLATE_UUID:  mailtrapTemplateUuid,
		JWT_SIGN_SECRET:         jwtSignSecret,
		TLS_CERT_FILE_PATH:      tlsCertFilePath,
		TLS_KEY_FILE_PATH:       tlsKeyFilePath,
	}
}
