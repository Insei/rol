package infrastructure

func generateOrderString(orderBy string, orderDirection string) string {
	order := ""
	if len(orderBy) > 0 {
		order = orderBy
		if len(orderDirection) > 0 {
			order = order + " " + orderDirection
		}
	}
	if len(order) < 1 {
		order = "id"
	}
	return order
}
