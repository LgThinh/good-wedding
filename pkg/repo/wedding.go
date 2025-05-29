package repo

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"good-wedding/conf"
	"good-wedding/pkg/errors"
	"good-wedding/pkg/model"
	"good-wedding/pkg/utils"
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
	SaveFileToDB(tx *gorm.DB, creatorID uuid.UUID, file *multipart.FileHeader, url, fileType, objectType, customName string) error
	CreateUser(tx *gorm.DB, ob *model.User) (*model.User, error)
	CreateComment(tx *gorm.DB, ob *model.Comment) (*model.Comment, error)
	CreateWish(tx *gorm.DB, ob *model.WeddingWish) (*model.WeddingWish, error)
	CommentFilter(tx *gorm.DB, f *model.CommentFilter) (*model.CommentFilterResult, error)
	WeddingWishFilter(tx *gorm.DB, f *model.WeddingWishFilter) (*model.WeddingWishFilterResult, error)
	GetObjectMedia(tx *gorm.DB, id uuid.UUID) (*model.ObjectMedia, error)
	UserFilter(tx *gorm.DB, f *model.UserFilter) (*model.UserFilterResult, error)
	ObjectMediaFilter(tx *gorm.DB, f *model.ObjectMediaFilter) (*model.ObjectMediaFilterResult, error)
	GetOneImage(tx *gorm.DB, customName string) (*model.GetOneImageResult, error)
	GetOneVideo(tx *gorm.DB, customName string) (*model.GetOneVideoResult, error)
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

func (r *WeddingRepo) SaveFileToDB(tx *gorm.DB, creatorID uuid.UUID, file *multipart.FileHeader, url, fileType, objectType, customName string) error {
	log := logger.WithTag("WeddingRepo|SaveFileToDB")

	fileContent, err := file.Open()
	if err != nil {
		logger.LogError(log, err, "unable to open file")
		return errors.FeAppError("File lỗi", errors.UnknownError)
	}
	defer fileContent.Close()

	newFile := model.ObjectMedia{
		BaseModel: model.BaseModel{
			CreatorID: &creatorID,
		},
		Name:       file.Filename,
		CustomName: customName,
		ObjectType: objectType,
		FileType:   fileType,
		Url:        url,
		Size:       file.Size,
	}

	if err := tx.Create(&newFile).Error; err != nil {
		logger.LogError(log, err, "fail to save file info")
		return errors.FeAppError("Lỗi khi lưu file", errors.UnknownError)
	}

	return nil
}

func (r *WeddingRepo) CreateUser(tx *gorm.DB, ob *model.User) (*model.User, error) {
	log := logger.WithTag("WeddingRepo|CreateUser")
	err := tx.Create(&ob).Error
	if err != nil {
		logger.LogError(log, err, "Error when get create")
		appErr := errors.FeAppError("Không tạo được", errors.UnknownError)
		return nil, appErr
	}

	return ob, nil
}

func (r *WeddingRepo) CreateComment(tx *gorm.DB, ob *model.Comment) (*model.Comment, error) {
	log := logger.WithTag("WeddingRepo|CreateComment")
	err := tx.Create(&ob).Error
	if err != nil {
		logger.LogError(log, err, "Error when get create")
		appErr := errors.FeAppError("Không tạo được", errors.UnknownError)
		return nil, appErr
	}

	return ob, nil
}

func (r *WeddingRepo) CreateWish(tx *gorm.DB, ob *model.WeddingWish) (*model.WeddingWish, error) {
	log := logger.WithTag("WeddingRepo|CreateWish")
	err := tx.Create(&ob).Error
	if err != nil {
		logger.LogError(log, err, "Error when get create")
		appErr := errors.FeAppError("Không tạo được", errors.UnknownError)
		return nil, appErr
	}

	return ob, nil
}

