package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Registrar interface {
   RegisterRoutes(router *gin.Engine)
}

func New(r []Registrar) http.Handler {
   var engine *gin.Engine
   switch os.Getenv("GIN_MODE") {
      case "test":
         // Set Gin to test mode
         gin.SetMode(gin.TestMode)
         engine = gin.New()
      default:
         // Set Gin to release mode
         gin.SetMode(gin.ReleaseMode)
         engine = gin.Default()
   }

   setupRouter(engine, r)
   return engine
}

func setupRouter(router *gin.Engine, r []Registrar){
   for _, registrar := range r {
      registrar.RegisterRoutes(router)
   }
}
