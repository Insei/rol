package interfaces

type IQueryBuilder interface {
	Where(fieldName, comparator string, value interface{}) IQueryBuilder
	WhereQuery(builder IQueryBuilder) IQueryBuilder
	Or(fieldName, comparator string, value interface{}) IQueryBuilder
	OrQuery(builder IQueryBuilder) IQueryBuilder
	Build() (interface{}, error)
}
