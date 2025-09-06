package model

type Item struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Quantity int      `json:"quantity"`
	Price    float64  `json:"price"`
	Tags     []string `json:"tags"`
	Status   string   `json:"status"`
}

type ItemCreate struct {
	Name     string   `json:"name" validate:"required,min=1,max=64"`
	Quantity int      `json:"quantity" validate:"required,gte=0"`
	Price    float64  `json:"price" validate:"required,gt=0"`
	Tags     []string `json:"tags" validate:"required"`
	Status   string   `json:"status" validate:"required,oneof=active archived"`
}

type ItemUpdate struct {
	Name     *string   `json:"name,omitempty" validate:"omitempty,min=1,max=64"`
	Quantity *int      `json:"quantity,omitempty" validate:"omitempty,gte=0"`
	Price    *float64  `json:"price,omitempty" validate:"omitempty,gt=0"`
	Tags     *[]string `json:"tags,omitempty" validate:"omitempty"`
	Status   *string   `json:"status,omitempty" validate:"omitempty,oneof=active archived"`
}
