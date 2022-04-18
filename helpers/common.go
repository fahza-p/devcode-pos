package helpers

import (
	"strconv"

	"devcode-pos/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Paginate(c *fiber.Ctx, p *models.Pagination, model interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		p.Skip, _ = strconv.Atoi(c.Query("skip"))
		p.Limit, _ = strconv.Atoi(c.Query("limit"))
		categoryId, _ := strconv.Atoi(c.Query("categoryId"))
		searchQuery := c.Query("q")

		if categoryId > 0 {
			db.Where("product_category_id = ?", categoryId)
			models.DB.Where("product_category_id = ?", categoryId)
		}

		if searchQuery != "" {
			db.Where("product_name LIKE ?", "%"+searchQuery+"%")
			models.DB.Where("product_name LIKE ?", "%"+searchQuery+"%")
		}

		models.DB.Model(model).Count(&p.TotalData)

		return db.Offset(p.Skip).Limit(p.Limit)
	}
}
