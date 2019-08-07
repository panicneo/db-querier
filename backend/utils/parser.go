package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	ErrInternal = func(err error) error {
		return errors.New("数据库异常：" + err.Error())
	}
)

type RowsFormatter interface {
	Format(rows *sqlx.Rows) (interface{}, error)
}

var TableLikeFormatter RowsFormatter = new(tableLikeFormat)
var SliceMapFormatter RowsFormatter = new(sliceMapFormat)

type sliceMapFormat struct{}

func (sliceMapFormat) Format(rows *sqlx.Rows) (interface{}, error) {
	defer func() {
		_ = rows.Close()
	}()
	result := make([]gin.H, 0)
	for rows.Next() {
		entry := make(gin.H)
		if err := rows.MapScan(entry); err != nil {
			return nil, ErrInternal(err)
		}
		for k, encoded := range entry {
			switch encoded.(type) {
			case []byte:
				entry[k] = string(encoded.([]byte))
			}
		}
		result = append(result, entry)
	}

	return result, nil
}

type tableLikeFormat struct{}

func (tableLikeFormat) Format(rows *sqlx.Rows) (interface{}, error) {
	defer func() {
		_ = rows.Close()
	}()
	headers, _ := rows.Columns()
	result := make([][]interface{}, 0)
	for rows.Next() {
		entry, err := rows.SliceScan()
		if err != nil {
			return nil, ErrInternal(err)
		}
		for k, encoded := range entry {
			switch encoded.(type) {
			case []byte:
				entry[k] = string(encoded.([]byte))
			}
		}
		result = append(result, entry)
	}
	return gin.H{"headers": headers, "rows": result}, nil
}
