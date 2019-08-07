package query

import (
	"db-querier/service"
	"db-querier/utils"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrParamMissed = func(name string) error {
		return errors.New("missed parameter: " + name)
	}
	ErrParamTypeUnMatch = func(name string, expectType string) error {
		return errors.New(fmt.Sprintf("invalid parameter [%s]'s type, expect a %s", name, expectType))
	}
)

func Query(ctx *gin.Context) {
	name := ctx.Param("name")
	keys := service.Q.QueryParams(name)
	params := make([]interface{}, 0)

	for n, t := range keys.NameTypePairs() {
		q, ok := ctx.GetQuery(n)
		if !ok {
			_ = ctx.Error(ErrParamMissed(n))
			continue
		}
		switch strings.ToLower(t) {
		case "number":
			v, e := cast.ToIntE(q)
			if e != nil {
				_ = ctx.Error(ErrParamTypeUnMatch(n, t))
				continue
			}
			params = append(params, v)
		default:
			params = append(params, q)
		}
	}
	if len(ctx.Errors) != 0 {
		return
	}
	rows, err := service.Q.Query(name, params...)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}

	var formatter utils.RowsFormatter
	if ctx.Query("format") == "table" {
		formatter = utils.TableLikeFormatter
	} else {
		formatter = utils.SliceMapFormatter
	}
	result, err := formatter.Format(rows)
	if err != nil {
		_ = ctx.AbortWithError(400, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func List(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, service.Q.List())
}
