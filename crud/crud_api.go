package crud

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type APICrudModel interface {
	GetDatabase() *gorm.DB
	GetActiveConditionMap() map[string]interface{}
	GetID() int64
	Delete() (bool, error)
	List() interface{}
}

func CrudCreator(m APICrudModel) map[string]func(c *gin.Context) {
	return map[string]func(c *gin.Context){
		http.MethodGet:    getModel(m),
		http.MethodPost:   createModel(m),
		http.MethodDelete: deleteModel(m),
		http.MethodPut:    updateModel(m),
	}
}

func GetLimitOffset(c *gin.Context) (int, int) {
	limitRaw := c.Query("limit")
	limit, err := strconv.Atoi(limitRaw)
	if err != nil {
		limit = 0
	}
	offsetRaw := c.Query("offset")
	offset, err := strconv.Atoi(offsetRaw)
	if err != nil {
		offset = 0
	}
	return limit, offset
}
func GetID(c *gin.Context) (int64, error) {
	IdRaw := c.Param("id")
	id, err := strconv.ParseInt(IdRaw, 0, 10)
	return int64(id), err
}

func getModel(m APICrudModel) func(c *gin.Context) {
	return func(c *gin.Context) {
		IdRaw := c.Param("id")
		var (
			res interface{}
			err error
		)
		dbManager := m.GetDatabase()
		q := dbManager.Where(m.GetActiveConditionMap())
		if IdRaw == "" {
			limit, offset := GetLimitOffset(c)
			res = m.List()
			if limit > 0 {
				q.Limit(limit)
			}
			if offset > 0 {
				q.Offset(offset)
			}
			err = q.Find(res).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		} else {
			id, err_ := GetID(c)
			if err_ != nil {
				c.JSON(http.StatusBadRequest, err)
			}
			err = q.First(m, id).Error
			res = m
			if err != nil {
				c.JSON(http.StatusNotFound, err)
				return
			}
		}
		c.JSON(http.StatusOK, res)
	}
}

func createModel(m APICrudModel) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := c.ShouldBindJSON(&m)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		dbManager := m.GetDatabase()
		err = dbManager.Create(m).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = dbManager.First(m, m.GetID()).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, m)
	}
}

func updateModel(m APICrudModel) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := GetID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		err = c.ShouldBindJSON(&m)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		if id != m.GetID() {
			c.JSON(http.StatusBadRequest, "url id is not equal with body id")
			return
		}
		dbManager := m.GetDatabase()
		err = dbManager.Model(m).Where("id=?", id).Updates(m).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		err = dbManager.First(m, m.GetID()).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, m)
	}
}

func deleteModel(m APICrudModel) func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := GetID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		dbManager := m.GetDatabase()
		soft, err := m.Delete()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		if soft {
			err = dbManager.Model(m).Where("id=?", id).Updates(m).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, err)
				return
			}
		}
		c.Status(http.StatusOK)
	}
}
