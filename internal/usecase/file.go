package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/gw123/glog"
	"github.com/mytoolzone/task-mini-program/internal/entity"
)

// FileUseCase 文件处理
type FileUseCase struct {
	RootDir string
}

func NewFileUseCase(rootDir string) *FileUseCase {
	if rootDir == "" {
		rootDir = "./public/upload"
	}

	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rootDir, os.ModePerm); err != nil {
			glog.WithErr(err).Errorf("create rootDir")
		}
	}

	return &FileUseCase{
		RootDir: rootDir,
	}
}

// UploadFile 上传文件
func (f FileUseCase) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader) (*entity.File, error) {
	// 保存 file 到本地
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dstFullPath, err := f.writeFile(f.createFilename(fileHeader.Filename), file)
	if err != nil {
		return nil, err
	}

	fileEntity := &entity.File{
		Name:      fileHeader.Filename,
		Path:      dstFullPath,
		Size:      fileHeader.Size,
		CreatedAt: time.Now(),
	}
	return fileEntity, nil
}

func (f FileUseCase) writeFile(filename string, file multipart.File) (string, error) {
	dstFullPath := f.getFileFullPath()
	// 判断目录是否存在
	if _, err := os.Stat(dstFullPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dstFullPath, os.ModePerm); err != nil {
			return "", err
		}
	}

	dstFullFilePath := dstFullPath + "/" + filename
	newFile, err := os.OpenFile(dstFullFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(newFile, file)
	if err != nil {
		return "", err
	}
	return dstFullFilePath, nil
}

func (f FileUseCase) createFilename(tmpName string) string {
	return time.Now().Format("150405") + "_" + tmpName
}

func (f FileUseCase) getFileFullPath() string {
	return fmt.Sprintf("%s/%s", f.RootDir, time.Now().Format("20060102"))
}
