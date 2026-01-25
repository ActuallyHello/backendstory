package resources

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	MaxFileSize = 5 << 20 // 5MB

	fileServiceCode = "FILE_SERVICE"
)

type FileService interface {
	CreateImage(file multipart.File, header *multipart.FileHeader, staticFilePath string) (string, error)
	DeleteImage(mediaPath, staticFilePath string) error
}

func NewFileService() FileService {
	return &fileService{}
}

type fileService struct {
}

// returns relative path for saved file
func (fs *fileService) CreateImage(file multipart.File, header *multipart.FileHeader, staticFilesPath string) (string, error) {
	defer file.Close()

	// Проверяем размер файла
	if header.Size > MaxFileSize {
		return "", core.NewLogicalError(nil, fileServiceCode, "File too large")
	}

	// Проверяем тип файла
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", core.NewTechnicalError(err, fileServiceCode, "Failed to read file")
	}

	fileType := http.DetectContentType(buffer)
	if !strings.HasPrefix(fileType, "image/") {
		return "", core.NewLogicalError(nil, fileServiceCode, "Unsupported file type")
	}

	// Возвращаем указатель на начало файла
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", core.NewTechnicalError(err, fileServiceCode, "Failed to reset file pointer")
	}

	// Определяем расширение по типу контента
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		switch fileType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		case "image/gif":
			ext = ".gif"
		default:
			ext = ".bin"
		}
	}

	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("image_%d%s", timestamp, ext)

	// Полный путь к файлу на диске
	fullFilePath := filepath.Join(staticFilesPath, "imgs", filename)

	// Относительный путь для БД (тот, по которому файл будет доступен через HTTP)
	relativePath := filepath.Join("/static/imgs", filename)

	// Создаем директорию если она не существует
	dir := filepath.Dir(fullFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", core.NewTechnicalError(err, fileServiceCode, "Failed to create directory")
	}

	// Создаем файл на диске
	dst, err := os.Create(fullFilePath)
	if err != nil {
		return "", core.NewTechnicalError(err, fileServiceCode, "Failed to create file: "+err.Error())
	}
	defer dst.Close()

	// Копируем содержимое файла
	_, err = io.Copy(dst, file)
	if err != nil {
		// Удаляем файл в случае ошибки
		os.Remove(fullFilePath)
		return "", core.NewTechnicalError(err, fileServiceCode, "Failed to save file: "+err.Error())
	}

	return relativePath, nil
}

func (fs *fileService) DeleteImage(mediaPath, staticFilesPath string) error {
	// Убираем /static из пути, так как у нас уже есть базовый путь
	relativePath := strings.TrimPrefix(mediaPath, "/static")
	filePath := filepath.Join(staticFilesPath, relativePath)

	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			return core.NewTechnicalError(err, fileServiceCode, "Ошибка при удалении")
		}
	}
	return nil
}
