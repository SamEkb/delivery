package queries

type GetAllCouriersQuery struct {
	isValid bool
}

func NewGetAllCouriersQuery() (*GetAllCouriersQuery, error) {
	return &GetAllCouriersQuery{
		isValid: true,
	}, nil
}

func (c *GetAllCouriersQuery) IsValid() bool {
	return c.isValid
}
