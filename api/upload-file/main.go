package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func setUpMinio(ctx context.Context, bucket string) (*minio.Client, error) {
	endpoint := "minio:9000"
	accessKey, ok := os.LookupEnv("MINIO_ACCESS_KEY")
	if !ok {
		return nil, fmt.Errorf("MINIO_ACCESS_KEY env variable is required")
	}

	secretKey, ok := os.LookupEnv("MINIO_SECRET_KEY")
	if !ok {
		return nil, fmt.Errorf("MINIO_SECRET_KEY env variable is required")
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	policy := `{
		"Version":"2012-10-17",
		"Statement":[
			{
				"Effect":"Allow",
				"Principal":{"AWS":["*"]},
				"Action":["s3:GetBucketLocation","s3:ListBucket"],
				"Resource":["arn:aws:s3:::` + bucket + `"]
			},
			{
				"Effect":"Allow",
				"Principal":{"AWS":["*"]},
				"Action":["s3:GetObject"],
				"Resource":["arn:aws:s3:::` + bucket + `/*"]
			}
		]
	}`

	err = client.SetBucketPolicy(ctx, bucket, policy)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func main() {
	ctx := context.Background()
	bucket, ok := os.LookupEnv("MINIO_BUCKET")
	if !ok {
		panic("MINIO_BUCKET env variable is required")
	}

	minioClient, err := setUpMinio(ctx, bucket)
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// TODO: Move to common limits
		var maxSize int64 = 5 * 1024 * 1024
		if file.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{"message": "File exceeded limit: " + fmt.Sprintf("%d", maxSize)})
			return
		}

		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer srcFile.Close()

		zipFile, err := ioutil.TempFile("/tmp", "minio-*")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		defer zipFile.Close()
		defer os.Remove(zipFile.Name())

		_, err = io.Copy(zipFile, srcFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		kind, err := filetype.MatchFile(zipFile.Name())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if kind.MIME.Type != "application" || kind.MIME.Subtype != "zip" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Unsupported file type"})
			return
		}

		id := uuid.New().String()
		_, err = minioClient.FPutObject(ctx, bucket, id, zipFile.Name(), minio.PutObjectOptions{ContentType: "application/zip"})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	addr := "0.0.0.0:" + port

	log.WithField("addr", addr).Info("Server has been started")

	router.Use(gin.Recovery())
	router.Run(addr)
}
