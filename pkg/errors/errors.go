package errors

import (
	"errors"
	"github.com/gin-gonic/gin"
	"good-wedding/pkg/model"
	"net/http"
)

const (
	UnknownError          = "UnknownError"
	VnUnknownErrorMessage = "Có lỗi xảy ra"

	NotFound   = "NotFound"
	VNNotFound = "Không tìm thấy dữ liệu"

	ValidationError          = "ValidationError"
	VnValidationErrorMessage = "Dữ liệu không hợp lệ"

	NotAuthenticated          = "NotAuthenticated"
	VnNotAuthenticatedMessage = "Không xác thực"

	UnAuthorized          = "UnAuthorized"
	VnUnAuthorizedMessage = "Không có quyền"

	BadRequest   = "BadRequest"
	VNBadRequest = "Yêu cầu không hợp lệ"

	InternalServerError   = "InternalServerError"
	VNInternalServerError = "Lỗi server"

	PermissionDenied   = "PermissionDenied"
	VnPermissionDenied = "Không có quyền truy cập"

	AccountNotFound          = "AccountNotFound"
	VnAccountNotFoundMessage = "Không tìm thấy tài khoản"

	MissingRequiredFields          = "MissingRequiredFields"
	VnMissingRequiredFieldsMessage = "Thiếu dữ liệu nhập vào "

	MissingAuthorizationHeader   = "MissingAuthorizationHeader"
	VnMissingAuthorizationHeader = "Thiếu thông tin xác thực"

	InvalidAuthorizationFormat   = "InvalidAuthorizationFormat"
	VnInvalidAuthorizationFormat = "Định dạng Authorization không hợp lệ"

	TokenInvalid   = "TokenInvalid"
	VnTokenInvalid = "Token hết hạn hoặc không hợp lệ"

	MissingObjectID   = "MissingObjectID"
	VNMissingObjectID = "Thiếu link ảnh"

	MissingUsername   = "MissingUsername"
	VNMissingUsername = "Tên không được để trống"

	MissingComment   = "MissingComment"
	VNMissingComment = "Bình luận không được để trống"

	MissingWish   = "MissingWish"
	VNMissingWish = "Lời chúc không được để trống"
)

func FeAppError(err string, errType string) *ResponseError {
	return &ResponseError{
		ErrorResp: ErrorResponse{
			Code:    errType,
			Message: err,
		},
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ResponseError struct {
	ErrorResp ErrorResponse `json:"error"`
}

func (er *ResponseError) Error() string {
	return er.ErrorResp.Message
}

type MetaResponse struct {
	TraceID string `json:"traceId"`
	Success bool   `json:"success"`
}

type MessagesResponse struct {
	Meta MetaResponse  `json:"meta"`
	Err  ErrorResponse `json:"error"`
}

func ErrorHandlerMiddleware(c *gin.Context) {
	c.Next()
	errs := c.Errors

	if len(errs) > 0 {
		var err *ResponseError
		ok := errors.As(errs[0].Err, &err)
		if ok {
			meta := model.NewMetaDataWithTraceID(c.Request.Context())

			resp := MessagesResponse{
				Meta: MetaResponse{
					TraceID: meta.TraceID,
				},
				Err: ErrorResponse{
					Code:    err.ErrorResp.Code,
					Message: err.ErrorResp.Message,
				},
			}

			switch err.ErrorResp.Code {
			case NotFound:
				c.JSON(http.StatusNotFound, resp)
				return
			case ValidationError:
				c.JSON(http.StatusBadRequest, resp)
				return
			case NotAuthenticated:
				c.JSON(http.StatusUnauthorized, resp)
				return
			case UnAuthorized:
				c.JSON(http.StatusForbidden, resp)
				return
			case AccountNotFound:
				c.JSON(http.StatusNotFound, resp)
				return
			case MissingRequiredFields:
				c.JSON(http.StatusBadRequest, resp)
				return
			case BadRequest:
				c.JSON(http.StatusBadRequest, resp)
				return
			case MissingAuthorizationHeader:
				c.JSON(http.StatusBadRequest, resp)
				return
			case InvalidAuthorizationFormat:
				c.JSON(http.StatusBadRequest, resp)
				return
			case TokenInvalid:
				c.JSON(http.StatusBadRequest, resp)
				return
			case MissingObjectID:
				c.JSON(http.StatusBadRequest, resp)
				return
			case MissingUsername:
				c.JSON(http.StatusBadRequest, resp)
				return
			case MissingComment:
				c.JSON(http.StatusBadRequest, resp)
				return
			case MissingWish:
				c.JSON(http.StatusBadRequest, resp)
				return
			case InternalServerError:
				c.JSON(http.StatusInternalServerError, resp)
				return
			default:
				c.JSON(http.StatusInternalServerError, resp)
				return
			}
		}
		return
	}
}
