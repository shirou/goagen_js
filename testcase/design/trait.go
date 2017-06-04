package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

func define_int_trait() {
	Trait("IntTrait", func() {
		Attributes(func() {
			Attribute("int", Integer)
			Attribute("int_required", Integer)
			Attribute("int_max", Integer, func() {
				Description("max")
				Maximum(10)
			})
			Attribute("int_min", Integer, func() {
				Minimum(-1)
			})
			Attribute("int_minmax", Integer, func() {
				Minimum(0)
				Maximum(10)
			})
			Attribute("int_enum", Integer, func() {
				Enum(1, 2, 3)
			})
			Attribute("int_array", ArrayOf(Integer))
			Attribute("int_secret", Integer, func() {
				Description("not included in default")
				Example(0)
			})

			Required("int_required")
		})
	})
}
