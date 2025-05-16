package queries

type GetAllUncompletedOrdersQuery struct {
	isValid bool
}

func NewGetAllUncompletedOrdersQuery() (*GetAllUncompletedOrdersQuery, error) {
	return &GetAllUncompletedOrdersQuery{
		isValid: true,
	}, nil
}

func (c *GetAllUncompletedOrdersQuery) IsValid() bool {
	return c.isValid
}
