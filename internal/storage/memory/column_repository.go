package memory

import (
	"github.com/ovk741/TasksStream/internal/domain"
)

type ColumnRepository struct {
	columns map[string]domain.Column
}

func NewColumnRepository() *ColumnRepository {
	return &ColumnRepository{
		columns: make(map[string]domain.Column),
	}
}

func (r *ColumnRepository) Create(column domain.Column) {
	r.columns[column.ID] = column
}

func (r *ColumnRepository) GetByBoardID(boardID string) []domain.Column {
	result := []domain.Column{}

	for _, column := range r.columns {
		if column.BoardID == boardID {
			result = append(result, column)
		}
	}

	return result
}

func (r *ColumnRepository) Update(column domain.Column) error {
	if _, ok := r.columns[column.ID]; !ok {
		return ErrNotFound
	}
	r.columns[column.ID] = column
	return nil
}

func (r *ColumnRepository) Delete(columnID string) error {
	if _, ok := r.columns[columnID]; !ok {
		return ErrNotFound
	}

	delete(r.columns, columnID)
	return nil
}

func (r *ColumnRepository) GetByID(columnID string) (domain.Column, error) {
	column, ok := r.columns[columnID]
	if !ok {
		return domain.Column{}, ErrNotFound
	}
	return column, nil
}
