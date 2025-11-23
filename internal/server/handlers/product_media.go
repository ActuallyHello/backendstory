// internal/server/handlers/product_media_handler.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	appErr "github.com/ActuallyHello/backendstory/internal/core/errors"
	"github.com/ActuallyHello/backendstory/internal/dto"
	"github.com/ActuallyHello/backendstory/internal/server/middleware"
	"github.com/ActuallyHello/backendstory/internal/services"
	"github.com/ActuallyHello/backendstory/internal/services/resources"
	"github.com/ActuallyHello/backendstory/internal/store/entities"
	"github.com/go-playground/validator/v10"
)

const (
	productMediaHandlerCode = "PRODUCT_MEDIA_HANDLER"
)

// ProductMediaHandler provides handlers for product media operations
// @Description Handler for managing product images and media files
type ProductMediaHandler struct {
	validate            *validator.Validate
	productMediaService services.ProductMediaService
	productService      services.ProductService
	fileService         resources.FileService
	staticFilesPath     string
}

// NewProductMediaHandler creates a new ProductMediaHandler instance
// @Summary Create product media handler
// @Description Initializes a new product media handler with required dependencies
func NewProductMediaHandler(
	productMediaService services.ProductMediaService,
	productService services.ProductService,
	fileService resources.FileService,
	staticFilesPath string,
) *ProductMediaHandler {
	return &ProductMediaHandler{
		validate:            validator.New(),
		productMediaService: productMediaService,
		productService:      productService,
		// TODO refactor
		fileService:     resources.FileService{},
		staticFilesPath: staticFilesPath,
	}
}

// UploadImage загружает изображение для товара
// @Summary Загрузить изображение товара
// @Description Загружает изображение для указанного товара. Поддерживаемые форматы: JPEG, PNG, GIF. Максимальный размер файла: 10MB.
// @Tags Product Media
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param product_id formData integer true "ID товара" minimum(1)
// @Param file formData file true "Изображение товара"
// @Success 201 {object} dto.ProductMediaDTO "Изображение успешно загружено"
// @Failure 400 {object} dto.ErrorResponse "Неверный запрос: отсутствует product_id или файл, неверный формат данных"
// @Failure 401 {object} dto.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} dto.ErrorResponse "Товар не найден"
// @Failure 413 {object} dto.ErrorResponse "Превышен максимальный размер файла (10MB)"
// @Failure 415 {object} dto.ErrorResponse "Неподдерживаемый тип файла"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/product-media/upload [post]
// @OperationId uploadProductMediaImage
func (h *ProductMediaHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Парсим multipart форму
	if err := r.ParseMultipartForm(resources.MaxFileSize); err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productMediaHandlerCode, "Failed to parse form"))
		return
	}

	// Получаем product_id из формы
	productIDStr := r.FormValue("product_id")
	if productIDStr == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productMediaHandlerCode, "product_id is required"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productMediaHandlerCode, "Invalid product_id"))
		return
	}

	// Проверяем существование товара
	_, err = h.productService.GetByID(ctx, uint(productID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productMediaHandlerCode, "Failed to get file from form"))
		return
	}
	defer file.Close()

	filePath, err := h.fileService.CreateFromWeb(file, header, h.staticFilesPath)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	productMedia := entities.ProductMedia{
		Link:      filePath,
		ProductID: uint(productID),
	}
	createdMedia, err := h.productMediaService.Create(ctx, productMedia)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToProductMediaDTO(createdMedia))
}

// GetByProductID возвращает все изображения для товара
// @Summary Получить изображения товара
// @Description Возвращает все изображения для указанного товара. Для доступа требуется аутентификация.
// @Tags Product Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path integer true "ID товара" minimum(1)
// @Success 200 {array} dto.ProductMediaDTO "Список изображений товара"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID товара"
// @Failure 401 {object} dto.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} dto.ErrorResponse "Товар не найден"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/product-media/product/{product_id} [get]
// @OperationId getProductMediaByProductId
func (h *ProductMediaHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productIDStr := r.PathValue("product_id")
	if productIDStr == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productMediaHandlerCode, "product_id parameter missing"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productMediaHandlerCode, "Invalid product_id"))
		return
	}

	// Проверяем существование товара
	_, err = h.productService.GetByID(ctx, uint(productID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	mediaList, err := h.productMediaService.GetByProductID(ctx, uint(productID))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	dtos := make([]dto.ProductMediaDTO, len(mediaList))
	for i, media := range mediaList {
		dtos[i] = dto.ToProductMediaDTO(media)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dtos)
}

// Delete удаляет изображение товара
// @Summary Удалить изображение товара
// @Description Удаляет изображение товара по ID. Удаляет как запись из БД, так и файл с диска.
// @Tags Product Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "ID изображения" minimum(1)
// @Success 204 "Изображение успешно удалено"
// @Failure 400 {object} dto.ErrorResponse "Неверный ID изображения"
// @Failure 401 {object} dto.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} dto.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} dto.ErrorResponse "Изображение не найдено"
// @Failure 500 {object} dto.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/product-media/{id} [delete]
// @OperationId deleteProductMedia
func (h *ProductMediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		middleware.HandleError(w, r, appErr.NewLogicalError(nil, productMediaHandlerCode, "id parameter missing"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		middleware.HandleError(w, r, appErr.NewLogicalError(err, productMediaHandlerCode, "Invalid id"))
		return
	}

	// Получаем информацию об изображении
	media, err := h.productMediaService.GetByID(ctx, uint(id))
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	// Удаляем файл с диска
	if media.Link != "" {
		h.fileService.Delete(media.Link, h.staticFilesPath)
	}

	// Удаляем запись из БД
	err = h.productMediaService.Delete(ctx, media)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