func (r *WeddingRepo) CommentFilter(tx *gorm.DB, f *model.CommentFilter) (*model.CommentFilterResult, error) {
	log := logger.WithTag("WeddingRepo|CommentFilter")

	tx = tx.Model(&model.Comment{})

	if f.FromDate != nil {
		fromDate := utils.ConvertUnixMilliToTime(*f.FromDate)
		tx = tx.Where("created_at >= ?", fromDate)
	}

	if f.ToDate != nil {
		toDate := utils.ConvertUnixMilliToTime(*f.ToDate)
		tx = tx.Where("created_at <= ?", toDate)
	}

	if f.ObjectID != nil {
		tx = tx.Where("object_id = ?", *f.ObjectID)
	}

	var comments []*model.Comment

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}

	pager := f.Pager

	tx = pager.DoQuery(&comments, tx)
	if tx.Error != nil {
		logger.LogError(log, tx.Error, "Error when get list")
		appErr := errors.FeAppError("Không tìm thấy danh sách", errors.NotFound)
		return nil, appErr
	}

	var records []*model.CommentDataResponse

	for _, comment := range comments {
		var (
			user  model.User
			media model.ObjectMedia
		)
		newTx := tx.Session(&gorm.Session{NewDB: true})
		ts := newTx.Table("user").Where("id = ?", comment.UserID).First(&user)
		if ts.Error != nil {
			logger.LogError(log, ts.Error, "Error when getting user")
			err := errors.FeAppError(errors.VNNotFound, errors.NotFound)
			return nil, err
		}

		ts = newTx.Table("object_media").Where("id = ?", comment.ObjectID).First(&media)
		if ts.Error != nil {
			logger.LogError(log, ts.Error, "Error when getting object media")
			err := errors.FeAppError(errors.VNNotFound, errors.NotFound)
			return nil, err
		}

		records = append(records, &model.CommentDataResponse{
			InitTime:  utils.ConvertTimeToMillisString(&comment.CreatedAt),
			ObjectID:  &comment.ObjectID,
			ObjectUrl: media.Url,
			UserID:    user.ID,
			UserName:  user.UserName,
			Comment:   comment.Comment,
		})
	}

	result := &model.CommentFilterResult{
		Filter:  f,
		Records: records,
	}

	if result.Records == nil {
		result.Records = []*model.CommentDataResponse{}
	}

	return result, nil
}

func (r *WeddingRepo) WeddingWishFilter(tx *gorm.DB, f *model.WeddingWishFilter) (*model.WeddingWishFilterResult, error) {
	log := logger.WithTag("WeddingRepo|WeddingWishFilter")

	tx = tx.Model(&model.WeddingWish{})

	if f.FromDate != nil {
		fromDate := utils.ConvertUnixMilliToTime(*f.FromDate)
		tx = tx.Where("created_at >= ?", fromDate)
	}

	if f.ToDate != nil {
		toDate := utils.ConvertUnixMilliToTime(*f.ToDate)
		tx = tx.Where("created_at <= ?", toDate)
	}

	var wishes []*model.WeddingWish

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}

	pager := f.Pager

	tx = pager.DoQuery(&wishes, tx)
	if tx.Error != nil {
		logger.LogError(log, tx.Error, "Error when get list")
		appErr := errors.FeAppError("Không tìm thấy danh sách", errors.NotFound)
		return nil, appErr
	}

	var records []*model.WeddingWishDataResponse

	for _, wish := range wishes {
		var (
			user model.User
		)
		newTx := tx.Session(&gorm.Session{NewDB: true})
		ts := newTx.Table("user").Where("id = ?", wish.UserID).First(&user)
		if ts.Error != nil {
			logger.LogError(log, ts.Error, "Error when getting user")
			err := errors.FeAppError(errors.VNNotFound, errors.NotFound)
			return nil, err
		}

		records = append(records, &model.WeddingWishDataResponse{
			InitTime: utils.ConvertTimeToMillisString(&wish.CreatedAt),
			UserID:   wish.UserID,
			UserName: user.UserName,
			Comment:  wish.Comment,
		})
	}

	result := &model.WeddingWishFilterResult{
		Filter:  f,
		Records: records,
	}

	if result.Records == nil {
		result.Records = []*model.WeddingWishDataResponse{}
	}

	return result, nil
}

func (r *WeddingRepo) GetObjectMedia(tx *gorm.DB, id uuid.UUID) (*model.ObjectMedia, error) {
	log := logger.WithTag("WeddingRepo|GetObjectMedia")
	var todo model.ObjectMedia
	err := tx.Where("id = ?", id).First(&todo).Error
	if err != nil {
		logger.LogError(log, err, "Error when get media")
		appErr := errors.FeAppError("Không tìm thấy ảnh hay video", errors.UnknownError)
		return nil, appErr
	}

	return &todo, nil
}

