package repository

import (
	"database/sql"
	"fmt"
	"github.com/shadkain/db_hw/internal/vars"
)

func Error(err error) error {
	if err == sql.ErrNoRows {
		return vars.ErrNotFound
	}

	return err
}

func (this *repositoryImpl) getOrder(desc bool) string {
	if desc {
		return " desc"
	}

	return ""
}
func (this *repositoryImpl) getLimit(limit int) string {
	if limit > 0 {
		return fmt.Sprintf(" limit %d", limit)
	}

	return ""
}

func (this *repositoryImpl) getSinceOperator(desc bool) string {
	if desc {
		return "<"
	}

	return ">"
}
