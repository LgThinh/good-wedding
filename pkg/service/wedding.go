package service

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"good-wedding/conf"
	"good-wedding/pkg/errors"
	"good-wedding/pkg/model"
	"good-wedding/pkg/repo"
	"good-wedding/pkg/utils"
	"mime/multipart"
	"strings"
)

type WeddingService struct {
	weddingRepo repo.WeddingRepoInterface
}

func NewWeddingService(weddingRepo repo.WeddingRepoInterface) *WeddingService {
	return &WeddingService{
		weddingRepo: weddingRepo,
	}
}

type WeddingServiceInterface interface {
	UploadImageToS3(ctx context.Context, file *multipart.FileHeader, adminID uuid.UUID) (*model.UploadImageSuccessResponse, error)
	UploadVideoToS3(ctx context.Context, file *multipart.FileHeader, adminID uuid.UUID) (*model.UploadVideoSuccessResponse, error)
}

func (s *WeddingService) UploadImageToS3(ctx context.Context, file *multipart.FileHeader, adminID uuid.UUID) (*model.UploadImageSuccessResponse, error) {
	txWithTimeout, cancel := s.weddingRepo.DBWithTimeout(ctx)
	defer cancel()

	tx := txWithTimeout.Begin()
	defer tx.Rollback()

	fileContent, err := file.Open()
	if err != nil {
		appErr := errors.FeAppError("Lỗi hình ảnh", errors.UnknownError)
		return nil, appErr
	}
	defer fileContent.Close()

	filename := file.Filename
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" {
		appErr := errors.FeAppError("Sai định dạng, chỉ chấp nhận: jpg, jpeg, png", errors.UnknownError)
		return nil, appErr
	}

	prefix := utils.RandStringBytes(12, false)
	publicFileName := utils.Image + "/" + prefix + filename

	var contentType string
	switch fileExt {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	default:
		appErr := errors.FeAppError("Sai định dạng, chỉ chấp nhận:  jpg, jpeg, png", errors.UnknownError)
		return nil, appErr
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(conf.GetConfig().AWSBucketName),
		Key:         aws.String(publicFileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}

	url, err := s.weddingRepo.UploadToS3(publicFileName, uploadInput)
	if err != nil {
		return nil, err
	}

	ext := strings.TrimPrefix(fileExt, ".")
	err = s.weddingRepo.SaveFileToDB(tx, adminID, file, *url, utils.Image, ext)
	if err != nil {
		return nil, err
	}

	result := model.UploadImageSuccessResponse{
		Meta: model.NewMetaData(),
		Data: model.UploadImageUrl{
			Url: *url,
		},
	}

	tx.Commit()
	return &result, nil
}

func (s *WeddingService) UploadVideoToS3(ctx context.Context, file *multipart.FileHeader, adminID uuid.UUID) (*model.UploadVideoSuccessResponse, error) {
	txWithTimeout, cancel := s.weddingRepo.DBWithTimeout(ctx)
	defer cancel()

	tx := txWithTimeout.Begin()
	defer tx.Rollback()

	fileContent, err := file.Open()
	if err != nil {
		appErr := errors.FeAppError("Video lỗi", errors.UnknownError)
		return nil, appErr
	}
	defer fileContent.Close()

	filename := file.Filename
	fileExt := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	if fileExt != ".mp4" && fileExt != ".avi" && fileExt != ".mov" && fileExt != ".mkv" {
		appErr := errors.FeAppError("Sai định dạng, chỉ chấp nhận:  mp4, avi, mov, mkv", errors.UnknownError)
		return nil, appErr
	}

	prefix := utils.RandStringBytes(12, false)
	publicFileName := utils.Video + "/" + prefix + filename

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
	default:
		appErr := errors.FeAppError("Sai định dạng, chỉ chấp nhận: mp4, avi, mov, mkv", errors.UnknownError)
		return nil, appErr
	}

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(conf.GetConfig().AWSBucketName),
		Key:         aws.String(publicFileName),
		Body:        fileContent,
		ContentType: aws.String(contentType),
	}

	url, err := s.weddingRepo.UploadToS3(publicFileName, uploadInput)
	if err != nil {
		return nil, err
	}

	ext := strings.TrimPrefix(fileExt, ".")
	err = s.weddingRepo.SaveFileToDB(tx, adminID, file, *url, utils.Video, ext)
	if err != nil {
		return nil, err
	}

	result := model.UploadVideoSuccessResponse{
		Meta: model.NewMetaData(),
		Data: model.UploadVideoUrl{
			Url: *url,
		},
	}

	tx.Commit()
	return &result, nil
}
