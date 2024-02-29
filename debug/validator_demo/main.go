// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/zerolog/log"
// )

// type RegisterRequest struct {
// 	Username string `json:"username" binding:"required"`
// 	Nickname string `json:"nickname" binding:"required"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// 	Age      uint8  `json:"age" binding:"gte=1,lte=120"`
// }

// func main() {

// 	router := gin.Default()

// 	router.POST("register", Register)

// 	router.Run(":9999")
// }

// func Register(c *gin.Context) {
// 	var r RegisterRequest
// 	err := c.ShouldBindJSON(&r)
// 	if err != nil {
// 		log.Info().Msgf("register failed, error is: %v", err.Error())
// 		c.JSON(http.StatusOK, gin.H{"msg": err.Error()})
// 		return
// 	}
// 	//验证 存储操作省略.....
// 	fmt.Println("register success")
// 	c.JSON(http.StatusOK, "successful")
// }

// curl --location --request POST 'http://localhost:9999/register' --header 'Content-Type: application/json' --data-raw '{    "username": "asong",    "nickname": "golang梦工厂",    "email": "7418111@111.com",    "password": "123",    "age": 140}'

package main

// 参考: https://juejin.cn/post/6863765115456454664

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type User struct {
	Username string `validate:"min=6,max=10"`
	Age      uint8  `validate:"gte=1,lte=10"`
	Sex      string `validate:"oneof=female male"`
}

func main() {
	route := gin.Default()
	validate := validator.New()

	user1 := User{Username: "asong", Age: 11, Sex: "null"}
	err := validate.Struct(user1)
	if err != nil {
		fmt.Println(err)
	}

	user2 := User{Username: "asong111", Age: 8, Sex: "male"}
	err = validate.Struct(user2)
	if err != nil {
		fmt.Println(err)
	}

	route.GET("/time", getTime)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("timing", timing)
		if err != nil {
			fmt.Println("success")
		}
	}

	route.Run(":9080")

}

// 自定义验证器
type Info struct {
	CreateTime time.Time `form:"create_time" binding:"required,timing" time_format:"2006-01-02"`
	UpdateTime time.Time `form:"update_time" binding:"required,timing" time_format:"2006-01-02"`
}

// 自定义验证规则断言
func timing(fl validator.FieldLevel) bool {
	if date, ok := fl.Field().Interface().(time.Time); ok {
		today := time.Now()
		log.Info().Msgf("today.After(date) is: %v", today.After(date))
		if today.After(date) {
			return false
		}
	}
	return true
}

func getTime(c *gin.Context) {
	var b Info
	// 数据模型绑定查询字符串验证
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "time are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
