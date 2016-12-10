package help

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
)

type QorHelpEntry struct {
	gorm.Model
	Title      string
	Categories Categories
	Body       string `gorm:"size:65532"`
}

type Categories struct {
	RawValue   string
	Categories []string
}

func (categories *Categories) Scan(data interface{}) (err error) {
	switch values := data.(type) {
	case []byte:
		if string(values) != "" {
			var strs []string
			for _, c := range strings.Split(string(values), ";") {
				strs = append(strs, strings.TrimFunc(c, func(v rune) bool {
					r := strings.TrimSpace(string(v))
					return r == "" || r == "[" || r == "]"
				}))
			}
			categories.Scan(strs)
		}
	case string:
		return categories.Scan([]byte(values))
	case []string:
		var strs []string
		for _, v := range values {
			if strings.TrimSpace(v) != "" {
				strs = append(strs, v)
			}
		}
		sort.Strings(strs)

		categories.Categories = strs
	default:
		err = errors.New("unsupported driver -> Scan pair for Categories")
	}
	return
}

func (categories Categories) Value() (driver.Value, error) {
	var cs []string
	for _, c := range categories.Categories {
		cs = append(cs, fmt.Sprintf("[%v]", c))
	}
	return strings.Join(cs, "; "), nil
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
		res.UseTheme("help")

		if res.GetMeta("Body") == nil {
			res.Meta(&admin.Meta{Name: "Body", Type: "rich_editor"})
		}

		if res.GetMeta("Categories") == nil {
			res.Meta(&admin.Meta{
				Name: "Categories",
				Valuer: func(record interface{}, context *qor.Context) interface{} {
					if field, ok := context.GetDB().NewScope(record).FieldByName("Categories"); ok {
						if categories, ok := field.Field.Interface().(Categories); ok {
							return categories.Categories
						}
					}
					return []string{}
				},
				Config: &admin.SelectManyConfig{Collection: func(record interface{}, context *qor.Context) [][]string {
					var results = [][]string{{"dashboard", string(Admin.T(context, "qor_help.categories.dashboard", "Dashboard"))}}
					for _, r := range Admin.GetResources() {
						results = append(results, []string{r.ToParam(), string(Admin.T(context, fmt.Sprintf("qor_help.categories.%v", r.ToParam()), r.Name))})
					}
					return results
				}}})
		}

		res.ShowAttrs("Body")

		Admin.RegisterViewPath("github.com/qor/help/views")
		Admin.RegisterResourceRouters(res, "create", "update", "read", "delete")
	}
}
