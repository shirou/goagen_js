package design

import (
	//	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var IntTest = MediaType("application/vnd.inttest+json", func() {
	Description("Validate test")
	UseTrait("IntTrait")
	View("default", func() {
		Attribute("int")
		Attribute("int_required")
		Attribute("int_max")
		Attribute("int_min")
		Attribute("int_minmax")
		Attribute("int_enum")
		Attribute("int_array")
	})
	View("secret", func() {
		Attribute("int")
		Attribute("int_secret")
	})
})
