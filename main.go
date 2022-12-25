package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type User struct {
	ID               int
	Name             string
	Password         string
	SecretProtection string
}

var registerID, loginID, postID int
var registerName, registerPassword, registerSecretProtection string
var loginName, loginPassword, loginSecretProtection string
var message, postUerName string

// 定义一个全局对象db
var db *sql.DB

// 连接数据库库
func initDB() (err error) {
	// dsn := "root:root@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4&parseTime=True"
	dsn := "root:xian712525@tcp(127.0.0.1:3306)/go_db?charset=utf8mb4"
	// open函数只是验证格式是否正确，并不是创建数据库连接
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		return err2
	}
	return nil
}
func IsDB() {
	//验证是否成功
	err := initDB()
	if err != nil {
		fmt.Printf("*********err: %v\n", err)
	} else {
		fmt.Println("***************连接成功***************")
	}
	fmt.Printf("db: %v\n", db)
}

// 注册用户
func insertData(i int, n string, p string, ps string) {
	sqlStr := "insert into userMassage(ID,username,password,secretProtection) values (?,?,?,?)"
	r, err := db.Exec(sqlStr, i, n, p, ps)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	ID, err2 := r.LastInsertId()
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
		return
	}
	fmt.Printf("ID: %v\n", ID)
}

// 查询密码
func queryManyData() bool {
	is := false
	sqlStr := "select username,password from usermassage"
	r, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	defer r.Close()
	// 循环读取结果集中的数据
	for r.Next() {
		var u2 User
		err2 := r.Scan(&u2.Name, &u2.Password)
		if err2 != nil {
			fmt.Printf("err: %v\n", err2)
		}
		if u2.Name == loginName && u2.Password == loginPassword {
			is = true
			break
		}
	}
	return is
}
func main() {
	IsDB() //验证连接数据库是否成功
	r := gin.Default()
	//注册
	r.POST("/register", func(c *gin.Context) {
		//获取数据
		ri := c.PostForm("ID")
		registerID, _ = strconv.Atoi(ri)
		registerName = c.PostForm("name")
		registerPassword = c.PostForm("password")
		registerSecretProtection = c.PostForm("secretProtection")
		//添加注册信息到文件
		insertData(registerID, registerName, registerPassword, registerSecretProtection)
		//返回结果
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "注册成功",
		})
	})
	//登录界面
	r.GET("/login", func(c *gin.Context) {
		c.SetCookie("login", "yes", 60, "/", "localhost", false, true)
		// 返回信息
		c.String(200, "Login success!")
	})
	//登录验证
	r.POST("/login", func(c *gin.Context) {
		//获取参数
		loginName = c.PostForm("name")
		loginPassword = c.PostForm("password")
		//验证(只需要用户名和密码)
		is := queryManyData()
		if is {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "登录成功",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    404,
				"message": "登录失败",
			})
			return
		}

	})
	//新增仓库
	r.POST("/new-house", func(c *gin.Context) {
		// 解析请求参数
		houseID := c.PostForm("houseID")
		name := c.PostForm("name")
		address := c.PostForm("address")

		// 将仓库信息写入数据库
		_, err := db.Exec("insert into warehouses (houseID,name, address) values (?,?, ?)", houseID, name, address)
		if err != nil {
			// 写入数据库失败，返回错误响应
			c.JSON(http.StatusOK, gin.H{
				"code":    404,
				"message": "创建失败",
			})
			return
		}

		// 写入数据库成功，返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "创建成功",
		})
	})
	//商品的上下架

	r.POST("/new-product", func(c *gin.Context) {

		warehouseID := c.PostForm("warehouseID") //所属仓库
		productID := c.PostForm("productID")
		money := c.PostForm("money")
		status := c.PostForm("status") // 上架或下架
                 _, err1 := db.Exec("insert into new-products (warehouseID,productID,money) values (?,?,?)", warehouseID, productID, money)
		if err1 != nil {
			fmt.Println(err1)
		}
		// 更新数据库中货物的状态
		_, err := db.Exec("update new-products set status = ? where warehouse_id = ? and product_id = ? and money =?", status, warehouseID, productID, money)
		if err != nil {
			// 更新数据库失败，返回错误响应
			c.JSON(http.StatusOK, gin.H{
				"code":    404,
				"message": "更新商品失败",
			})
			return
		}

		// 更新数据库成功，返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "更新商品成功",
		})
	})

	r.Run(":/9090")
}
