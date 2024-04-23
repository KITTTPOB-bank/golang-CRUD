package routes

import (
	"github.com/KITTTPOB-bank/Carnan_backend/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	//get
	r.GET("/getallcar", controllers.Getallcar)
	r.GET("/cardataadmin", controllers.Getallcaradmin)
	r.GET("/cardata", controllers.Getallcarname)
	r.GET("/carpart/:brand/:name/:year/:check", controllers.GetMyCar)

	r.GET("/getaws3image/:brand/:name/:year", controllers.Getaws3image)

	r.GET("/carpartall", controllers.Getcarpartall)

	//Post
	r.POST("/addcar", controllers.Addcar)
	r.POST("/jsonfileaddcarpart", controllers.AddJsonfileMycarpart)
	r.POST("/jsonfile/:typesend", controllers.AddJsonfilecarpart)
	r.POST("/uploadscars", controllers.Uploadcarimage)

	//delete
	r.DELETE("/dropcar/:delid", controllers.Dropcar)
	r.DELETE("/deleteallpart/:id/:typesend", controllers.DeleteMypart)

	//put
	r.PUT("/addpart", controllers.Addpart)
	r.PUT("/deletepart/", controllers.DropcarPart)
	r.PUT("/addallpart", controllers.Addallpart)

}
