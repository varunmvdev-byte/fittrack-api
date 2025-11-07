package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/varunmvdev-byte/fittrack-api/internal/database"
	"github.com/varunmvdev-byte/fittrack-api/internal/handlers"
	"github.com/varunmvdev-byte/fittrack-api/internal/middleware"
)

func main() {
	_ = godotenv.Load() // load .env if present

	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		// Auth
		auth := handlers.NewAuthHandler(db)
		api.POST("/auth/register", auth.Register)
		api.POST("/auth/login", auth.Login)

		// Protected
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())

		workout := handlers.NewWorkoutHandler(db)
		protected.GET("/workouts", workout.ListWorkouts)
		protected.POST("/workouts", workout.CreateWorkout)
		protected.GET("/workouts/:id", workout.GetWorkout)
		protected.PUT("/workouts/:id", workout.UpdateWorkout)
		protected.DELETE("/workouts/:id", workout.DeleteWorkout)

		protected.POST("/workouts/:id/exercises", workout.AddExercise)
		protected.PUT("/exercises/:id", workout.UpdateExercise)
		protected.DELETE("/exercises/:id", workout.DeleteExercise)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server running on :%s", port)
	r.Run(":" + port)
}
