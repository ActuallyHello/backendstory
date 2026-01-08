package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/ActuallyHello/backendstory/pkg/backendstory/enum"
	"github.com/ActuallyHello/backendstory/pkg/backendstory/enumvalue"
	"github.com/ActuallyHello/backendstory/pkg/config"
	"github.com/ActuallyHello/backendstory/pkg/container"
)

func main() {
	slog.Info("Loading dbinit application...")

	config := config.MustLoadConfig(".")

	slog.Info("Config was loaded...")

	container := container.NewAppContainer(config)

	slog.Info("Dependency container uploaded! DBinit ready to start!")

	ctx := context.Background()

	setupEnumValues(ctx, container)
}

func setupEnumValues(ctx context.Context, container *container.AppContainer) {
	enumService := container.GetEnumService()
	enumValueServicce := container.GetEnumValueService()

	// Товары
	productStatus := enum.Enum{
		Code:  "ProductStatus",
		Label: "Статус товара",
	}
	product, err := enumService.Create(ctx, productStatus)
	if err != nil {
		log.Fatalf("Ошибка при установки статуса товара: %v", err)
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
			log.Fatalf("Ошибка при установки статуса товара: %v", err)
		}
	}

	// Заказы
	orderStatus := enum.Enum{
		Code:  "OrderStatus",
		Label: "Статус заказа",
	}
	order, err := enumService.Create(ctx, orderStatus)
	if err != nil {
		log.Fatalf("Ошибка при установки статуса заказа: %v", err)
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
			log.Fatalf("Ошибка при установки статуса заказа: %v", err)
		}
	}

	// Элементы заказа
	orderItemStatus := enum.Enum{
		Code:  "OrderItemStatus",
		Label: "Статус элемента заказа",
	}
	orderItem, err := enumService.Create(ctx, orderItemStatus)
	if err != nil {
		log.Fatalf("Ошибка при установки статуса элемента заказа: %v", err)
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
			log.Fatalf("Ошибка при установки статуса элемента заказа: %v", err)
		}
	}
}
