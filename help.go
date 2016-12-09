package help

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

type QorHelpEntry struct {
	gorm.Model
	Title string
	Body  string `gorm:"size:65532"`
}

func (QorHelpEntry) ResourceName() string {
	return "Help"
}

func (QorHelpEntry) ToParam() string {
	return "!help"
}

func (qorHelpEntry *QorHelpEntry) ConfigureQorResource(res resource.Resourcer) {
	if res, ok := res.(*admin.Resource); ok {
		Admin := res.GetAdmin()
		router := Admin.GetRouter()

		Admin.RegisterViewPath("github.com/qor/help/views")

		helpController := controller{}
		router.Get("!help", helpController.Index, admin.RouteConfig{
			PermissionMode: roles.Read,
			Resource:       res,
		})

		router.Get("!help/new", helpController.New, admin.RouteConfig{
			PermissionMode: roles.Create,
			Resource:       res,
		})
	}
}
