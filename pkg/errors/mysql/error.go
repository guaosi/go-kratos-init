package mysql

import (
	"github.com/go-kratos/kratos/v2/errors"
	"strings"
)

var (
	ErrMySQLDataDuplicate = errors.New(400, "MYSQL_DATA_DUPLICATE", "mysql data is duplication")
	ErrMySQLDataNotFound  = errors.New(404, "MYSQL_DATA_NOT_FOUND", "mysql data is not found")
	ErrMySQL              = errors.New(500, "MYSQL_ERROR", "mysql error")
)

func JudgeRecordDuplicate(err error) bool {
	if err != nil {
		return strings.Contains(err.Error(), "1062")
	}
	return false
}
