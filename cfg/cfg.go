package configs

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	//"github.com/joho/godotenv"
)

type Config struct {
	DbName     string
	DbPort     string
	DbHost     string
	DbUsername string
	DbPassword string

	APPPort string
	ENV     string

	JWTSecret  string
	BcryptSalt int

	//s3 to upload, all uploaded files will available just for only a day
	S3_ID          string
	S3_SECRET_KEY  string
	S3_BUCKET_NAME string
	S3_REGION      string
}

var (
	config     Config
	configOnce sync.Once
)

func LoadConfig() (Config, error) {
	var err error

	configOnce.Do(func() {
		// Load the .env file before reading the environment variables.
		// err = godotenv.Load(".env")
		// if err != nil {
		// 	fmt.Printf("No .env file found or error loading it: %v\n", err)
		// }

		config = Config{
			DbName:     GetEnvOrDefault("DB_NAME", "default_db_name"),
			DbHost:     GetEnvOrDefault("DB_HOST", "localhost"),
			DbPort:     GetEnvOrDefault("DB_PORT", "5432"),
			DbUsername: GetEnvOrDefault("DB_USERNAME", "user"),
			DbPassword: GetEnvOrDefault("DB_PASSWORD", "password"),

			APPPort: GetEnvOrDefault("APP_PORT", "8080"),
			ENV:     GetEnvOrDefault("ENV", "development"),

			JWTSecret: GetEnvOrDefault("JWT_SECRET", "sec"),

			S3_ID:          GetEnvOrDefault("S3_ID", ""),
			S3_SECRET_KEY:  GetEnvOrDefault("S3_SECRET_KEY", ""),
			S3_BUCKET_NAME: GetEnvOrDefault("S3_BUCKET_NAME", ""),
			S3_REGION:      GetEnvOrDefault("S3_REGION", ""),
		}

		config.BcryptSalt, err = strconv.Atoi(GetEnvOrDefault("BCRYPT_SALT", "10"))
		if err != nil {
			err = fmt.Errorf("failed to convert BCRYPT_SALT to int: %v", err)
		}
	})

	return config, err
}

// If the variable is empty, the default value is returned.
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
