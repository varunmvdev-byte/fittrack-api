package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/varunmvdev-byte/fittrack-api/internal/models"
)

type WorkoutHandler struct{ db *gorm.DB }

func NewWorkoutHandler(db *gorm.DB) *WorkoutHandler { return &WorkoutHandler{db: db} }

func idParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err
}

type workoutCreateReq struct {
	Date  string `json:"date" binding:"required"` // ISO RFC3339 date
	Notes string `json:"notes"`
}

type exerciseReq struct {
	Name   string  `json:"name" binding:"required"`
	Sets   int     `json:"sets"`
	Reps   int     `json:"reps"`
	Weight float64 `json:"weight"`
}

func (h *WorkoutHandler) ListWorkouts(c *gin.Context) {
	userID := c.GetUint("userID")
	var workouts []models.Workout

	if err := h.db.Where("user_id = ?", userID).
		Preload("Exercises").
		Order("date desc").
		Find(&workouts).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) CreateWorkout(c *gin.Context) {
	userID := c.GetUint("userID")

	var req workoutCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use RFC3339"})
		return
	}

	w := models.Workout{
		UserID: userID, // ✅ Store owner
		Date:   date,
		Notes:  req.Notes,
	}

	if err := h.db.Create(&w).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func (h *WorkoutHandler) GetWorkout(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := idParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var w models.Workout
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).
		Preload("Exercises").
		First(&w).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) UpdateWorkout(c *gin.Context) {
	userID := c.GetUint("userID")
	id, err := idParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var w models.Workout
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&w).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req workoutCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Date != "" {
		date, err := time.Parse(time.RFC3339, req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
			return
		}
		w.Date = date
	}

	w.Notes = req.Notes

	if err := h.db.Save(&w).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	id, err := idParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.db.Delete(&models.Workout{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *WorkoutHandler) AddExercise(c *gin.Context) {
	id, err := idParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workout id"})
		return
	}

	var req exerciseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ex := models.Exercise{
		WorkoutID: id,
		Name:      req.Name,
		Sets:      req.Sets,
		Reps:      req.Reps,
		Weight:    req.Weight,
	}

	if err := h.db.Create(&ex).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusCreated, ex)
}

func (h *WorkoutHandler) UpdateExercise(c *gin.Context) {
	exID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req exerciseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ex models.Exercise
	if err := h.db.First(&ex, uint(exID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	ex.Name = req.Name
	ex.Sets = req.Sets
	ex.Reps = req.Reps
	ex.Weight = req.Weight

	if err := h.db.Save(&ex).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, ex)
}

func (h *WorkoutHandler) DeleteExercise(c *gin.Context) {
	exID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.db.Delete(&models.Exercise{}, uint(exID)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.Status(http.StatusNoContent)
}
