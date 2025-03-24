package repository

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/wiqwi12/effective-mobile-test/internal/domain/models"
	"github.com/wiqwi12/effective-mobile-test/pkg/logger"
)

type GroupRepository struct {
	db     *sql.DB
	Logger *logger.Logger
}

func NewGroupRepository(db *sql.DB, logger *logger.Logger) *GroupRepository {
	return &GroupRepository{
		db:     db,
		Logger: logger,
	}
}

func (r *GroupRepository) CreateGroup(group models.Group) error {
	query, args, err := squirrel.Insert("groups").Columns("name, id").
		Values(group.Name, group.Id).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for group creation",
			"error", err,
			"group_name", group.Name,
			"group_id", group.Id)
		return err
	}

	_, err = r.db.Exec(query, args...)
	if err != nil {
		r.Logger.Info.Error("Failed to execute group creation query",
			"error", err,
			"group_name", group.Name,
			"group_id", group.Id)
		return err
	}

	return nil
}

func (r *GroupRepository) GetGroupByName(name string) (models.Group, error) {
	query, args, err := squirrel.Select("id, name").
		From("groups").
		Where(squirrel.Eq{
			"name": name,
		}).PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		r.Logger.Info.Error("Failed to build SQL query for group lookup",
			"error", err,
			"group_name", name)
		return models.Group{}, err
	}

	var group models.Group
	err = r.db.QueryRow(query, args...).Scan(&group.Id, &group.Name)
	if err != nil {
		if err == sql.ErrNoRows {
		} else {
			r.Logger.Info.Error("Error executing group lookup query",
				"error", err,
				"group_name", name)
		}
		return models.Group{}, err
	}

	return group, nil
}

func (r *GroupRepository) GroupExsist(ctx context.Context, name string) (bool, error) {
	query, args, err := squirrel.Select("1").
		From("groups").
		Where(squirrel.Eq{"name": name}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		r.Logger.Info.Error("Failed to build query for group existence check",
			"error", err,
			"group_name", name)
		return false, err
	}

	var exists bool
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		r.Logger.Info.Error("Error checking if group exists",
			"error", err,
			"group_name", name)
		return false, err
	}

	return true, nil
}
