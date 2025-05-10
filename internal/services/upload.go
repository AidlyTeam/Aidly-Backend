package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	serviceErrors "github.com/AidlyTeam/Aidly-Backend/internal/errors"
	"github.com/AidlyTeam/Aidly-Backend/pkg/file"
	"github.com/google/uuid"
)

type uploadService struct {
	utilService IUtilService
	uploadDir   string
}

func newUploadService(
	utilService IUtilService,
) *uploadService {
	return &uploadService{
		utilService: utilService,
		uploadDir:   "uploads",
	}
}

func (s *uploadService) SaveImage(file *multipart.FileHeader, filePath string) error {
	if err := s.createDirectories(); err != nil {
		return err
	}

	// Dosya uzantısını kontrol ediyoruz
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return serviceErrors.NewServiceErrorWithMessage(400, serviceErrors.ErrInvalidFileType)
	}

	// Dosyayı açıyoruz
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Kaydetmek için hedef dosyayı oluşturuyoruz
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Dosya içeriğini kopyalıyoruz
	_, err = io.Copy(outFile, src)
	if err != nil {
		return err
	}

	return nil
}

func (s *uploadService) createDirectories() error {
	if err := file.CheckDir(s.uploadDir); err != nil {
		if err := file.CreateDir(s.uploadDir); err != nil {
			return serviceErrors.NewServiceErrorWithMessageAndError(
				500,
				"CREATE_DIRECTORY_ERROR",
				err,
			)
		}
	}
	return nil
}

func (s *uploadService) UploadDir() string {
	return s.uploadDir
}

func (s *uploadService) DeleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (s *uploadService) CreatePath(fileName string) string {
	ext := filepath.Ext(fileName)

	fmt.Println(ext, "dsa")
	var path string
	if ext != "" {
		path = s.UploadDir() + "/" + uuid.New().String() + ext
	} else {
		path = s.UploadDir() + "/" + uuid.New().String()
	}

	return path
}
