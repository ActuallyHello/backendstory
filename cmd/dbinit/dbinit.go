package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/config"
	"github.com/ActuallyHello/backendstory/pkg/container"
	"github.com/ActuallyHello/backendstory/pkg/core"
)

const (
	dbInitCommand = "DB_INIT_COMMAND"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	slog.Info("Loading dbinit application...")

	config := config.MustLoadConfig(".")

	slog.Info("Config was loaded...")

	container, err := container.NewAppContainer(ctx, config)
	if err != nil {
		slog.Error("Failed to init container...")
		os.Exit(1)
	}
	defer container.Close()

	slog.Info("Dependency container uploaded! DBinit ready to start!")

	err = setupEnumValues(ctx, container)
	if err != nil {
		slog.Error("Setup values error", "error", err)
	}
	stop()

	<-ctx.Done()
	slog.Info("Application stopped gracefully")
}

func setupEnumValues(ctx context.Context, container *container.AppContainer) error {
	enumService := container.GetEnumService()
	enumValueServicce := container.GetEnumValueService()

	// Товары
	productStatus := enum.Enum{
		Code:  "ProductStatus",
		Label: "Статус товара",
	}
	product, err := enumService.Create(ctx, productStatus)
	if err != nil {
		return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса товара")
	}
	for _, enumValue := range []enumvalue.EnumValue{
		{
			Code:   "Available",
			Label:  "В наличии",
			EnumID: product.ID,
		},
		{
			Code:   "Unavailable",
			Label:  "Нет в наличии",
			EnumID: product.ID,
		},
	} {
		_, err := enumValueServicce.Create(ctx, enumValue)
		if err != nil {
			return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса товара")
		}
	}

	// Заказы
	orderStatus := enum.Enum{
		Code:  "OrderStatus",
		Label: "Статус заказа",
	}
	order, err := enumService.Create(ctx, orderStatus)
	if err != nil {
		return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса заказа")
	}
	for _, enumValue := range []enumvalue.EnumValue{
		{
			Code:   "InProgress",
			Label:  "В обработке",
			EnumID: order.ID,
		},
		{
			Code:   "Approved",
			Label:  "Подтвержден",
			EnumID: order.ID,
		},
		{
			Code:   "Cancelled",
			Label:  "Отменен",
			EnumID: order.ID,
		},
	} {
		_, err := enumValueServicce.Create(ctx, enumValue)
		if err != nil {
			return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса заказа")
		}
	}

	// Элементы заказа
	orderItemStatus := enum.Enum{
		Code:  "OrderItemStatus",
		Label: "Статус элемента заказа",
	}
	orderItem, err := enumService.Create(ctx, orderItemStatus)
	if err != nil {
		return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса элемента заказа")
	}
	for _, enumValue := range []enumvalue.EnumValue{
		{
			Code:   "InProgress",
			Label:  "В обработке",
			EnumID: orderItem.ID,
		},
		{
			Code:   "Approved",
			Label:  "Подтвержден",
			EnumID: orderItem.ID,
		},
		{
			Code:   "Cancelled",
			Label:  "Отменен",
			EnumID: orderItem.ID,
		},
	} {
		_, err := enumValueServicce.Create(ctx, enumValue)
		if err != nil {
			return core.NewTechnicalError(err, dbInitCommand, "Ошибка при установки статуса элемента заказа")
		}
	}
	return nil
}
