package help

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

type Help struct {
	Config *Config
}

type Config struct {
	Resource *admin.Resource
}

type QorHelpEntry struct {
	gorm.Model
	Title string
	Body  string `gorm:"size:65532"`
}

func New(config *Config) *Help {
	return &Help{Config: config}
}

func (help *Help) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		Admin := res.GetAdmin()
		router := Admin.GetRouter()

		if help.Config.Resource == nil {
			help.Config.Resource = Admin.NewResource(&QorHelpEntry{})
			Admin.Config.DB.AutoMigrate(&QorHelpEntry{})
		}

		Admin.RegisterViewPath("github.com/qor/help/views")

		helpController := controller{Help: help}
		router.Get("!help", helpController.Index, admin.RouteConfig{
			PermissionMode: roles.Read,
			Resource:       res,
		})
	}
}
