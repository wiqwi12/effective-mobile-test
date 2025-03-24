package repository

import "github.com/wiqwi12/effective-mobile-test/internal/domain/models"

type GroupRepository interface {
	GetGroupByName(name string) (models.Group, error) //todo заменить на Реквест с фильтрами
	CreateGroup(group models.Group) error
}
