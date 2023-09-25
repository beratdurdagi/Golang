package routes

import (
	"log"
	"os"

	controllers "github.com/Karalakrepp/Golang/Ecommerce_GO/controller"
	database "github.com/Karalakrepp/Golang/Ecommerce_GO/db"
	"github.com/Karalakrepp/Golang/Ecommerce_GO/middleware"
	"github.com/gin-gonic/gin"
)

type APIServer struct {
	listenAddr string
}

var (
	db      = database.DBSet()
	logger  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	handler = controllers.NewHandlerAPI(logger, db)
	client  = db.Client
)

func NewAPIServer(l *log.Logger, listenaddr string) *APIServer {
	return &APIServer{

		listenAddr: listenaddr,
	}
}

// this is subRouter before login
func (s *APIServer) NewRouter(subRouter *gin.Engine) {

	subRouter.POST("/users/register", handler.SignUp())
	subRouter.POST("/users/login", handler.Login())
	subRouter.POST("/admin/addProduct", handler.ProductViewerAdmin())
	subRouter.GET("/users/productview", handler.SeachProduct())
	subRouter.GET("/users/search", handler.SearchProductByQuery())

}

// It starts HTTP and use middleware authentication for login and include after login routes
func (s *APIServer) Run() {
	app := controllers.NewApplication(db.ProductData(client, "Products"), db.UserData(client, "Users"))

	router := gin.New()

	router.Use(gin.Logger())

	s.NewRouter(router)

	router.Use(middleware.JWT_Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controllers.GetItemFromCart())
	router.POST("/addaddress", handler.AddAddress())
	router.PUT("/edithomeaddress", handler.EditHomeAddress())
	router.PUT("/editworkaddress", handler.EditWorkAddress())
	router.GET("/deleteaddresses", handler.DeleteAddress())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(s.listenAddr))
}
