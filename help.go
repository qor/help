package help

import (
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

type Help struct {
	Config *Config
}

type Config struct {
}

func New(config *Config) *Help {
	return &Help{Config: config}
}

func (help *Help) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		Admin := res.GetAdmin()
		router := Admin.GetRouter()

		Admin.RegisterViewPath("github.com/qor/help/views")

		helpController := controller{Help: help}
		router.Get("!help", helpController.Index, admin.RouteConfig{
			PermissionMode: roles.Read,
			Resource:       res,
		})
	}
}
