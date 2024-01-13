package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// create schema crud collate utf8mb4_0900_ai_ci;

func main() {
	dsn := "root:Yizhili80@tcp(127.0.0.1:3306)/crud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 设置表名为单数
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	fmt.Println(db)  // &{0xc0001762d0 <nil> 0 0xc000232000 1}
	fmt.Println(err) // <nil>

	// 结构体
	type List struct {
		gorm.Model
		Name    string `gorm:"type:varchar(20); not null" json:"name" binding:"required"`
		State   string `gorm:"type:varchar(20); not null" json:"state" binding:"required"`
		Phone   string `gorm:"type:varchar(20); not null" json:"phone" binding:"required"`
		Email   string `gorm:"type:varchar(40); not null" json:"email" binding:"required"`
		Address string `gorm:"type:varchar(200); not null" json:"address" binding:"required"`
	}

	// AutoMigrate 创建表、缺失的外键、约束、列和索引。
	err = db.AutoMigrate(&List{})
	if err != nil {
		panic(err.Error())
	}

	// GORM 使用 database/sql 来维护连接池
	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second) // 10 秒

	r := gin.Default()
	// 增
	r.POST("/user/add", func(c *gin.Context) {
		var data List
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": "400",
			})
		} else {
			// 数据库操作
			db.Create(&data)
			c.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": "200",
			})
		}
	})

	// 删
	// 1. 找到 id 对应的数据条目
	// 2. 判断 id 是否存在
	// 3. 从数据库中删除
	// 4. 返回，id 没有找到
	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		var data []List
		// 接收 id
		id := c.Param("id")
		// 判断 id 是否存在
		db.Where("id = ?", id).Find(&data) // 从data 指针指向的表查询 id 是否存在
		// id 存在则删除，不存在则报错
		if len(data) == 0 {
			c.JSON(200, gin.H{
				"msg":  "id 没有找到，删除失败",
				"code": 400,
			})
		} else {
			// 操作数据库删除
			db.Delete(&data, id)
			//db.Where("id = ?", id).Delete(&data)

			c.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}
	})

	// 改
	r.PUT("/user/update/:id", func(c *gin.Context) {
		// 1. 找到 id 对应的数据条目
		// 2. 判断 id 是否存在
		// 3. 更新
		// 4. 返回，id 没有找到
		var data List
		id := c.Param("id")
		// 2. 判断 id 是否存在
		db.Select("id").Where("id = ?", id).Find(&data)
		if data.ID == 0 {
			c.JSON(200, gin.H{
				"msg":  "用户id没有找到",
				"code": 400,
			})
		} else {
			err := c.ShouldBindJSON(&data)
			// Update with struct
			db.Where("id = ?", id).Updates(&data)
			// UPDATE users WHERE id = 'id';

			if err != nil {
				c.JSON(200, gin.H{
					"msg":  "更新失败",
					"code": 400,
				})
			} else {
				c.JSON(200, gin.H{
					"msg":  "更新完成",
					"code": 200,
				})
			}
		}
	})

	// 查 (1. 条件查询 2. 全部查询、分页查询）
	// 1. 条件查询
	r.GET("/user/list/:name", func(c *gin.Context) {
		// 获取路径参数
		name := c.Param("name")
		var dataList []List
		// 查询数据库
		db.Where("name = ?", name).Find(&dataList)

		// 判断是否查询到
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"msg":  "没有查找到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "查找到数据",
				"code": 200,
				"data": dataList,
			})
		}
	})

	// 2. 全部查询、分页查询

	r.GET("/user/list", func(c *gin.Context) {
		var dataList []List
		// 返回一个总数
		var total int64

		pageSize, _ := strconv.Atoi(c.Query("pageSize"))
		pageNum, _ := strconv.Atoi(c.Query("pageNum"))
		//fmt.Println("pageSize", pageSize)
		//fmt.Println("pageNum", pageNum)

		// 判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}
		// 固定写法
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}

		// 查询数据库
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList) == 0 {
			c.JSON(200, gin.H{
				"code": 400,
				"msg":  "查询失败！",
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"code": 200,
				"msg":  "查询成功！",
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})
		}
	})

	err = r.Run(":8080")
	if err != nil {
		return
	}
}
