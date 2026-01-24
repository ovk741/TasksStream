package memory

import "github.com/ovk741/TasksStream/internal/domain"

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
