package validate

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/yourname/inventory-go/internal/model"
)

var V = validator.New()

type Errors map[string]string

func (e Errors) Error() string { return "validation failed" }

func ValidateCreateItem(in *model.ItemCreate) error {
	if err := V.Struct(in); err != nil {
		return toErrors(err)
	}
	return nil
}

func ValidateUpdateItem(in *model.ItemUpdate) error {
	// Валидация только тех полей, что пришли
	if in.Name != nil {
		if err := V.Var(*in.Name, "min=1,max=64"); err != nil {
			return toErrors(err)
		}
	}
	if in.Quantity != nil {
		if err := V.Var(*in.Quantity, "gte=0"); err != nil {
			return toErrors(err)
		}
	}
	if in.Price != nil {
		if err := V.Var(*in.Price, "gt=0"); err != nil {
			return toErrors(err)
		}
	}
	if in.Status != nil {
		if err := V.Var(*in.Status, "oneof=active archived"); err != nil {
			return toErrors(err)
		}
	}
	// tags — допускаем любой массив строк
	return nil
}

func toErrors(err error) error {
	var verr validator.ValidationErrors
	if errors.As(err, &verr) {
		out := Errors{}
		for _, fe := range verr {
			out[fe.Field()] = fe.Tag()
		}
		return out
	}
	return err
}
