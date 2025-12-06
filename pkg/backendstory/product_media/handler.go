package productmedia

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/product"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/resources"
	"github.com/ActuallyHello/backendstory/pkg/core"
	"github.com/go-playground/validator/v10"
)

const (
	productMediaHandlerCode = "PRODUCT_MEDIA_HANDLER"
)

// ProductMediaHandler provides handlers for product media operations
// @Description Handler for managing product images and media files
type ProductMediaHandler struct {
	validate            *validator.Validate
	productMediaService ProductMediaService
	productService      product.ProductService
	fileService         resources.FileService
	staticFilesPath     string
}

// NewProductMediaHandler creates a new ProductMediaHandler instance
// @Summary Create product media handler
// @Description Initializes a new product media handler with required dependencies
func NewProductMediaHandler(
	productMediaService ProductMediaService,
	productService product.ProductService,
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
// @Success 201 {object} ProductMediaDTO "Изображение успешно загружено"
// @Failure 400 {object} core.ErrorResponse "Неверный запрос: отсутствует product_id или файл, неверный формат данных"
// @Failure 401 {object} core.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} core.ErrorResponse "Товар не найден"
// @Failure 413 {object} core.ErrorResponse "Превышен максимальный размер файла (10MB)"
// @Failure 415 {object} core.ErrorResponse "Неподдерживаемый тип файла"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /product-media/upload [post]
// @ID uploadProductMediaImage
func (h *ProductMediaHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Парсим multipart форму
	if err := r.ParseMultipartForm(resources.MaxFileSize); err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productMediaHandlerCode, "Failed to parse form"))
		return
	}

	// Получаем product_id из формы
	productIDStr := r.FormValue("product_id")
	if productIDStr == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productMediaHandlerCode, "product_id is required"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		core.HandleError(w, r, core.NewLogicalError(err, productMediaHandlerCode, "Invalid product_id"))
		return
	}

	// Проверяем существование товара
	_, err = h.productService.GetByID(ctx, uint(productID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	// Получаем файл из формы
	file, header, err := r.FormFile("file")
	if err != nil {
		core.HandleError(w, r, core.NewLogicalError(err, productMediaHandlerCode, "Failed to get file from form"))
		return
	}
	defer file.Close()

	filePath, err := h.fileService.CreateFromWeb(file, header, h.staticFilesPath)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	productMedia := ProductMedia{
		Link:      filePath,
		ProductID: uint(productID),
	}
	createdMedia, err := h.productMediaService.Create(ctx, productMedia)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ToProductMediaDTO(createdMedia))
}

// GetByProductID возвращает все изображения для товара
// @Summary Получить изображения товара
// @Description Возвращает все изображения для указанного товара. Для доступа требуется аутентификация.
// @Tags Product Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param product_id path integer true "ID товара" minimum(1)
// @Success 200 {array} ProductMediaDTO "Список изображений товара"
// @Failure 400 {object} core.ErrorResponse "Неверный ID товара"
// @Failure 401 {object} core.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} core.ErrorResponse "Товар не найден"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /product-media/product/{product_id} [get]
// @ID getProductMediaByProductId
func (h *ProductMediaHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	productIDStr := r.PathValue("product_id")
	if productIDStr == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productMediaHandlerCode, "product_id parameter missing"))
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		core.HandleError(w, r, core.NewLogicalError(err, productMediaHandlerCode, "Invalid product_id"))
		return
	}

	// Проверяем существование товара
	_, err = h.productService.GetByID(ctx, uint(productID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	mediaList, err := h.productMediaService.GetByProductID(ctx, uint(productID))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	dtos := make([]ProductMediaDTO, len(mediaList))
	for i, media := range mediaList {
		dtos[i] = ToProductMediaDTO(media)
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
// @Failure 400 {object} core.ErrorResponse "Неверный ID изображения"
// @Failure 401 {object} core.ErrorResponse "Требуется аутентификация"
// @Failure 403 {object} core.ErrorResponse "Доступ запрещен: недостаточно прав"
// @Failure 404 {object} core.ErrorResponse "Изображение не найдено"
// @Failure 500 {object} core.ErrorResponse "Внутренняя ошибка сервера"
// @Router /product-media/{id} [delete]
// @ID deleteProductMedia
func (h *ProductMediaHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	if idStr == "" {
		core.HandleError(w, r, core.NewLogicalError(nil, productMediaHandlerCode, "id parameter missing"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		core.HandleError(w, r, core.NewLogicalError(err, productMediaHandlerCode, "Invalid id"))
		return
	}

	// Получаем информацию об изображении
	media, err := h.productMediaService.GetByID(ctx, uint(id))
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	// Удаляем файл с диска
	if media.Link != "" {
		h.fileService.Delete(media.Link, h.staticFilesPath)
	}

	// Удаляем запись из БД
	err = h.productMediaService.Delete(ctx, media)
	if err != nil {
		core.HandleError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
