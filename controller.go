package help

import "github.com/qor/admin"

type controller struct {
	Help *Help
}

func (controller) Index(context *admin.Context) {
	context.Execute("help/index", map[string]interface{}{})
}
