package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mytoolzone/task-mini-program/internal/app_code"
	"github.com/mytoolzone/task-mini-program/internal/controller/http/http_util"
	"github.com/mytoolzone/task-mini-program/internal/usecase"
)

type fileRoutes struct {
	fc usecase.File
}

func newFileRoutes(handler *gin.RouterGroup, auth gin.HandlerFunc, role gin.HandlerFunc, f usecase.File) fileRoutes {
	tr := fileRoutes{f}
	h := handler.Group("/file")
	{
		h.POST("/upload", tr.uploadFile)
	}

	return tr
}

type doUploadFileResponse struct {
	Url string `json:"url"`
}

// @Summary 上传文件
// @Description 上传文件
// @Tags 文件
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} http_util.Response{data=doUploadFileResponse}
// @Failure 400 {object} http_util.Response
// @Failure 500 {object} http_util.Response
// @Router /file/upload [post]
func (r fileRoutes) uploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		http_util.Error(context, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	// 检测文件类型是否合法,只允许上传图片
	if !http_util.IsImage(file) {
		http_util.Error(context, app_code.WithError(app_code.ErrorNotImage, errors.New("只允许上传图片")))
		return
	}

	fileEntity, err := r.fc.UploadFile(context, file)
	if err != nil {
		http_util.Error(context, app_code.WithError(app_code.ErrorBadRequest, err))
		return
	}

	http_util.Success(context, doUploadFileResponse{
		Url: fileEntity.Path,
	})
}
