package lib

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	RDB      *redis.Client
	DB       *gorm.DB
	S3       *s3.Client
	DEEPSEEK *openai.Client
)

func InstallRedis() {
	redisUrl, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		panic("REDIS_URL is not set")
	}
	opt, _ := redis.ParseURL(redisUrl)
	RDB = redis.NewClient(opt)
}

func InstallDB() {
	dsn, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		panic("DATABASE_URL is not set")
	}
	if strings.HasPrefix(dsn, "mysql") {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
		DB = db
		return
	}
	if strings.HasPrefix(dsn, "postgres") {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic("failed to get sql.DB: " + err.Error())
		}
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
		DB = db
		return
	}
	panic("unsupported database")
}

func InstallS3FromEnv() {
	var (
		endpoint        = os.Getenv("AWS_ENDPOINT_URL_S3")
		region          = os.Getenv("AWS_REGION")
		accessKeyId     = os.Getenv("AWS_ACCESS_KEY_ID")
		accessKeySecret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	)
	S3 = getS3Client(accessKeyId, accessKeySecret, region, endpoint)
}

func InstallDeepseekFromEnv() {
	DEEPSEEK = openai.NewClient(
		option.WithBaseURL(os.Getenv("DEEPSEEK_BASE_URL")),
		option.WithAPIKey(os.Getenv("DEEPSEEK_API_KEY")),
	)
}

func getS3Client(accessKeyId, accessKeySecret, region, endpoint string) *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion(region),
	)
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	return client
}
