package routes

import (
	"time"

	"health-tracker/handlers"
	"health-tracker/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Global rate limiting (100 requests per minute)
	r.Use(middleware.RateLimitMiddleware(100, time.Minute))

	// API group
	api := r.Group("/api")
	{
		// Public routes (no auth required) with stricter rate limiting
		auth := api.Group("/auth")
		auth.Use(middleware.StrictRateLimitMiddleware()) // 10 requests per minute for auth
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/reset-password", handlers.ResetPassword)
		}

		// Articles routes (public)
		articles := api.Group("/articles")
		{
			articles.GET("", handlers.GetArticles)
			articles.GET("/categories", handlers.GetArticleCategories)
			articles.GET("/search", handlers.SearchArticles)
			articles.GET("/:id", handlers.GetArticle)
		}

		// Protected routes (auth required)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// User routes
			protected.GET("/auth/me", handlers.GetCurrentUser)
			protected.PUT("/auth/profile", handlers.UpdateProfile)

			// Health data routes
			health := protected.Group("/health")
			{
				health.POST("", handlers.CreateHealthData)
				health.GET("", handlers.GetHealthData)
				health.GET("/latest", handlers.GetLatestHealthData)
				health.GET("/dashboard", handlers.GetDashboard)
				health.GET("/graph/:period", handlers.GetHealthGraph)
			}

			// Symptom routes
			symptoms := protected.Group("/symptoms")
			{
				symptoms.GET("/list", handlers.GetSymptomList)
				symptoms.POST("", handlers.LogSymptom)
				symptoms.POST("/batch", handlers.LogMultipleSymptoms)
				symptoms.GET("/history", handlers.GetSymptomHistory)
				symptoms.GET("/stats", handlers.GetSymptomStats)
			}

			// Family routes
			family := protected.Group("/family")
			{
				family.POST("/invite", handlers.InviteFamilyMember)
				family.GET("/members", handlers.GetFamilyMembers)
				family.GET("/requests", handlers.GetFamilyRequests)
				family.PUT("/approve/:id", handlers.ApproveFamilyRequest)
				family.PUT("/reject/:id", handlers.RejectFamilyRequest)
				family.GET("/:id/health", handlers.GetFamilyMemberHealth)
				family.DELETE("/:id", handlers.RemoveFamilyMember)
			}

			// Recommendation routes
			recommendations := protected.Group("/recommendations")
			{
				recommendations.GET("/food", handlers.GetFoodRecommendations)
				recommendations.GET("/exercise", handlers.GetExerciseRecommendations)
				recommendations.GET("/emotional", handlers.GetEmotionalRecommendations)
				recommendations.GET("/daily-menu", handlers.GetDailyMenu)
			}

			// Forum routes
			forum := protected.Group("/forum")
			{
				forum.GET("/posts", handlers.GetPosts)
				forum.POST("/posts", handlers.CreatePost)
				forum.GET("/posts/:id", handlers.GetPost)
				forum.DELETE("/posts/:id", handlers.DeletePost)
				forum.POST("/posts/:id/comments", handlers.AddComment)
				forum.POST("/posts/:id/like", handlers.ToggleLike)
			}

			// Water tracker routes
			water := protected.Group("/water")
			{
				water.GET("", handlers.GetWaterIntake)
				water.POST("/add", handlers.AddWaterGlass)
				water.POST("/remove", handlers.RemoveWaterGlass)
				water.PUT("/goal", handlers.UpdateWaterGoal)
				water.GET("/history", handlers.GetWaterHistory)
			}

			// Goals routes
			goals := protected.Group("/goals")
			{
				goals.GET("", handlers.GetGoals)
				goals.POST("", handlers.CreateGoal)
				goals.PUT("/:id/progress", handlers.UpdateGoalProgress)
				goals.PUT("/:id/toggle", handlers.ToggleGoalComplete)
				goals.DELETE("/:id", handlers.DeleteGoal)
				goals.GET("/stats", handlers.GetGoalStats)
			}

			// Reminders routes
			reminders := protected.Group("/reminders")
			{
				reminders.GET("", handlers.GetReminders)
				reminders.POST("", handlers.CreateReminder)
				reminders.PUT("/:id", handlers.UpdateReminder)
				reminders.DELETE("/:id", handlers.DeleteReminder)
				reminders.PUT("/:id/toggle", handlers.ToggleReminder)
			}
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Health Tracker API is running"})
	})
}