func (r *WeddingRepo) UserFilter(tx *gorm.DB, f *model.UserFilter) (*model.UserFilterResult, error) {
	log := logger.WithTag("WeddingRepo|UserFilter")

	tx = tx.Model(&model.User{})

	if f.FromDate != nil {
		fromDate := utils.ConvertUnixMilliToTime(*f.FromDate)
		tx = tx.Where("created_at >= ?", fromDate)
	}

	if f.ToDate != nil {
		toDate := utils.ConvertUnixMilliToTime(*f.ToDate)
		tx = tx.Where("created_at <= ?", toDate)
	}

	if f.UserName != nil {
		tx = tx.Where("username = ?", *f.UserName)
	}

	result := &model.UserFilterResult{
		Filter:  f,
		Records: []*model.User{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}

	pager := f.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		logger.LogError(log, tx.Error, "Error when get list")
		appErr := errors.FeAppError("Không tìm thấy danh sách", errors.NotFound)
		return nil, appErr
	}

	if result.Records == nil {
		result.Records = []*model.User{}
	}

	return result, nil
}

func (r *WeddingRepo) ObjectMediaFilter(tx *gorm.DB, f *model.ObjectMediaFilter) (*model.ObjectMediaFilterResult, error) {
	log := logger.WithTag("WeddingRepo|ObjectMediaFilter")

	tx = tx.Model(&model.ObjectMedia{})

	if f.FromDate != nil {
		fromDate := utils.ConvertUnixMilliToTime(*f.FromDate)
		tx = tx.Where("created_at >= ?", fromDate)
	}

	if f.ToDate != nil {
		toDate := utils.ConvertUnixMilliToTime(*f.ToDate)
		tx = tx.Where("created_at <= ?", toDate)
	}

	if f.Name != nil {
		tx = tx.Where("name = ?", *f.Name)
	}

	if f.Url != nil {
		tx = tx.Where("url = ?", *f.Url)
	}

	if f.FileType != nil {
		tx = tx.Where("file_type = ?", *f.FileType)
	}

	if f.ObjectType != nil {
		tx = tx.Where("object_type = ?", *f.ObjectType)
	}

	result := &model.ObjectMediaFilterResult{
		Filter:  f,
		Records: []*model.ObjectMedia{},
	}

	f.Pager.SortableFields = []string{"id", "created_at", "updated_at"}

	pager := f.Pager

	tx = pager.DoQuery(&result.Records, tx)
	if tx.Error != nil {
		logger.LogError(log, tx.Error, "Error when get list")
		appErr := errors.FeAppError("Không tìm thấy danh sách", errors.NotFound)
		return nil, appErr
	}

	if result.Records == nil {
		result.Records = []*model.ObjectMedia{}
	}

	return result, nil
}

func (r *WeddingRepo) GetOneImage(tx *gorm.DB, customName string) (*model.GetOneImageResult, error) {
	log := logger.WithTag("WeddingRepo|GetOneImage")

	var image model.ObjectMedia
	err := tx.Where("custom_name = ?", customName).First(&image).Error
	if err != nil {
		logger.LogError(log, err, "Error when getting image")
		return nil, errors.FeAppError(errors.VNNotFound, errors.NotFound)
	}

	result := model.GetOneImageResult{
		Meta: model.NewMetaData(),
		Data: image,
	}
	return &result, nil
}

func (r *WeddingRepo) GetOneVideo(tx *gorm.DB, customName string) (*model.GetOneVideoResult, error) {
	log := logger.WithTag("WeddingRepo|GetOneVideo")

	var image model.ObjectMedia
	err := tx.Where("custom_name = ?", customName).First(&image).Error
	if err != nil {
		logger.LogError(log, err, "Error when getting video")
		return nil, errors.FeAppError(errors.VNNotFound, errors.NotFound)
	}

	result := model.GetOneVideoResult{
		Meta: model.NewMetaData(),
		Data: image,
	}
	return &result, nil
}
