package help

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor/resource"
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
		res.UseTheme("help")
		res.ShowAttrs("Body")

		Admin := res.GetAdmin()
		Admin.RegisterViewPath("github.com/qor/help/views")
		Admin.RegisterResourceRouters(res, "create", "update", "read", "delete")
	}
}
