package infrastructure

import (
	"errors"
	"fmt"
	"regexp"
	"rol/app/interfaces"
	"strings"
)

type GormQueryBuilder struct {
	QueryString string
	Values      []interface{}
}

func NewGormQueryBuilder() *GormQueryBuilder {
	return &GormQueryBuilder{}
}

func toSnakeCase(fieldName string) string {
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snakeName := strings.ToLower(matchAllCap.ReplaceAllString(fieldName, "${1}_${2}"))

	return snakeName
}

func (gormQuery *GormQueryBuilder) addQuery(condition, fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	if len(gormQuery.QueryString) > 0 {
		gormQuery.QueryString += fmt.Sprintf(" %s ", condition)
	}
	gormQuery.QueryString += fmt.Sprintf("%s %s ?", toSnakeCase(fieldName), comparator)
	gormQuery.Values = append(gormQuery.Values, value)
	return gormQuery
}

func (gormQuery *GormQueryBuilder) addQueryBuilder(condition string, builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	if len(gormQuery.QueryString) > 0 {
		gormQuery.QueryString += fmt.Sprintf(" %s ", condition)
	}
	argsInterface, err := builder.Build()
	if err != nil {
		return gormQuery
	}
	argsArrInterface := argsInterface.([]interface{})
	switch v := argsArrInterface[0].(type) {
	case string:
		if len(argsArrInterface[0].(string)) < 1 {
			return gormQuery
		}
		gormQuery.QueryString += fmt.Sprintf("(%s)", strings.ReplaceAll(v, "WHERE ", ""))
	default:
		panic("[GormQueryBuilder] can't add passed query builder to current builder, check what you pass GormQueryBuilder")
	}
	for i := 1; i < len(argsArrInterface); i++ {
		gormQuery.Values = append(gormQuery.Values, argsArrInterface[i])
	}
	return gormQuery
}

func (gormQuery *GormQueryBuilder) Where(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return gormQuery.addQuery("AND", fieldName, comparator, value)
}

func (gormQuery *GormQueryBuilder) WhereQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return gormQuery.addQueryBuilder("AND", builder)
}
func (gormQuery *GormQueryBuilder) Or(fieldName, comparator string, value interface{}) interfaces.IQueryBuilder {
	return gormQuery.addQuery("OR", fieldName, comparator, value)
}
func (gormQuery *GormQueryBuilder) OrQuery(builder interfaces.IQueryBuilder) interfaces.IQueryBuilder {
	return gormQuery.addQueryBuilder("OR", builder)
}

func (gormQuery *GormQueryBuilder) Build() (interface{}, error) {
	if len(gormQuery.QueryString) < 1 {
		return nil, errors.New("queryBuilder is empty")
	}
	arr := make([]interface{}, 0)
	arr = append(arr, gormQuery.QueryString)
	arr = append(arr, gormQuery.Values...)
	return arr, nil
}
