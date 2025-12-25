package database

import (
	"health-tracker/models"
	"log"
)

func SeedData() {
	// Check if symptom templates already exist
	var count int64
	DB.Model(&models.SymptomTemplate{}).Count(&count)
	
	if count == 0 {
		log.Println("Seeding symptom templates...")

		// Seed physical symptoms
		for _, symptom := range models.PhysicalSymptoms {
			DB.Create(&models.SymptomTemplate{
				SymptomType: "physical",
				SymptomName: symptom,
				Description: "Gejala fisik: " + symptom,
			})
		}

		// Seed mental symptoms
		for _, symptom := range models.MentalSymptoms {
			DB.Create(&models.SymptomTemplate{
				SymptomType: "mental",
				SymptomName: symptom,
				Description: "Gejala mental: " + symptom,
			})
		}

		// Seed recommendations
		seedRecommendations()
	}

	// Seed articles
	var articleCount int64
	DB.Model(&models.Article{}).Count(&articleCount)
	if articleCount == 0 {
		log.Println("Seeding health articles...")
		articles := models.GetSampleArticles()
		for _, article := range articles {
			DB.Create(&article)
		}
	}

	log.Println("Seed data completed")
}

func seedRecommendations() {
	recommendations := []models.Recommendation{
		// Food recommendations by BMI
		{Category: "food", Condition: "Underweight", Title: "Tingkatkan Asupan Kalori", Description: "Makanan bergizi untuk menambah berat badan", Priority: 1},
		{Category: "food", Condition: "Normal", Title: "Pertahankan Pola Makan Seimbang", Description: "Lanjutkan pola makan sehat Anda", Priority: 2},
		{Category: "food", Condition: "Overweight", Title: "Kurangi Kalori & Tingkatkan Serat", Description: "Makanan rendah kalori dan tinggi serat", Priority: 1},
		{Category: "food", Condition: "Obese", Title: "Diet Ketat dengan Pengawasan", Description: "Konsultasi dengan ahli gizi", Priority: 1},
		
		// Exercise recommendations
		{Category: "exercise", Condition: "sedentary", Title: "Mulai dengan Jalan Kaki", Description: "30 menit jalan kaki setiap hari", Priority: 1},
		{Category: "exercise", Condition: "light", Title: "Tingkatkan Intensitas", Description: "Jogging atau bersepeda ringan", Priority: 2},
		{Category: "exercise", Condition: "moderate", Title: "Variasikan Latihan", Description: "Kombinasi kardio dan kekuatan", Priority: 2},
		{Category: "exercise", Condition: "active", Title: "Pertahankan Rutinitas", Description: "Lanjutkan pola olahraga Anda", Priority: 3},
		
		// Emotional recommendations
		{Category: "emotional", Condition: "stressed", Title: "Teknik Relaksasi", Description: "Meditasi dan pernapasan dalam", Priority: 1},
		{Category: "emotional", Condition: "anxious", Title: "Kelola Kecemasan", Description: "Latihan grounding dan mindfulness", Priority: 1},
		{Category: "emotional", Condition: "sad", Title: "Aktivitas Menyenangkan", Description: "Lakukan hobi yang Anda sukai", Priority: 1},
		{Category: "emotional", Condition: "happy", Title: "Pertahankan Mood Positif", Description: "Lanjutkan aktivitas yang membuat bahagia", Priority: 3},
	}

	for _, rec := range recommendations {
		DB.Create(&rec)
	}
}
