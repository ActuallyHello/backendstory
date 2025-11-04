package common

// import (
// 	"context"
// 	"log/slog"

// 	"github.com/ActuallyHello/backendstory/internal/store/entities"
// )

// // FindByCriteria универсальный метод поиска по критериям
// func (r *enumerationValueRepository) FindByCriteria(ctx context.Context, criteria SearchCriteria) (PaginatedResult[entities.EnumerationValue], error) {
// 	var enumerationValues []entities.EnumerationValue
// 	var totalCount int64

// 	// Создаем базовый запрос
// 	query := r.db.WithContext(ctx).Model(&entities.EnumerationValue{})

// 	// Если нужно общее количество, считаем до применения пагинации
// 	if criteria.WithTotal {
// 		countQuery := r.db.WithContext(ctx).Model(&entities.EnumerationValue{})
// 		criteria.Apply(countQuery) // применяем условия, но не пагинацию
// 		if err := countQuery.Count(&totalCount).Error; err != nil {
// 			slog.Error("Count failed in FindByCriteria", "error", err)
// 			return PaginatedResult[entities.EnumerationValue]{}, err
// 		}
// 	}

// 	// Применяем все критерии к основному запросу
// 	query = criteria.Apply(query)

// 	// Выполняем запрос
// 	if err := query.Find(&enumerationValues).Error; err != nil {
// 		slog.Error("FindByCriteria failed", "error", err, "criteria", criteria)
// 		return PaginatedResult[entities.EnumerationValue]{}, err
// 	}

// 	result := PaginatedResult[entities.EnumerationValue]{
// 		Data:       enumerationValues,
// 		TotalCount: totalCount,
// 		Limit:      criteria.Limit,
// 		Offset:     criteria.Offset,
// 		HasMore:    criteria.Limit > 0 && len(enumerationValues) == criteria.Limit,
// 	}

// 	slog.Info("Found records by criteria", "count", len(enumerationValues), "total", totalCount, "criteria", criteria)
// 	return result, nil
// }

// FindByNames поиск по имени и фамилии с LIKE
// func (r *personRepository) FindByNames(ctx context.Context, firstname, lastname string, limit, offset int) (PaginatedResult[entities.Person], error) {
// 	criteria := SearchCriteria{}.
// 		AddOrder("CREATEDAT", OrderDESC).
// 		SetPagination(limit, offset).
// 		SetWithTotal(true)

// 	if firstname != "" {
// 		criteria.AddCondition("FIRSTNAME", OpLike, firstname)
// 	}
// 	if lastname != "" {
// 		criteria.AddCondition("LASTNAME", OpLike, lastname)
// 	}

// 	return r.FindByCriteria(ctx, criteria)
// }

// // FindByUserID поиск по user ID
// func (r *personRepository) FindByUserID(ctx context.Context, userID uint) (entities.Person, error) {
// 	var person entities.Person

// 	criteria := SearchCriteria{}.
// 		AddCondition("USERID", OpEqual, userID)

// 	query := r.db.WithContext(ctx).Model(&entities.Person{})
// 	query = criteria.Apply(query)

// 	if err := query.First(&person).Error; err != nil {
// 		slog.Error("FindByUserID failed", "error", err, "user_id", userID)
// 		return entities.Person{}, err
// 	}

// // 	return person, nil
// // }

// // Пример 1: Простой поиск с пагинацией
// 	result1, err := personRepo.FindByNames(ctx, "John", "Doe", 10, 0)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Пример 2: Сложный поиск с использованием универсального метода
// 	criteria := SearchCriteria{}.
// 		AddCondition("FIRSTNAME", OpLike, "Joh").
// 		AddCondition("LASTNAME", OpLike, "Do").
// 		AddCondition("PHONE", OpLike, "123").
// 		AddOrder("FIRSTNAME", OrderASC).
// 		AddOrder("LASTNAME", OrderASC).
// 		SetPagination(20, 0).
// 		SetWithTotal(true)

// 	result2, err := personRepo.FindByCriteria(ctx, criteria)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Пример 3: Поиск с IN оператором
// 	statusCriteria := SearchCriteria{}.
// 		AddCondition("STATUSID", OpIn, []interface{}{1, 2, 3}).
// 		AddOrder("CREATEDAT", OrderDESC).
// 		SetPagination(50, 0).
// 		SetWithTotal(true)

// 	result3, err := userRepo.FindByCriteria(ctx, statusCriteria)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("Found %d users out of %d total", len(result3.Data), result3.TotalCount)
