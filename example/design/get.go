package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("get", func() {
	BasePath("get")
	Response(InternalServerError)

	Action("without", func() {
		Description("Get Method without params")
		Routing(GET(""))
		Response(OK, func() {
			Media(IntTest)
		})
		Response(Unauthorized)
	})
	Action("path_params", func() {
		Description("Get Method with params in path")
		Routing(GET("int/:ParamInt/:ParamStr"))
		Params(func() {
			Param("ParamInt", Integer, func() {
				Description("path_params param int")
				Maximum(10)
			})
			Param("ParamStr", String)
		})
		Response(OK, func() {
			Media(IntTest)
		})
		Response(Unauthorized)
	})

	Action("get_int", func() {
		Description("Get Method")
		Routing(GET("int"))
		Params(func() {
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
		Response(OK)
		Response(Unauthorized)
	})
})
