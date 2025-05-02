package repo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"good-wedding/conf"
	"good-wedding/pkg/errors"
	"good-wedding/pkg/model"
	"good-wedding/pkg/utils/logger"
	"gorm.io/gorm"
	"mime/multipart"
)

type WeddingRepo struct {
	DB       *gorm.DB
	S3Bucket *s3.S3
}

func (r *WeddingRepo) DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(ctx, generalQueryTimeout)
	return r.DB.WithContext(ctx), cancel
}

func NewWeddingRepo(imageRepo *gorm.DB, S3Bucket *s3.S3) WeddingRepoInterface {
	return &WeddingRepo{
		DB:       imageRepo,
		S3Bucket: S3Bucket,
	}
}

type WeddingRepoInterface interface {
	DBWithTimeout(ctx context.Context) (*gorm.DB, context.CancelFunc)
	UploadToS3(fileName string, file *s3.PutObjectInput) (*string, error)
	SaveFileToDB(tx *gorm.DB, creatorID uuid.UUID, file *multipart.FileHeader, url, fileType, objectType string) error
}

func (r *WeddingRepo) UploadToS3(fileName string, file *s3.PutObjectInput) (*string, error) {
	log := logger.WithTag("WeddingRepo|UploadToS3")
	_, err := r.S3Bucket.PutObject(file)
	if err != nil {
		logger.LogError(log, err, "unable to upload ")
		appErr := errors.FeAppError("Lỗi khi đăng tải file", errors.UnknownError)
		return nil, appErr
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", conf.GetConfig().AWSBucketName, fileName)

	return &fileURL, nil
}

func (r *WeddingRepo) SaveFileToDB(tx *gorm.DB, creatorID uuid.UUID, file *multipart.FileHeader, url, fileType, objectType string) error {
	log := logger.WithTag("WeddingRepo|SaveFileToDB")

	fileContent, err := file.Open()
	if err != nil {
		logger.LogError(log, err, "unable to open file")
		appErr := errors.FeAppError("File lỗi", errors.UnknownError)
		return appErr
	}
	defer fileContent.Close()

	newFile := model.ObjectMedia{
		BaseModel: model.BaseModel{
			CreatorID: &creatorID,
		},
		ObjectType: objectType,
		Name:       file.Filename,
		FileType:   fileType,
		Url:        url,
		Size:       file.Size,
	}

	tx = r.DB.Create(&newFile)
	if tx.Error != nil {
		logger.LogError(log, err, "fail to save file info")
		appErr := errors.FeAppError("Lỗi khi lưu file", errors.UnknownError)
		return appErr
	}

	return nil
}
