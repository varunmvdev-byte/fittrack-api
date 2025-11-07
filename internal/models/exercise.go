package models

import "time"

type Exercise struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	WorkoutID uint      `json:"workout_id"`
	Name      string    `json:"name"`
	Sets      int       `json:"sets"`
	Reps      int       `json:"reps"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
