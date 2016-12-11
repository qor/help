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

var Global = "dashboard"

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
					tx := context.GetDB()
					if tx.NewRecord(record) {
						if category := context.Request.URL.Query().Get("category"); category != "" {
							return []string{category}
						}
					}

					if field, ok := tx.NewScope(record).FieldByName("Categories"); ok {
						if categories, ok := field.Field.Interface().(Categories); ok {
							return categories.Categories
						}
					}
					return []string{}
				},
				Config: &admin.SelectManyConfig{Collection: func(record interface{}, context *qor.Context) [][]string {
					var results [][]string
					for _, r := range Admin.GetResources() {
						results = append(results, []string{r.ToParam(), string(Admin.T(context, fmt.Sprintf("qor_help.categories.%v", r.ToParam()), r.Name))})
					}
					return results
				}},
			})
		}

		res.IndexAttrs() // generate search attrs
		searchHandler := res.SearchHandler
		res.SearchHandler = func(keyword string, context *qor.Context) *gorm.DB {
			tx := searchHandler(keyword, context)
			if category := context.Request.URL.Query().Get("category"); category != "" {
				return tx.Where("categories LIKE ?", "%"+fmt.Sprintf("[%v]", category)+"%")
			}
			return tx
		}

		res.ShowAttrs("Body")

		Admin.RegisterViewPath("github.com/qor/help/views")
		Admin.RegisterResourceRouters(res, "create", "update", "read", "delete")

		Admin.RegisterFuncMap("get_help_categories", func(context *admin.Context) [][]string {
			var results [][]string
			for _, r := range Admin.GetResources() {
				results = append(results, []string{r.ToParam(), string(Admin.T(context.Context, fmt.Sprintf("qor_help.categories.%v", r.ToParam()), r.Name))})
			}
			return results
		})

		Admin.RegisterFuncMap("get_current_help_category", func(r *admin.Resource, context *admin.Context) string {
			if r != nil {
				return r.ToParam()
			}

			if category := context.Request.URL.Query().Get("category"); category != "" {
				return category
			}

			return ""
		})

		Admin.RegisterFuncMap("get_help_documents_count", func(r *admin.Resource, context *admin.Context) uint {
			var result uint
			if r != nil {
				context.GetDB().Model(res.NewStruct()).Where("categories LIKE ?", "%"+fmt.Sprintf("[%v]", r.ToParam())+"%").Count(&result)
			}
			return result
		})
	}
}
