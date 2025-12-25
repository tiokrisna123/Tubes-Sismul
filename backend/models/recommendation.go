package models

type Recommendation struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Category    string `json:"category"` // food, exercise, emotional
	Condition   string `json:"condition"` // BMI category, symptom, emotional state
	Title       string `json:"title"`
	Description string `json:"description"`
	Details     string `gorm:"type:text" json:"details"`
	Priority    int    `json:"priority"` // 1-5
}

type FoodRecommendation struct {
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Foods       []string `json:"foods"`
	Avoid       []string `json:"avoid"`
	Reason      string   `json:"reason"`
}

type ExerciseRecommendation struct {
	Category    string   `json:"category"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Exercises   []string `json:"exercises"`
	Duration    string   `json:"duration"`
	Frequency   string   `json:"frequency"`
	Intensity   string   `json:"intensity"`
	Reason      string   `json:"reason"`
}

type EmotionalRecommendation struct {
	EmotionalState string   `json:"emotional_state"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Activities     []string `json:"activities"`
	Tips           []string `json:"tips"`
	Reason         string   `json:"reason"`
}

type MealPlan struct {
	MealType      string   `json:"meal_type"` // breakfast, lunch, dinner, snack
	Title         string   `json:"title"`
	Foods         []string `json:"foods"`
	Ingredients   []string `json:"ingredients,omitempty"`
	Recipe        string   `json:"recipe,omitempty"`
	Calories      string   `json:"calories,omitempty"`
	Description   string   `json:"description"`
	EstimatedCost string   `json:"estimated_cost,omitempty"` // Estimasi biaya dalam Rupiah
}

type DailyMenu struct {
	Date               string     `json:"date"`
	HealthTip          string     `json:"health_tip"`
	Breakfast          MealPlan   `json:"breakfast"`
	BreakfastAlt       []MealPlan `json:"breakfast_alt,omitempty"` // Alternative breakfast options
	Lunch              MealPlan   `json:"lunch"`
	LunchAlt           []MealPlan `json:"lunch_alt,omitempty"` // Alternative lunch options
	Dinner             MealPlan   `json:"dinner"`
	DinnerAlt          []MealPlan `json:"dinner_alt,omitempty"` // Alternative dinner options
	Snacks             []MealPlan `json:"snacks"`
	Drinks             []string   `json:"drinks"`
	Fruits             []string   `json:"fruits"`
	AvoidDrinks        []string   `json:"avoid_drinks,omitempty"`
	AvoidFruits        []string   `json:"avoid_fruits,omitempty"`
	TotalCalories      string     `json:"total_calories"`
	TotalEstimatedCost string     `json:"total_estimated_cost,omitempty"` // Total estimasi biaya harian
}
