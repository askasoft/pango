package args

import (
	"github.com/askasoft/pango/asg"
	"github.com/askasoft/pango/str"
)

type Orders struct {
	Order string `json:"o,omitempty" form:"o"`
}

func (o *Orders) String() string {
	return o.Order
}

func (o *Orders) Normalize(columns ...string) {
	if len(columns) == 0 {
		return
	}

	orders := str.FieldsAny(o.String(), ", ")

	orders = asg.DeleteFunc(orders, func(o string) bool {
		return !asg.Contains(columns, str.TrimPrefix(o, "-"))
	})

	o.Order = str.Join(orders, ",")
}
