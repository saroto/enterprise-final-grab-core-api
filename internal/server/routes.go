package server

import (
	controller "enterprise_core/internal/controllers"
	"net/http"

	"enterprise_core/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

		// // Get allowed origins from environment variable
		// allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		// if allowedOrigins == "" {
		// 	allowedOrigins = "http://localhost:5173" // Default value
		// }
		   // Get allowed origins from environment variable (if set)
		//    allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		//    if allowedOrigins == "" {
		// 	   allowedOrigins = "*" // Allow all origins if not specified
		//    }

		   
		// Convert comma-separated origins into a slice
		// origins := strings.Split(allowedOrigins, ",")
		
	// r.Use(cors.New(cors.Config{
	// 	// AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 	AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
	// 	AllowCredentials: true, // Enable cookies/auth
	// }))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
	}))

	r.Use(middleware.LogRequestMiddleware())
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/login", controller.LoginUser(s.db))
	r.POST("/register", controller.RegisterUser(s.db))
	
	// Public routes (No authentication needed)
	r.POST("/acc", controller.CreateAccount(s.db))

	// Authenticated routes (Require AuthMiddleware)
	r.GET("/auth-test", middleware.AuthMiddleware(s.db), controller.AuthTest())
	r.GET("/getMe", middleware.AuthMiddleware(s.db), controller.GetMe(s.db))
	r.PUT("/update-user/:id", middleware.AuthMiddleware(s.db), controller.UpdateUser(s.db))
	r.GET("/getSSO", middleware.AuthMiddleware(s.db), controller.GetSSO(s.db))
	r.POST("/logout", middleware.AuthMiddleware(s.db), controller.LogoutUser(s.db))
	r.GET("/report", middleware.AuthMiddleware(s.db), controller.GetUserReport(s.db))

		
	r.GET("/all-users", middleware.AuthMiddleware(s.db), controller.GetAllUsers(s.db))
	
	
	// Account routes with authentication middleware
	r.POST("/accounts", middleware.AuthMiddleware(s.db), controller.CreateAccount(s.db))
	r.GET("/accounts", middleware.AuthMiddleware(s.db), controller.GetAllAccounts(s.db))
	r.GET("/own-accounts", middleware.AuthMiddleware(s.db), controller.GetOwnAccount(s.db))
	r.GET("/accounts/:id", middleware.AuthMiddleware(s.db), controller.GetAccount(s.db))
	r.PUT("/accounts/:id", middleware.AuthMiddleware(s.db), controller.UpdateAccount(s.db))
	r.DELETE("/accounts/:id", middleware.AuthMiddleware(s.db), controller.DeleteAccount(s.db))
	
	// Transaction routes with authentication middleware
	r.POST("/transactions", middleware.AuthMiddleware(s.db), controller.CreateTransaction(s.db))
	r.GET("/transactions", middleware.AuthMiddleware(s.db), controller.GetAllTransactions(s.db))
	r.GET("/transactions/:id", middleware.AuthMiddleware(s.db), controller.GetTransaction(s.db))
	r.PUT("/transactions/:id", middleware.AuthMiddleware(s.db), controller.UpdateTransaction(s.db))
	r.DELETE("/transactions/:id", middleware.AuthMiddleware(s.db), controller.DeleteTransaction(s.db))
	

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello Core Enterprise"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
