package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

func defineUserTrait() {
	Trait("UserTrait", func() {
		Attribute("name", String, func() {
			MinLength(4)
			MaxLength(16)
		})
		Attribute("age", Integer, func() {
			Minimum(20)
			Maximum(70)
		})
		Attribute("email", String, func() {
			Pattern(`/^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,4}$/`)
		})
		Attribute("sex", String, func() {
			Enum("male", "female", "other")
		})
		Required("name")
	})
}

var User = MediaType("application/vnd.user+json", func() {
	UseTrait("UserTrait")
	View("default", func() {
		Attribute("name")
		Attribute("email")
		Attribute("age")
		Attribute("sex")
	})
})

var UserCreatePayload = Type("UserCreatePayload", func() {
	UseTrait("UserTrait")
})

var _ = Resource("user", func() {
	BasePath("user")
	Response(InternalServerError)

	Action("list", func() {
		Routing(GET(""))
		Response(OK, func() {
			Media(CollectionOf(User))
		})
		Response(Unauthorized)
	})
	Action("get", func() {
		Routing(GET(":UserID"))
		Params(func() {
			Param("UserID", Integer, func() {
				Maximum(10000)
			})
			Required("UserID")
		})
		Response(OK, func() {
			Media(User)
		})
		Response(Unauthorized)
	})

	Action("create", func() {
		Routing(POST(""))
		Payload(UserCreatePayload, func() {
			Example(map[string]interface{}{
				"price": 15000,
			})
		})
		Response(OK)
	})
})
