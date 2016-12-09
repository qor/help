package help

import "github.com/qor/admin"

type controller struct {
	Help *Help
}

func (ctr controller) Index(context *admin.Context) {
	helpResource := ctr.Help.Config.Resource
	results := helpResource.NewSlice()

	context.Execute("help/index", map[string]interface{}{
		"HelpResults":  results,
		"HelpResource": helpResource,
	})
}

func (ctr controller) New(context *admin.Context) {
	context.Execute("help/new", context.Resource.NewStruct())
}
