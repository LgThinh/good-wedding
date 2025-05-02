package repo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"good-wedding/conf"
	"good-wedding/pkg/model"
	"good-wedding/pkg/utils"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
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
	UploadImageToS3(file *multipart.FileHeader) (*model.UploadImageSuccessResponse, error)
	UploadVideoToS3(file *multipart.FileHeader) (*model.UploadVideoSuccessResponse, error)
}

func (r *WeddingRepo) UploadImageToS3(file *multipart.FileHeader) (*model.UploadImageSuccessResponse, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open image: %v", err)
	}
	defer fileContent.Close()

	filename := file.Filename
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		return nil, fmt.Errorf("invalid file type, only jpg, jpeg, png are allowed")
	}

	prefix := utils.RandStringBytes(12, false)
	newFileName := "image/" + prefix + "." + filename

	var contentType string
	switch fileExt {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	default:
		return nil, fmt.Errorf("invalid file type, only jpg, jpeg, png are allowed")
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(conf.GetConfig().AWSBucketName),
		Key:         aws.String(newFileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}

	_, err = r.S3Bucket.PutObject(uploadInput)
	if err != nil {
		return nil, fmt.Errorf("unable to upload image: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", conf.GetConfig().AWSBucketName, newFileName)

	err = r.SaveFileToDB(file, fileURL, utils.Image)
	if err != nil {
		return nil, err
	}

	rs := &model.UploadImageSuccessResponse{
		Meta: model.NewMetaData(),
		Data: model.UploadImageUrl{
			Url: fileURL,
		},
	}

	return rs, nil
}

func (r *WeddingRepo) UploadVideoToS3(file *multipart.FileHeader) (*model.UploadVideoSuccessResponse, error) {
	fileContent, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("unable to open video: %v", err)
	}
	defer fileContent.Close()

	filename := file.Filename
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	if fileExt != ".mp4" && fileExt != ".avi" && fileExt != ".mov" && fileExt != ".mkv" && fileExt != ".ts" {
		return nil, fmt.Errorf("invalid file type, only mp4, avi, mov, mkv , ts are allowed")
	}

	prefix := utils.RandStringBytes(8, false)
	newFileName := "video/" + prefix + "." + filename

	var contentType string
	switch fileExt {
	case ".mp4":
		contentType = "video/mp4"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".mkv":
		contentType = "video/x-matroska"
	case ".ts":
		contentType = "video/MP2T"

	default:
		return nil, fmt.Errorf("invalid file type, only mp4, avi, mov, mkv , ts are allowed")
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(conf.GetConfig().AWSBucketName),
		Key:         aws.String(newFileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}

	_, err = r.S3Bucket.PutObject(uploadInput)
	if err != nil {
		return nil, fmt.Errorf("unable to upload video: %v", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", conf.GetConfig().AWSBucketName, newFileName)

	err = r.SaveFileToDB(file, fileURL, utils.Video)
	if err != nil {
		return nil, err
	}

	rs := &model.UploadVideoSuccessResponse{
		Meta: model.NewMetaData(),
		Data: model.UploadVideoUrl{
			Url: fileURL,
		},
	}

	return rs, nil
}

func (r *WeddingRepo) SaveFileToDB(file *multipart.FileHeader, url, fileType string) error {
	fileContent, err := file.Open()
	if err != nil {
		return fmt.Errorf("unable to open file: %v", err)
	}
	defer fileContent.Close()

	newFile := model.Object{
		BaseModel:  model.BaseModel{},
		Name:       file.Filename,
		ObjectType: fileType,
		Url:        url,
		Size:       file.Size,
	}

	tx := r.DB.Create(&newFile)
	if tx.Error != nil {
		return fmt.Errorf("fail to save file info: %v", err)
	}

	return nil
}
