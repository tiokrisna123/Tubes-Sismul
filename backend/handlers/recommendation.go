package handlers

import (
	"net/http"

	"health-tracker/database"
	"health-tracker/models"
	"health-tracker/utils"

	"github.com/gin-gonic/gin"
)

// GetFoodRecommendations returns personalized food recommendations
func GetFoodRecommendations(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get latest health data
	var health models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&health)

	// Get recent symptoms
	var symptoms []models.Symptom
	database.DB.Where("user_id = ?", userID).Order("logged_at desc").Limit(10).Find(&symptoms)

	recommendations := generateFoodRecommendations(health, symptoms)

	utils.SuccessResponse(c, http.StatusOK, "Food recommendations retrieved", recommendations)
}

// GetExerciseRecommendations returns personalized exercise recommendations
func GetExerciseRecommendations(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get user profile
	var user models.User
	database.DB.First(&user, userID)

	// Get latest health data
	var health models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&health)

	// Get recent symptoms
	var symptoms []models.Symptom
	database.DB.Where("user_id = ?", userID).Order("logged_at desc").Limit(10).Find(&symptoms)

	recommendations := generateExerciseRecommendations(user, health, symptoms)

	utils.SuccessResponse(c, http.StatusOK, "Exercise recommendations retrieved", recommendations)
}

// GetEmotionalRecommendations returns emotional activity recommendations
func GetEmotionalRecommendations(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get latest health data for emotional state
	var health models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&health)

	// Get mental symptoms
	var mentalSymptoms []models.Symptom
	database.DB.Where("user_id = ? AND symptom_type = ?", userID, "mental").
		Order("logged_at desc").Limit(5).Find(&mentalSymptoms)

	recommendations := generateEmotionalRecommendations(health.EmotionalState, mentalSymptoms)

	utils.SuccessResponse(c, http.StatusOK, "Emotional recommendations retrieved", recommendations)
}

// GetDailyMenu returns personalized daily menu based on health condition
func GetDailyMenu(c *gin.Context) {
	userID := c.GetUint("userID")

	// Get latest health data
	var health models.HealthData
	database.DB.Where("user_id = ?", userID).Order("record_date desc").First(&health)

	// Get recent symptoms
	var symptoms []models.Symptom
	database.DB.Where("user_id = ?", userID).Order("logged_at desc").Limit(10).Find(&symptoms)

	menu := generateDailyMenu(health, symptoms)

	utils.SuccessResponse(c, http.StatusOK, "Daily menu generated", menu)
}

func generateDailyMenu(health models.HealthData, symptoms []models.Symptom) models.DailyMenu {
	bmiCategory := models.GetBMICategory(health.BMI)
	
	// Check symptoms
	symptomNames := make(map[string]bool)
	for _, s := range symptoms {
		symptomNames[s.SymptomName] = true
	}

	var menu models.DailyMenu
	menu.Date = "Hari Ini"

	// Default healthy menu with detailed recipes - Indonesian style
	menu.Breakfast = models.MealPlan{
		MealType:      "breakfast",
		Title:         "🌅 Bubur Ayam Kampung Sehat",
		Foods:         []string{"Bubur ayam kampung", "Telur setengah matang", "Jus jeruk segar"},
		Ingredients:   []string{"Beras 100g", "Ayam kampung 100g", "Daun bawang 2 batang", "Bawang goreng 1 sdm", "Kecap manis 1 sdt", "Kerupuk emping", "Telur 1 butir", "Jeruk 2 buah"},
		Description:   "Sarapan tradisional Indonesia yang hangat dan bergizi",
		Calories:      "~450 kkal",
		EstimatedCost: "Rp 18.000 - 25.000",
		Recipe:        "1. Masak bubur dengan perbandingan 1:6 air hingga lembut. 2. Rebus ayam, suwir halus. 3. Taruh bubur di mangkuk, tambahkan ayam suwir, daun bawang, bawang goreng. 4. Siram kecap manis, sajikan dengan telur setengah matang dan kerupuk.",
	}

	menu.Lunch = models.MealPlan{
		MealType:      "lunch",
		Title:         "🍱 Nasi Liwet Komplit Sehat",
		Foods:         []string{"Nasi liwet", "Ayam goreng bumbu kuning", "Tempe mendoan", "Lalapan & sambal matah", "Sayur asem"},
		Ingredients:   []string{"Beras 200g", "Santan encer 200ml", "Daun salam 2 lembar", "Serai 1 batang", "Ayam 1 potong paha", "Tempe 100g", "Timun, kemangi, kol", "Tomat, bawang merah, cabai", "Sayur asem (kangkung, jagung, labu)"},
		Description:   "Menu makan siang Indonesia yang lengkap dan seimbang",
		Calories:      "~650 kkal",
		EstimatedCost: "Rp 25.000 - 35.000",
		Recipe:        "1. Masak nasi dengan santan encer, daun salam, serai. 2. Goreng ayam dengan bumbu kuning (kunyit, bawang). 3. Balut tempe tipis dengan tepung berbumbu, goreng. 4. Iris tomat, bawang merah, cabai untuk sambal matah. 5. Rebus sayur asem dengan asam jawa.",
	}

	menu.Dinner = models.MealPlan{
		MealType:      "dinner",
		Title:         "🌙 Pepes Ikan & Sayur Bening",
		Foods:         []string{"Pepes ikan mas", "Nasi merah", "Sayur bening bayam", "Tahu bacem"},
		Ingredients:   []string{"Ikan mas 200g", "Bumbu pepes (kemangi, daun salam, lengkuas)", "Nasi merah 150g", "Bayam 100g", "Jagung manis 1/2 buah", "Tahu 100g", "Kecap manis, gula jawa"},
		Description:   "Makan malam sehat rendah lemak dengan protein ikan",
		Calories:      "~450 kkal",
		EstimatedCost: "Rp 30.000 - 40.000",
		Recipe:        "1. Balut ikan dengan bumbu pepes dalam daun pisang, kukus 20 menit. 2. Rebus sayur bening dari bayam dan jagung. 3. Bacem tahu dengan kecap manis dan gula jawa. 4. Sajikan dengan nasi merah hangat.",
	}

	menu.Snacks = []models.MealPlan{
		{
			MealType:      "snack",
			Title:         "🍌 Snack Pagi: Pisang Rebus & Kacang Hijau",
			Foods:         []string{"Pisang rebus", "Bubur kacang hijau"},
			Ingredients:   []string{"Pisang kepok 2 buah", "Kacang hijau 50g", "Gula aren 2 sdm", "Santan encer 100ml"},
			Description:   "Camilan tradisional tinggi energi dan protein nabati",
			Calories:      "~200 kkal",
			EstimatedCost: "Rp 8.000 - 12.000",
			Recipe:        "Rebus pisang 10 menit. Masak kacang hijau dengan gula aren dan santan hingga lembut.",
		},
		{
			MealType:      "snack",
			Title:         "🥜 Snack Siang: Gado-Gado Mini",
			Foods:         []string{"Gado-gado porsi kecil"},
			Ingredients:   []string{"Tahu goreng 50g", "Tempe goreng 50g", "Kol, tauge, kacang panjang rebus", "Bumbu kacang 3 sdm", "Kerupuk"},
			Description:   "Makanan bergizi dengan protein nabati dan sayuran",
			Calories:      "~180 kkal",
			EstimatedCost: "Rp 10.000 - 15.000",
			Recipe:        "Rebus sayuran, potong tahu tempe. Siram dengan bumbu kacang, tabur kerupuk.",
		},
		{
			MealType:      "snack",
			Title:         "🍠 Snack Sore: Kolak Pisang Ubi",
			Foods:         []string{"Kolak pisang ubi"},
			Ingredients:   []string{"Pisang raja 1 buah", "Ubi ungu 50g", "Gula aren 2 sdm", "Santan 100ml", "Daun pandan"},
			Description:   "Dessert tradisional yang menghangatkan tubuh",
			Calories:      "~180 kkal",
			EstimatedCost: "Rp 8.000 - 12.000",
			Recipe:        "Rebus santan dengan gula aren dan pandan. Masukkan ubi dan pisang, masak hingga empuk.",
		},
		{
			MealType:      "snack",
			Title:         "🥤 Snack Malam: Es Buah Segar",
			Foods:         []string{"Es buah campur"},
			Ingredients:   []string{"Semangka", "Melon", "Pepaya", "Nata de coco", "Susu kental manis 1 sdm", "Es batu"},
			Description:   "Camilan segar tinggi vitamin dan serat",
			Calories:      "~120 kkal",
			EstimatedCost: "Rp 10.000 - 15.000",
			Recipe:        "Potong buah-buahan, campur dengan nata de coco, siram susu, tambahkan es.",
		},
	}

	menu.TotalCalories = "~1850 kkal"
	menu.TotalEstimatedCost = "Rp 89.000 - 139.000"
	menu.HealthTip = "💡 Tips: Makan dalam porsi seimbang dengan prinsip 'Isi Piringku' - 1/3 karbohidrat, 1/3 sayuran, 1/3 protein."
	
	// Expanded drinks and fruits with Indonesian options
	menu.Drinks = []string{
		"💧 Air putih 8-10 gelas/hari", 
		"🍵 Teh hijau/jahe tanpa gula", 
		"🥛 Susu segar/kedelai", 
		"🥥 Air kelapa muda",
		"🍋 Jus lemon-madu hangat",
		"🫖 Wedang uwuh (herbal Jawa)",
		"🌿 Jamu kunyit asam",
		"🍊 Jus jeruk segar tanpa gula",
	}
	menu.Fruits = []string{
		"🍌 Pisang - sumber kalium & energi",
		"🍎 Apel - tinggi serat pectin", 
		"🥭 Pepaya - enzim pencernaan",
		"🍊 Jeruk - vitamin C",
		"🍉 Semangka - hidrasi",
		"🫐 Blueberry - antioksidan",
		"🥑 Alpukat - lemak sehat & serat",
		"🍇 Anggur merah - resveratrol",
		"🥝 Kiwi - vitamin C & serat",
		"🍐 Pir - indeks glikemik rendah",
		"🍓 Stroberi - antioksidan",
		"🥥 Kelapa muda - elektrolit alami",
	}
	menu.AvoidDrinks = []string{
		"🚫 Minuman bersoda & energi", 
		"🚫 Alkohol berlebihan", 
		"🚫 Kopi >3 cangkir/hari",
		"🚫 Teh manis berlebihan",
		"🚫 Jus kemasan dengan pemanis",
		"🚫 Susu full cream (jika kolesterol tinggi)",
	}
	menu.AvoidFruits = []string{
		"🚫 Buah kalengan dengan sirup gula",
		"🚫 Durian berlebihan (tinggi kalori)",
		"🚫 Nangka berlebihan (tinggi gula)",
	}

	// Alternative Breakfast Options
	menu.BreakfastAlt = []models.MealPlan{
		{
			MealType:      "breakfast",
			Title:         "🥣 Oatmeal Pisang Madu",
			Foods:         []string{"Oatmeal dengan pisang", "Telur rebus", "Teh hijau"},
			Ingredients:   []string{"Oatmeal 4 sdm", "Pisang 1 buah", "Madu 1 sdm", "Susu 200ml", "Telur 2 butir"},
			Description:   "Sarapan tinggi serat untuk energi stabil",
			Calories:      "~400 kkal",
			EstimatedCost: "Rp 15.000 - 20.000",
			Recipe:        "Masak oatmeal dengan susu, tambahkan pisang dan madu. Rebus telur 10 menit.",
		},
		{
			MealType:      "breakfast",
			Title:         "🍳 Nasi Uduk Betawi",
			Foods:         []string{"Nasi uduk", "Telur balado", "Tempe orek", "Bihun goreng"},
			Ingredients:   []string{"Beras 150g", "Santan 100ml", "Serai", "Daun salam", "Telur 2 butir", "Tempe 50g"},
			Description:   "Sarapan khas Betawi yang mengenyangkan",
			Calories:      "~500 kkal",
			EstimatedCost: "Rp 12.000 - 18.000",
			Recipe:        "Masak nasi dengan santan dan rempah. Goreng telur, buat balado. Orek tempe manis.",
		},
		{
			MealType:      "breakfast",
			Title:         "🥪 Roti Bakar Telur Keju",
			Foods:         []string{"Roti bakar", "Telur orak-arik", "Keju slice", "Jus jeruk"},
			Ingredients:   []string{"Roti tawar 2 lembar", "Telur 2 butir", "Keju 1 slice", "Mentega", "Jeruk 2 buah"},
			Description:   "Sarapan praktis tinggi protein",
			Calories:      "~450 kkal",
			EstimatedCost: "Rp 15.000 - 22.000",
			Recipe:        "Panggang roti dengan mentega. Orak-arik telur, taruh keju di atas roti panas.",
		},
		{
			MealType:      "breakfast",
			Title:         "🍜 Lontong Sayur",
			Foods:         []string{"Lontong", "Sayur labu santan", "Telur rebus", "Kerupuk"},
			Ingredients:   []string{"Lontong 2 potong", "Labu siam 100g", "Tahu 50g", "Santan 200ml", "Telur 1 butir"},
			Description:   "Sarapan tradisional Jawa yang hangat",
			Calories:      "~480 kkal",
			EstimatedCost: "Rp 10.000 - 15.000",
			Recipe:        "Masak sayur labu dengan santan dan bumbu. Potong lontong, sajikan dengan telur rebus.",
		},
		{
			MealType:      "breakfast",
			Title:         "🥞 Pancake Pisang Oat",
			Foods:         []string{"Pancake pisang oat", "Madu", "Buah segar"},
			Ingredients:   []string{"Oat 50g", "Pisang 2 buah", "Telur 1 butir", "Susu 100ml", "Madu 2 sdm"},
			Description:   "Sarapan sehat tanpa tepung terigu",
			Calories:      "~380 kkal",
			EstimatedCost: "Rp 12.000 - 18.000",
			Recipe:        "Blender pisang, oat, telur, susu. Panggang seperti pancake. Sajikan dengan madu.",
		},
	}

	// Alternative Lunch Options
	menu.LunchAlt = []models.MealPlan{
		{
			MealType:      "lunch",
			Title:         "🍛 Nasi Padang Sehat",
			Foods:         []string{"Nasi putih", "Rendang daging", "Sayur daun singkong", "Sambal hijau"},
			Ingredients:   []string{"Nasi 200g", "Daging sapi 100g", "Daun singkong 100g", "Cabai hijau", "Bumbu rendang"},
			Description:   "Menu Padang dengan porsi protein tinggi",
			Calories:      "~650 kkal",
			EstimatedCost: "Rp 25.000 - 35.000",
			Recipe:        "Masak rendang dengan bumbu tradisional. Rebus daun singkong dengan santan.",
		},
		{
			MealType:      "lunch",
			Title:         "🍲 Soto Ayam Lamongan",
			Foods:         []string{"Soto ayam kuah kuning", "Nasi", "Telur", "Koya", "Sambal"},
			Ingredients:   []string{"Ayam 150g", "Kunyit, lengkuas, serai", "Nasi 150g", "Telur 1 butir", "Koya (kerupuk bubuk)"},
			Description:   "Sup hangat khas Jawa Timur yang menyegarkan",
			Calories:      "~500 kkal",
			EstimatedCost: "Rp 18.000 - 25.000",
			Recipe:        "Rebus ayam dengan bumbu kuning. Suwir ayam, sajikan dengan kuah kuning, nasi, koya.",
		},
		{
			MealType:      "lunch",
			Title:         "🥗 Gado-Gado Jakarta",
			Foods:         []string{"Gado-gado sayuran", "Lontong", "Telur", "Kerupuk emping"},
			Ingredients:   []string{"Kol, tauge, kacang panjang, kentang", "Lontong 2 potong", "Telur 1 butir", "Bumbu kacang 100g"},
			Description:   "Salad Indonesia tinggi protein nabati",
			Calories:      "~550 kkal",
			EstimatedCost: "Rp 15.000 - 22.000",
			Recipe:        "Rebus sayuran dan kentang. Potong lontong. Siram dengan bumbu kacang kental.",
		},
		{
			MealType:      "lunch",
			Title:         "🍱 Ayam Geprek Sambal Matah",
			Foods:         []string{"Ayam geprek", "Nasi", "Sambal matah", "Lalapan"},
			Ingredients:   []string{"Dada ayam 150g", "Tepung bumbu", "Nasi 200g", "Bawang merah, cabai, jeruk limau"},
			Description:   "Ayam crispy dengan sambal segar Bali",
			Calories:      "~600 kkal",
			EstimatedCost: "Rp 18.000 - 25.000",
			Recipe:        "Goreng ayam crispy, geprek. Buat sambal matah dari bawang iris dan cabai rawit.",
		},
		{
			MealType:      "lunch",
			Title:         "🐟 Ikan Bakar Jimbaran",
			Foods:         []string{"Ikan kakap bakar", "Nasi", "Sambal plecing", "Sayur kangkung"},
			Ingredients:   []string{"Ikan kakap 250g", "Bumbu bakar", "Nasi 200g", "Kangkung 100g", "Sambal tomat"},
			Description:   "Ikan bakar segar khas Bali dengan sambal pedas",
			Calories:      "~520 kkal",
			EstimatedCost: "Rp 30.000 - 45.000",
			Recipe:        "Bakar ikan dengan bumbu bali. Rebus kangkung, sajikan dengan sambal plecing.",
		},
	}

	// Alternative Dinner Options
	menu.DinnerAlt = []models.MealPlan{
		{
			MealType:      "dinner",
			Title:         "🍜 Mie Ayam Bakso",
			Foods:         []string{"Mie ayam", "Bakso sapi", "Pangsit goreng", "Sawi hijau"},
			Ingredients:   []string{"Mie telur 150g", "Ayam cincang 100g", "Bakso 4 butir", "Sawi 50g", "Kaldu ayam"},
			Description:   "Comfort food Indonesia yang menghangatkan",
			Calories:      "~480 kkal",
			EstimatedCost: "Rp 15.000 - 22.000",
			Recipe:        "Rebus mie, tumis ayam cincang dengan kecap. Panaskan bakso dalam kaldu.",
		},
		{
			MealType:      "dinner",
			Title:         "🥣 Sup Iga Sapi",
			Foods:         []string{"Sup iga sapi", "Nasi hangat", "Sambal kecap", "Emping"},
			Ingredients:   []string{"Iga sapi 200g", "Kentang, wortel, tomat", "Seledri, daun bawang", "Nasi 150g"},
			Description:   "Sup bening bergizi tinggi kolagen",
			Calories:      "~550 kkal",
			EstimatedCost: "Rp 35.000 - 50.000",
			Recipe:        "Rebus iga hingga empuk 2 jam. Tambahkan sayuran, sajikan dengan sambal kecap.",
		},
		{
			MealType:      "dinner",
			Title:         "🍛 Nasi Goreng Kampung",
			Foods:         []string{"Nasi goreng kampung", "Telur mata sapi", "Kerupuk", "Acar"},
			Ingredients:   []string{"Nasi 200g", "Telur 2 butir", "Cabai rawit", "Bawang merah, putih", "Kecap manis"},
			Description:   "Nasi goreng pedas dengan bumbu sederhana",
			Calories:      "~520 kkal",
			EstimatedCost: "Rp 12.000 - 18.000",
			Recipe:        "Tumis bumbu, masukkan nasi dan cabai. Goreng telur mata sapi di atasnya.",
		},
		{
			MealType:      "dinner",
			Title:         "🥗 Salad Ayam Mediterranean",
			Foods:         []string{"Salad sayuran", "Dada ayam panggang", "Kentang wedges"},
			Ingredients:   []string{"Selada, tomat, mentimun", "Dada ayam 150g", "Kentang 100g", "Minyak zaitun, lemon"},
			Description:   "Makan malam ringan ala barat",
			Calories:      "~400 kkal",
			EstimatedCost: "Rp 25.000 - 35.000",
			Recipe:        "Panggang ayam dengan herbs. Buat salad dengan dressing olive oil lemon.",
		},
		{
			MealType:      "dinner",
			Title:         "🦐 Capcay Seafood",
			Foods:         []string{"Capcay seafood", "Nasi putih", "Acar kuning"},
			Ingredients:   []string{"Udang 50g", "Bakso ikan 50g", "Sayuran campur 150g", "Saus tiram", "Nasi 150g"},
			Description:   "Tumis sayuran Chinese-Indonesian style",
			Calories:      "~450 kkal",
			EstimatedCost: "Rp 22.000 - 32.000",
			Recipe:        "Tumis seafood dengan sayuran, tambahkan saus tiram. Sajikan dengan nasi.",
		},
	}

	// Customize based on BMI
	switch bmiCategory {
	case "Underweight":
		menu.HealthTip = "💪 Fokus menambah asupan kalori sehat. Tambahkan alpukat, kacang-kacangan, dan susu full cream."
		menu.Breakfast = models.MealPlan{
			MealType:      "breakfast",
			Title:         "🌅 Sarapan Tinggi Kalori",
			Foods:         []string{"Nasi goreng telur", "Susu full cream", "Pisang"},
			Ingredients:   []string{"Nasi 200g", "Telur 2 butir", "Minyak 2 sdm", "Bawang putih 2 siung", "Kecap manis", "Susu full cream 250ml", "Pisang 1 buah"},
			Description:   "Sarapan padat kalori untuk menambah berat badan sehat",
			Calories:      "~600 kkal",
			EstimatedCost: "Rp 12.000 - 18.000",
			Recipe:        "1. Tumis bawang putih hingga harum. 2. Masukkan nasi, aduk rata. 3. Buat lubang, masukkan telur, orak-arik. 4. Tambahkan kecap manis, aduk rata. 5. Sajikan dengan susu dan pisang.",
		}
		menu.Lunch = models.MealPlan{
			MealType:      "lunch",
			Title:         "🍱 Makan Siang Berenergi",
			Foods:         []string{"Nasi putih porsi besar", "Ayam goreng", "Tempe goreng", "Sayur santan"},
			Ingredients:   []string{"Nasi 250g", "Ayam 1 potong paha", "Tempe 100g", "Sayur nangka muda", "Santan 100ml", "Bumbu lengkap"},
			Description:   "Makan siang tinggi karbohidrat dan protein",
			Calories:      "~750 kkal",
			EstimatedCost: "Rp 20.000 - 30.000",
			Recipe:        "1. Goreng ayam dengan bumbu kuning. 2. Goreng tempe tipis-tipis. 3. Masak sayur nangka dengan santan. 4. Sajikan dengan nasi hangat porsi besar.",
		}
		menu.Dinner = models.MealPlan{
			MealType:      "dinner",
			Title:         "🌙 Makan Malam Bergizi",
			Foods:         []string{"Nasi tim ayam", "Sup daging", "Alpukat jus"},
			Ingredients:   []string{"Beras 150g", "Ayam cincang 100g", "Daging sapi 100g", "Wortel, kentang", "Alpukat 1 buah", "Susu kental manis"},
			Description:   "Makan malam mudah dicerna tapi tinggi kalori",
			Calories:      "~650 kkal",
			EstimatedCost: "Rp 35.000 - 45.000",
			Recipe:        "1. Masak nasi tim dengan ayam cincang dan kaldu. 2. Rebus daging dengan sayuran untuk sup. 3. Blender alpukat dengan susu kental manis.",
		}
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🥜 Snack Pagi", Foods: []string{"Roti selai kacang", "Susu cokelat"}, Ingredients: []string{"Roti 2 lembar", "Selai kacang 2 sdm", "Susu cokelat 200ml"}, Calories: "~350 kkal", EstimatedCost: "Rp 10.000 - 15.000", Recipe: "Oleskan selai kacang pada roti, sajikan dengan susu cokelat."},
			{MealType: "snack", Title: "🍌 Snack Sore", Foods: []string{"Pisang goreng", "Teh manis"}, Ingredients: []string{"Pisang 2 buah", "Tepung 50g", "Minyak goreng"}, Calories: "~300 kkal", EstimatedCost: "Rp 8.000 - 12.000", Recipe: "Balut pisang dengan tepung, goreng hingga kecokelatan."},
		}
		menu.TotalCalories = "~2650 kkal"
		menu.TotalEstimatedCost = "Rp 85.000 - 120.000"
	case "Overweight", "Obese":
		menu.HealthTip = "🥗 Fokus pada makanan rendah kalori tapi mengenyangkan. Perbanyak sayuran dan protein tanpa lemak."
		menu.Breakfast = models.MealPlan{
			MealType:      "breakfast",
			Title:         "🌅 Sarapan Rendah Kalori",
			Foods:         []string{"Telur rebus", "Salad sayur", "Teh hijau"},
			Ingredients:   []string{"Telur 2 butir", "Selada 1 mangkuk", "Tomat 1 buah", "Mentimun 1/2 buah", "Perasan lemon", "Teh hijau 1 kantong"},
			Description:   "Sarapan tinggi protein rendah karbohidrat",
			Calories:      "~200 kkal",
			EstimatedCost: "Rp 10.000 - 15.000",
			Recipe:        "1. Rebus telur 10 menit. 2. Potong sayuran untuk salad. 3. Siram dengan perasan lemon. 4. Seduh teh hijau tanpa gula.",
		}
		menu.Lunch = models.MealPlan{
			MealType:      "lunch",
			Title:         "🍱 Makan Siang Diet",
			Foods:         []string{"Salad sayuran besar", "Dada ayam panggang", "Sup sayur tanpa santan"},
			Ingredients:   []string{"Dada ayam 150g", "Sayuran campur 200g", "Wortel, brokoli, bayam", "Minyak zaitun 1 sdt", "Lemon, garam, lada"},
			Description:   "Rendah karbohidrat, tinggi protein dan serat",
			Calories:      "~350 kkal",
			EstimatedCost: "Rp 20.000 - 28.000",
			Recipe:        "1. Panggang dada ayam tanpa kulit. 2. Rebus sayuran untuk sup tanpa santan. 3. Buat salad dengan dressing minyak zaitun + lemon.",
		}
		menu.Dinner = models.MealPlan{
			MealType:      "dinner",
			Title:         "🌙 Makan Malam Super Ringan",
			Foods:         []string{"Sup sayuran", "Tahu kukus", "Sayuran kukus"},
			Ingredients:   []string{"Tahu 150g", "Wortel, brokoli, kembang kol 150g", "Bawang putih", "Garam & lada secukupnya"},
			Description:   "Makan malam sangat ringan untuk tidur nyenyak",
			Calories:      "~200 kkal",
			EstimatedCost: "Rp 12.000 - 18.000",
			Recipe:        "1. Kukus tahu 10 menit. 2. Kukus sayuran sampai empuk. 3. Sajikan dengan sedikit kecap asin rendah garam.",
		}
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🥒 Snack Pagi", Foods: []string{"Mentimun", "Wortel"}, Ingredients: []string{"Mentimun 1 buah", "Wortel 1 buah"}, Calories: "~50 kkal", EstimatedCost: "Rp 3.000 - 5.000", Description: "Sayuran segar tanpa saus", Recipe: "Cuci dan potong, makan langsung."},
			{MealType: "snack", Title: "🍎 Snack Sore", Foods: []string{"Apel", "Air putih"}, Ingredients: []string{"Apel 1 buah"}, Calories: "~80 kkal", EstimatedCost: "Rp 5.000 - 8.000", Description: "Buah rendah kalori", Recipe: "Cuci apel, makan dengan kulitnya untuk serat maksimal."},
		}
		menu.TotalCalories = "~880 kkal"
		menu.TotalEstimatedCost = "Rp 50.000 - 74.000"
	}

	// Customize based on symptoms
	if symptomNames["Diabetes"] || symptomNames["Gula Darah Tinggi"] {
		menu.HealthTip = "🩸 Pilih makanan dengan indeks glikemik rendah. Hindari gula dan makanan olahan."
		menu.Breakfast = models.MealPlan{
			MealType:    "breakfast",
			Title:       "🌅 Sarapan Diabetes-Friendly",
			Foods:       []string{"Oatmeal tanpa gula", "Telur dadar sayuran", "Alpukat"},
			Ingredients: []string{"Oatmeal 4 sdm", "Telur 2 butir", "Bayam 50g", "Alpukat 1/2 buah", "Susu almond tanpa gula 100ml"},
			Description: "Sarapan rendah gula, tinggi serat dan protein",
			Calories:    "~350 kkal",
			Recipe:      "1. Masak oatmeal dengan air/susu almond tanpa gula. 2. Dadar telur dengan bayam cincang. 3. Iris alpukat sebagai pendamping.",
		}
		menu.Lunch = models.MealPlan{
			MealType:    "lunch",
			Title:       "🍱 Makan Siang Gula Darah Stabil",
			Foods:       []string{"Nasi merah porsi kecil", "Ikan panggang", "Tumis sayuran", "Tahu kukus"},
			Ingredients: []string{"Nasi merah 100g", "Ikan kakap 150g", "Brokoli, wortel 100g", "Tahu 100g", "Bawang putih, jahe"},
			Description: "Rendah karbohidrat dengan protein tinggi",
			Calories:    "~400 kkal",
			Recipe:      "1. Panggang ikan dengan bumbu jahe. 2. Tumis sayuran dengan sedikit minyak. 3. Kukus tahu sajikan dengan kecap rendah gula.",
		}
		menu.Dinner = models.MealPlan{
			MealType:    "dinner",
			Title:       "🌙 Makan Malam Ringan Diabetesi",
			Foods:       []string{"Salad sayuran", "Dada ayam panggang", "Sup sayur bening"},
			Ingredients: []string{"Selada, tomat, mentimun", "Dada ayam 100g", "Wortel, kentang 50g", "Minyak zaitun 1 sdt"},
			Description: "Makan malam yang tidak menaikkan gula darah",
			Calories:    "~300 kkal",
			Recipe:      "1. Panggang dada ayam tanpa kulit. 2. Buat salad dengan dressing minyak zaitun. 3. Rebus sup sayur tanpa gula.",
		}
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🥒 Snack Pagi", Foods: []string{"Mentimun", "Kacang almond 10 butir"}, Ingredients: []string{"Mentimun 1 buah", "Almond 15g"}, Calories: "~100 kkal", Recipe: "Cuci mentimun, potong-potong, sajikan dengan almond."},
			{MealType: "snack", Title: "🍐 Snack Sore", Foods: []string{"Pir", "Keju rendah lemak"}, Ingredients: []string{"Pir 1 buah", "Keju slice 1 lembar"}, Calories: "~120 kkal", Recipe: "Iris pir, sajikan dengan keju. Pir memiliki IG rendah."},
		}
		menu.Drinks = []string{"☕ Air putih 8-10 gelas/hari", "🍵 Teh hijau tanpa gula", "☕ Kopi hitam tanpa gula (maks 2 cangkir)", "🥛 Susu almond tanpa gula", "🍋 Air lemon hangat", "🥒 Infused water mentimun"}
		menu.Fruits = []string{"🍐 Pir - IG rendah", "🍎 Apel hijau - serat tinggi", "🫐 Blueberry - antioksidan", "🍓 Stroberi - IG rendah", "🥑 Alpukat - lemak sehat", "🍒 Ceri - antiinflamasi"}
		menu.AvoidDrinks = []string{"🚫 Jus kemasan dengan gula", "🚫 Minuman bersoda", "🚫 Teh manis", "🚫 Kopi dengan gula/krimer", "🚫 Susu full cream", "🚫 Smoothie dengan es krim"}
		menu.AvoidFruits = []string{"🚫 Semangka - IG tinggi", "🚫 Nanas matang", "🚫 Mangga matang berlebihan", "🚫 Durian", "🚫 Buah kalengan sirup", "🚫 Kurma berlebihan"}
		menu.TotalCalories = "~1270 kkal"
	}

	if symptomNames["Maag"] || symptomNames["Gangguan Pencernaan"] {
		menu.HealthTip = "🍵 Makan dalam porsi kecil tapi sering (5-6x sehari). Hindari makanan pedas, asam, dan berminyak."
		menu.Breakfast = models.MealPlan{
			MealType:    "breakfast",
			Title:       "🌅 Sarapan Ramah Lambung",
			Foods:       []string{"Bubur ayam lembut", "Pisang matang", "Teh hangat tawar"},
			Ingredients: []string{"Beras 50g", "Ayam suwir 50g", "Daun bawang", "Pisang 1 buah", "Air 500ml"},
			Description: "Sarapan lembut mudah dicerna untuk lambung sensitif",
			Calories:    "~300 kkal",
			Recipe:      "1. Masak bubur dari beras dengan air hingga lembut. 2. Tambahkan ayam suwir dan daun bawang. 3. Sajikan dengan pisang matang.",
		}
		menu.Lunch = models.MealPlan{
			MealType:    "lunch",
			Title:       "🍱 Makan Siang Anti Maag",
			Foods:       []string{"Nasi putih lembek", "Ayam rebus", "Sayur bening bayam", "Tahu kukus"},
			Ingredients: []string{"Nasi 150g", "Ayam 100g", "Bayam 50g", "Tahu 100g", "Air kaldu"},
			Description: "Makanan rebus dan kukus yang tidak merangsang lambung",
			Calories:    "~400 kkal",
			Recipe:      "1. Rebus ayam hingga empuk tanpa bumbu pedas. 2. Buat sayur bening dari bayam. 3. Kukus tahu, sajikan dengan nasi lembek.",
		}
		menu.Dinner = models.MealPlan{
			MealType:    "dinner",
			Title:       "🌙 Makan Malam Gentle",
			Foods:       []string{"Sup kentang wortel", "Roti tawar", "Pisang"},
			Ingredients: []string{"Kentang 100g", "Wortel 50g", "Roti tawar 2 lembar", "Pisang 1 buah"},
			Description: "Makan malam ringan 3 jam sebelum tidur",
			Calories:    "~300 kkal",
			Recipe:      "1. Rebus kentang dan wortel hingga lembut. 2. Haluskan sedikit untuk tekstur soup. 3. Sajikan dengan roti dan pisang.",
		}
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍌 Snack Pagi (10:00)", Foods: []string{"Pisang", "Biskuit tawar"}, Ingredients: []string{"Pisang 1 buah", "Biskuit 2 keping"}, Calories: "~150 kkal", Recipe: "Makan biskuit dengan pisang untuk menetralkan asam lambung."},
			{MealType: "snack", Title: "🥛 Snack Sore (15:00)", Foods: []string{"Susu hangat", "Roti panggang"}, Ingredients: []string{"Susu rendah lemak 200ml", "Roti 1 lembar"}, Calories: "~180 kkal", Recipe: "Hangatkan susu, panggang roti tanpa mentega."},
		}
		menu.Drinks = []string{"🥛 Susu hangat", "🍵 Teh chamomile", "☕ Air putih hangat", "🥥 Air kelapa muda", "🍯 Air madu hangat (pagi)", "🥒 Jus lidah buaya"}
		menu.Fruits = []string{"🍌 Pisang matang - menetralkan asam", "🍐 Pir - serat lembut", "🍎 Apel tanpa kulit - mudah dicerna", "🍈 Melon - menenangkan lambung", "🥭 Pepaya - enzim pencernaan", "🥑 Alpukat - melindungi lambung"}
		menu.AvoidDrinks = []string{"🚫 Kopi", "🚫 Teh terlalu pekat", "🚫 Minuman bersoda", "🚫 Alkohol", "🚫 Jus jeruk/asam", "🚫 Minuman terlalu dingin"}
		menu.AvoidFruits = []string{"🚫 Jeruk - terlalu asam", "🚫 Lemon langsung", "🚫 Nanas - tinggi asam", "🚫 Tomat mentah", "🚫 Mangga muda"}
		menu.TotalCalories = "~1330 kkal"
	}

	if symptomNames["Tekanan Darah Tinggi"] || symptomNames["Hipertensi"] {
		menu.HealthTip = "❤️ Kurangi garam! Maksimal 1 sendok teh per hari. Perbanyak kalium dari pisang."
		menu.Breakfast.Foods = []string{"Oatmeal", "Pisang", "Yogurt rendah lemak"}
		menu.Breakfast.Recipe = "Masak oatmeal tanpa garam, tambahkan pisang iris dan yogurt."
		menu.Lunch.Foods = []string{"Nasi merah", "Ikan kukus lemon", "Sayur bayam", "Jus bit"}
		menu.Lunch.Recipe = "Kukus ikan dengan perasan lemon, rebus bayam tanpa garam. Blender bit dengan air."
		menu.Dinner.Foods = []string{"Salad sayuran tanpa garam", "Tahu panggang", "Kentang rebus"}
		menu.Dinner.Recipe = "Panggang tahu dengan bumbu herbal tanpa garam. Dressing salad: minyak zaitun + lemon."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍌 Snack Pagi", Foods: []string{"Pisang", "Air kelapa"}, Calories: "~120 kkal", Description: "Tinggi kalium untuk tekanan darah"},
		}
		menu.Drinks = []string{"☕ Air putih 10+ gelas/hari", "🥥 Air kelapa - tinggi kalium", "🍵 Teh hibiscus - menurunkan TD", "🥤 Jus bit - vasodilator alami", "🍶 Susu skim", "🫖 Teh hijau tanpa gula"}
		menu.Fruits = []string{"🍌 Pisang - tinggi kalium", "🍊 Jeruk - kalium+vitamin C", "🫐 Blueberry - antioksidan", "🥝 Kiwi - menurunkan TD", "🍉 Semangka - citrulline", "🍇 Anggur merah"}
		menu.AvoidDrinks = []string{"🚫 Kopi >2 cangkir", "🚫 Alkohol", "🚫 Minuman energi", "🚫 Minuman bersoda", "🚫 Teh manis"}
		menu.AvoidFruits = []string{"🚫 Buah kalengan tinggi natrium", "🚫 Acar buah"}
	}

	// KOLESTEROL TINGGI
	if symptomNames["Kolesterol Tinggi"] || symptomNames["Kolesterol"] {
		menu.HealthTip = "Hindari lemak jenuh dan trans. Perbanyak serat larut dan omega-3."
		menu.Breakfast.Foods = []string{"Oatmeal dengan chia seed", "Apel", "Teh hijau"}
		menu.Breakfast.Recipe = "Masak oatmeal, tambahkan 1 sdm chia seed dan apel potong."
		menu.Lunch.Foods = []string{"Nasi merah", "Salmon panggang", "Tumis brokoli", "Sup kacang merah"}
		menu.Lunch.Recipe = "Panggang salmon dengan minyak zaitun. Tumis brokoli dengan bawang putih tanpa minyak berlebihan."
		menu.Dinner.Foods = []string{"Salad alpukat", "Dada ayam panggang tanpa kulit", "Sayuran kukus"}
		menu.Dinner.Recipe = "Buang kulit ayam, panggang dengan herbs. Sajikan dengan salad alpukat."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🥜 Snack Sehat", Foods: []string{"Kacang walnut 10 butir", "Apel"}, Calories: "~150 kkal", Description: "Omega-3 untuk menurunkan kolesterol"},
		}
		menu.TotalCalories = "~1400 kkal"
	}

	// ANEMIA
	if symptomNames["Anemia"] || symptomNames["Kurang Darah"] || symptomNames["Pucat"] {
		menu.HealthTip = "Konsumsi makanan tinggi zat besi + vitamin C. Hindari teh/kopi saat makan utama."
		menu.Breakfast.Foods = []string{"Telur dadar bayam", "Roti gandum", "Jus jeruk"}
		menu.Breakfast.Recipe = "Dadar telur dengan bayam cincang. Minum jus jeruk untuk bantu penyerapan zat besi."
		menu.Lunch.Foods = []string{"Nasi merah", "Daging sapi tumis paprika", "Brokoli", "Sup kacang merah"}
		menu.Lunch.Recipe = "Tumis daging sapi dengan paprika merah (vitamin C). Rebus kacang merah untuk sup."
		menu.Dinner.Foods = []string{"Hati ayam goreng sedikit minyak", "Tempe bacem", "Tumis kangkung", "Nasi"}
		menu.Dinner.Recipe = "Goreng hati dengan sedikit minyak. Bacem tempe dengan bumbu manis."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍫 Snack Zat Besi", Foods: []string{"Kurma 5 buah", "Kismis"}, Calories: "~120 kkal", Description: "Tinggi zat besi alami"},
		}
	}

	// ASAM URAT
	if symptomNames["Asam Urat"] || symptomNames["Gout"] {
		menu.HealthTip = "Hindari makanan tinggi purin. Minum air putih minimal 10 gelas sehari."
		menu.Breakfast.Foods = []string{"Nasi dengan telur mata sapi", "Sayur bening labu", "Teh herbal"}
		menu.Breakfast.Recipe = "Goreng telur dengan sedikit minyak. Buat sayur bening dari labu siam."
		menu.Lunch.Foods = []string{"Nasi", "Tahu bacem", "Tempe goreng", "Sayur lodeh tanpa kacang"}
		menu.Lunch.Recipe = "Bacem tahu dengan bumbu manis. Lodeh dari labu dan wortel tanpa kacang."
		menu.Dinner.Foods = []string{"Sup sayuran (wortel, kentang, labu)", "Telur rebus", "Nasi sedikit"}
		menu.Dinner.Recipe = "Rebus sayuran dengan kaldu non-daging. Tambahkan telur rebus."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍒 Snack Anti Asam Urat", Foods: []string{"Buah ceri", "Air putih"}, Calories: "~80 kkal", Description: "Ceri membantu menurunkan asam urat"},
		}
	}

	// STRES & KECEMASAN  
	if symptomNames["Stres"] || symptomNames["Cemas"] || symptomNames["Kecemasan"] {
		menu.HealthTip = "Makanan kaya magnesium dan omega-3 membantu menenangkan pikiran."
		menu.Breakfast.Foods = []string{"Oatmeal dengan pisang", "Teh chamomile", "Almond"}
		menu.Breakfast.Recipe = "Masak oatmeal, tambahkan pisang. Seduh teh chamomile hangat."
		menu.Lunch.Foods = []string{"Nasi merah", "Salmon panggang", "Salad bayam alpukat", "Air lemon"}
		menu.Lunch.Recipe = "Panggang salmon. Campurkan bayam segar dengan irisan alpukat."
		menu.Dinner.Foods = []string{"Sup ayam hangat", "Kentang tumbuk", "Sayuran kukus", "Teh herbal"}
		menu.Dinner.Recipe = "Buat sup ayam dengan wortel dan seledri. Tumbuk kentang dengan sedikit susu."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍫 Snack Mood Booster", Foods: []string{"Cokelat hitam 2 kotak", "Kacang almond"}, Calories: "~100 kkal", Description: "Cokelat hitam meningkatkan serotonin"},
		}
	}

	// GANGGUAN TIDUR
	if symptomNames["Gangguan Tidur"] || symptomNames["Insomnia"] || symptomNames["Sulit Tidur"] {
		menu.HealthTip = "Hindari kafein setelah jam 2 siang. Makan malam ringan 3 jam sebelum tidur."
		menu.Breakfast.Foods = []string{"Roti gandum", "Telur rebus", "Pisang", "Susu hangat"}
		menu.Breakfast.Recipe = "Sarapan yang mengenyangkan untuk energi stabil sepanjang hari."
		menu.Lunch.Foods = []string{"Nasi merah", "Ikan panggang", "Tumis sayuran", "Sup"}
		menu.Lunch.Recipe = "Makan siang normal, hindari kafein setelah makan siang."
		menu.Dinner.Foods = []string{"Oatmeal ringan", "Kiwi 2 buah", "Madu", "Susu hangat"}
		menu.Dinner.Recipe = "Makan malam ringan. Kiwi mengandung serotonin. Susu hangat dengan madu."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🥛 Snack Sebelum Tidur", Foods: []string{"Susu hangat", "Madu 1 sdt"}, Calories: "~100 kkal", Description: "2 jam sebelum tidur untuk kualitas tidur"},
		}
		menu.TotalCalories = "~1400 kkal"
	}

	// DEMAM & FLU
	if symptomNames["Demam"] || symptomNames["Flu"] || symptomNames["Pilek"] {
		menu.HealthTip = "Perbanyak cairan hangat dan istirahat. Vitamin C untuk daya tahan tubuh."
		menu.Breakfast.Foods = []string{"Bubur ayam hangat", "Teh jahe madu", "Jeruk"}
		menu.Breakfast.Recipe = "Bubur dengan ayam suwir dan bawang goreng. Seduh jahe dengan madu."
		menu.Lunch.Foods = []string{"Sup ayam hangat", "Nasi lembek", "Sayur bening"}
		menu.Lunch.Recipe = "Sup ayam dengan banyak kaldu. Nasi yang agak lembek mudah dicerna."
		menu.Dinner.Foods = []string{"Bubur", "Telur rebus", "Air jahe hangat"}
		menu.Dinner.Recipe = "Makan ringan saat malam. Perbanyak minum air hangat."
		menu.Snacks = []models.MealPlan{
			{MealType: "snack", Title: "🍊 Vitamin C Boost", Foods: []string{"Jeruk 2 buah", "Air hangat"}, Calories: "~100 kkal", Description: "Vitamin C untuk daya tahan tubuh"},
		}
		menu.TotalCalories = "~1200 kkal"
	}

	return menu
}

func generateFoodRecommendations(health models.HealthData, symptoms []models.Symptom) []models.FoodRecommendation {
	var recommendations []models.FoodRecommendation

	bmiCategory := models.GetBMICategory(health.BMI)

	// BMI-based recommendations
	switch bmiCategory {
	case "Underweight":
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "weight_gain",
			Title:       "Makanan untuk Menambah Berat Badan",
			Description: "Tingkatkan asupan kalori dengan makanan bergizi tinggi",
			Foods:       []string{"Alpukat", "Kacang-kacangan", "Susu full cream", "Nasi merah", "Daging tanpa lemak", "Telur", "Keju", "Yogurt"},
			Avoid:       []string{"Makanan cepat saji", "Minuman bersoda"},
			Reason:      "BMI Anda di bawah normal, perlu menambah asupan kalori sehat",
		})
	case "Overweight", "Obese":
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "weight_loss",
			Title:       "Makanan untuk Menurunkan Berat Badan",
			Description: "Fokus pada makanan rendah kalori dan tinggi serat",
			Foods:       []string{"Sayuran hijau", "Buah-buahan segar", "Ikan", "Dada ayam", "Oatmeal", "Quinoa", "Kacang almond"},
			Avoid:       []string{"Gorengan", "Makanan tinggi gula", "Minuman manis", "Fast food", "Makanan olahan"},
			Reason:      "BMI Anda di atas normal, perlu mengurangi asupan kalori",
		})
	case "Normal":
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "maintenance",
			Title:       "Pertahankan Pola Makan Sehat",
			Description: "Lanjutkan konsumsi makanan seimbang untuk gaya hidup sehat",
			Foods:       []string{"Sayuran beragam warna", "Protein seimbang", "Karbohidrat kompleks", "Buah segar", "Air putih cukup"},
			Avoid:       []string{"Makanan ultra-proses", "Gula berlebihan"},
			Reason:      "BMI Anda normal, pertahankan pola makan sehat",
		})
	}

	// Symptom-based recommendations
	symptomNames := make(map[string]bool)
	for _, s := range symptoms {
		symptomNames[s.SymptomName] = true
	}

	// DEMAM DAN FLU
	if symptomNames["Demam"] || symptomNames["Flu"] || symptomNames["Pilek"] || symptomNames["Batuk"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "fever_flu",
			Title:       "🤒 Makanan untuk Demam & Flu",
			Description: "Makanan yang membantu pemulihan dari demam dan flu",
			Foods:       []string{"Sup ayam hangat", "Air putih hangat", "Teh jahe madu", "Buah jeruk (Vitamin C)", "Pisang", "Bubur ayam", "Kaldu tulang", "Lemon hangat"},
			Avoid:       []string{"Makanan berminyak", "Gorengan", "Es/minuman dingin", "Makanan pedas", "Susu (dapat memperbanyak lendir)"},
			Reason:      "Anda mengalami demam/flu. Perbanyak cairan hangat, istirahat cukup, dan konsumsi makanan berkuah. Obat pereda panas ringan seperti parasetamol dapat membantu.",
		})
	}

	// TEKANAN DARAH TINGGI (Referensi: Alodokter)
	if symptomNames["Tekanan Darah Tinggi"] || symptomNames["Hipertensi"] || symptomNames["Pusing"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "blood_pressure",
			Title:       "❤️ Makanan untuk Hipertensi (Ref: Alodokter)",
			Description: "Diet DASH - Batasi garam, perbanyak kalium dan magnesium. Masak dengan cara dikukus atau direbus.",
			Foods:       []string{"Pisang (tinggi kalium)", "Sayuran hijau (bayam, brokoli, kangkung)", "Ikan omega-3 (salmon, tuna, sarden)", "Buah-buahan segar (jeruk, semangka, melon, pepaya)", "Yogurt rendah lemak", "Kacang-kacangan", "Oatmeal & biji-bijian utuh", "Daging tanpa lemak (direbus/dikukus)", "Bawang putih"},
			Avoid:       []string{"Garam berlebihan (maks 1 sdt/hari)", "Makanan kaleng & acar", "Daging olahan (sosis, kornet)", "Makanan cepat saji", "Keripik asin", "Mie instan", "Saus & kecap kemasan", "Alkohol", "Makanan tinggi lemak jenuh"},
			Reason:      "Sumber: Alodokter. Batasi garam maksimal 1 sendok teh/hari. Perbanyak buah pisang dan sayuran hijau untuk kalium. Olahraga teratur 30-45 menit, 3-5 kali/minggu. Jaga berat badan ideal dan berhenti merokok.",
		})
	}

	// STRES DAN KECEMASAN
	if symptomNames["Stres"] || symptomNames["Cemas"] || symptomNames["Kecemasan"] || symptomNames["Gelisah"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "stress_anxiety",
			Title:       "🧘 Makanan Pereda Stres & Kecemasan",
			Description: "Makanan yang membantu menenangkan pikiran dan mengurangi stres",
			Foods:       []string{"Cokelat hitam (70%+ kakao)", "Alpukat", "Teh chamomile", "Kacang almond", "Salmon (Omega-3)", "Blueberry", "Bayam", "Oatmeal", "Pisang", "Teh hijau"},
			Avoid:       []string{"Kafein berlebihan", "Alkohol", "Gula berlebihan", "Makanan olahan", "Minuman energi"},
			Reason:      "Anda mengalami stres/kecemasan. Selain makanan, disarankan untuk meditasi 10 menit, mendengarkan musik tenang, dan menulis jurnal perasaan.",
		})
	}

	// GANGGUAN TIDUR (INSOMNIA)
	if symptomNames["Gangguan Tidur"] || symptomNames["Insomnia"] || symptomNames["Sulit Tidur"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "sleep_disorder",
			Title:       "😴 Makanan untuk Kualitas Tidur",
			Description: "Makanan yang membantu meningkatkan kualitas tidur Anda",
			Foods:       []string{"Susu hangat", "Kacang almond", "Pisang", "Kiwi", "Ceri", "Teh chamomile", "Ikan salmon", "Nasi putih (porsi kecil)", "Oatmeal", "Madu"},
			Avoid:       []string{"Kafein (kopi, teh, cokelat) setelah jam 2 siang", "Alkohol", "Makanan pedas malam hari", "Makanan berat sebelum tidur", "Minuman energi"},
			Reason:      "Anda mengalami gangguan tidur. Hindari kafein setelah jam 2 siang, jaga suhu kamar sejuk (18-22°C), dan lakukan relaksasi ringan sebelum tidur seperti pernapasan 4-7-8.",
		})
	}

	// KELELAHAN DAN BURNOUT
	if symptomNames["Kelelahan Fisik"] || symptomNames["Kelelahan Emosional (Burnout)"] || symptomNames["Burnout"] || symptomNames["Lemas"] || symptomNames["Lesu"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "fatigue_burnout",
			Title:       "⚡ Makanan Penambah Energi",
			Description: "Nutrisi untuk melawan kelelahan dan memulihkan energi",
			Foods:       []string{"Bayam (zat besi)", "Pisang", "Kacang almond", "Telur", "Salmon", "Ubi jalar", "Cokelat hitam", "Quinoa", "Air kelapa", "Kurma", "Daging sapi tanpa lemak"},
			Avoid:       []string{"Gula berlebihan (spike energi)", "Kafein berlebihan", "Alkohol", "Makanan cepat saji", "Minuman bersoda"},
			Reason:      "Anda mengalami kelelahan/burnout. Selain pola makan, disarankan untuk istirahat cukup, lakukan aktivitas luar ruangan, dan luangkan waktu berkualitas bersama keluarga.",
		})
	}

	// MAAG DAN GANGGUAN PENCERNAAN
	if symptomNames["Maag"] || symptomNames["Gangguan Pencernaan"] || symptomNames["Mual"] || symptomNames["Perut Kembung"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "digestive",
			Title:       "🍵 Makanan untuk Pencernaan Sehat",
			Description: "Makanan yang mudah dicerna dan menenangkan lambung",
			Foods:       []string{"Pisang", "Nasi putih", "Roti tawar", "Ayam rebus", "Jahe hangat", "Pepaya", "Yogurt probiotik", "Oatmeal", "Kentang rebus"},
			Avoid:       []string{"Makanan pedas", "Kopi", "Alkohol", "Makanan berminyak", "Jeruk/makanan asam", "Cokelat", "Minuman bersoda"},
			Reason:      "Anda mengalami gangguan pencernaan. Makan dalam porsi kecil tapi sering, hindari makan terlalu cepat, dan jangan langsung berbaring setelah makan.",
		})
	}

	// KOLESTEROL TINGGI (Referensi: Alodokter)
	if symptomNames["Kolesterol Tinggi"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "cholesterol",
			Title:       "💚 Makanan untuk Kolesterol Tinggi (Ref: Alodokter)",
			Description: "Perbanyak serat larut dan omega-3. Hindari lemak jenuh dan trans. Gunakan minyak zaitun untuk memasak.",
			Foods:       []string{"Oatmeal & biji-bijian utuh (serat larut)", "Ikan omega-3 (salmon, makarel, sarden)", "Kacang walnut & almond", "Alpukat", "Minyak zaitun", "Sayuran hijau (bayam, brokoli)", "Buah tinggi serat (apel, pir, stroberi)", "Tahu & tempe", "Bawang putih", "Teh hijau"},
			Avoid:       []string{"Daging merah berlemak", "Jeroan (hati, otak, ampela)", "Kulit ayam & bebek", "Kuning telur berlebihan", "Makanan gorengan", "Mentega & margarin", "Santan kental", "Susu full cream & es krim", "Makanan olahan (lemak trans)"},
			Reason:      "Sumber: Alodokter. Serat larut dari oatmeal membantu menurunkan kolesterol. Omega-3 dari ikan meningkatkan kolesterol baik (HDL). Selalu baca label kemasan untuk menghindari lemak trans. Olahraga teratur dan jaga berat badan ideal.",
		})
	}

	// SAKIT KEPALA / MIGRAIN
	if symptomNames["Sakit Kepala"] || symptomNames["Migrain"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "headache",
			Title:       "🧠 Makanan Pereda Sakit Kepala",
			Description: "Makanan yang dapat membantu mengurangi sakit kepala",
			Foods:       []string{"Air putih (dehidrasi sering sebabkan sakit kepala)", "Magnesium (kacang, bayam)", "Jahe", "Ikan berlemak", "Semangka", "Kentang", "Pisang", "Kopi (secukupnya)"},
			Avoid:       []string{"Keju tua", "Makanan fermentasi berlebihan", "Alkohol (terutama red wine)", "MSG", "Pemanis buatan", "Cokelat berlebihan"},
			Reason:      "Anda mengalami sakit kepala. Pastikan cukup minum air, istirahat di ruangan gelap, dan hindari trigger makanan.",
		})
	}

	// DIABETES / GULA DARAH TINGGI (Referensi: Alodokter)
	if symptomNames["Diabetes"] || symptomNames["Gula Darah Tinggi"] || symptomNames["Hiperglikemia"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "diabetes",
			Title:       "🩸 Makanan untuk Diabetes (Ref: Alodokter)",
			Description: "Pilih karbohidrat kompleks dan makanan indeks glikemik rendah. Masak dengan cara dikukus, direbus, atau dipanggang.",
			Foods:       []string{"Beras merah (pengganti nasi putih)", "Oatmeal & gandum utuh", "Sayuran hijau (bayam, brokoli, kangkung, sawi)", "Ikan omega-3 (salmon, makarel, tuna, sarden)", "Tahu, tempe, edamame", "Buah segar (apel, pir, jeruk, stroberi, alpukat)", "Yogurt rendah lemak tanpa gula", "Kacang-kacangan", "Telur (putih telur)"},
			Avoid:       []string{"Nasi putih berlebihan", "Roti putih & kue manis", "Gula pasir & gula jawa", "Minuman manis & bersoda", "Jus kemasan", "Makanan gorengan", "Daging berlemak & jeroan", "Susu full cream", "Makanan dengan indeks glikemik tinggi"},
			Reason:      "Sumber: Alodokter. Penderita diabetes tetap boleh makan nasi, tapi pilih nasi merah dan batasi porsinya. Masak makanan dengan cara dikukus, direbus, atau dipanggang - hindari digoreng. Kontrol porsi makan dan rutin cek gula darah.",
		})
	}

	// ANEMIA / KURANG DARAH
	if symptomNames["Anemia"] || symptomNames["Kurang Darah"] || symptomNames["Pucat"] || symptomNames["Lemas"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "anemia",
			Title:       "🩸 Makanan untuk Anemia",
			Description: "Makanan tinggi zat besi untuk meningkatkan produksi sel darah merah",
			Foods:       []string{"Daging sapi tanpa lemak", "Hati ayam/sapi", "Bayam", "Brokoli", "Kacang merah", "Tahu tempe", "Telur", "Kerang", "Kurma", "Bit", "Kismis", "Vitamin C untuk penyerapan zat besi"},
			Avoid:       []string{"Teh bersamaan dengan makan (menghambat penyerapan zat besi)", "Kopi bersamaan dengan makan", "Susu bersamaan dengan suplemen zat besi", "Makanan tinggi kalsium saat makan"},
			Reason:      "Anda mengalami anemia. Konsumsi makanan tinggi zat besi bersama vitamin C untuk penyerapan optimal. Hindari teh/kopi 1 jam sebelum dan sesudah makan.",
		})
	}

	// SEMBELIT / KONSTIPASI
	if symptomNames["Sembelit"] || symptomNames["Konstipasi"] || symptomNames["Susah BAB"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "constipation",
			Title:       "🥗 Makanan untuk Melancarkan Pencernaan",
			Description: "Makanan tinggi serat untuk mengatasi sembelit",
			Foods:       []string{"Pepaya", "Pisang matang", "Sayuran hijau", "Kacang-kacangan", "Oatmeal", "Buah pir", "Prune (buah plum kering)", "Air putih minimal 8 gelas", "Chia seed", "Yogurt probiotik", "Ubi jalar"},
			Avoid:       []string{"Makanan olahan", "Daging merah berlebihan", "Makanan berlemak tinggi", "Alkohol", "Minuman berkafein berlebihan", "Makanan cepat saji"},
			Reason:      "Anda mengalami sembelit. Perbanyak serat dan air putih, serta lakukan olahraga ringan seperti jalan kaki untuk membantu pergerakan usus.",
		})
	}

	// ASAM URAT
	if symptomNames["Asam Urat"] || symptomNames["Gout"] || symptomNames["Nyeri Sendi"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "gout",
			Title:       "🦴 Makanan untuk Asam Urat",
			Description: "Makanan rendah purin untuk mengontrol asam urat",
			Foods:       []string{"Air putih minimal 10 gelas", "Sayuran (kecuali tertentu)", "Buah ceri", "Apel", "Pisang", "Susu rendah lemak", "Telur", "Tahu", "Kentang", "Roti gandum"},
			Avoid:       []string{"Jeroan (hati, ampela, otak)", "Daging merah", "Seafood (udang, kerang, kepiting)", "Alkohol (terutama bir)", "Minuman manis fruktosa tinggi", "Kacang-kacangan berlebihan", "Bayam dan asparagus berlebihan", "Sarden dan ikan teri"},
			Reason:      "Anda memiliki asam urat tinggi. Hindari makanan tinggi purin, perbanyak minum air putih, dan jaga berat badan ideal.",
		})
	}

	// KEHAMILAN
	if symptomNames["Hamil"] || symptomNames["Kehamilan"] || symptomNames["Morning Sickness"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "pregnancy",
			Title:       "🤰 Makanan untuk Ibu Hamil",
			Description: "Nutrisi penting untuk kesehatan ibu dan janin",
			Foods:       []string{"Sayuran hijau (asam folat)", "Salmon (omega-3, DHA)", "Telur", "Susu dan produk susu", "Daging tanpa lemak", "Kacang-kacangan", "Buah-buahan segar", "Ubi jalar", "Yogurt", "Alpukat", "Kurma"},
			Avoid:       []string{"Ikan tinggi merkuri (hiu, king mackerel)", "Daging/telur mentah", "Keju lunak tidak dipasteurisasi", "Kafein berlebihan", "Alkohol", "Jamu-jamuan tanpa resep dokter", "Makanan laut mentah (sushi)", "Nanas muda berlebihan"},
			Reason:      "Anda sedang hamil. Pastikan mendapat asam folat, zat besi, kalsium, dan protein cukup. Konsultasikan dengan dokter untuk suplemen prenatal.",
		})
	}

	// ALERGI MAKANAN
	if symptomNames["Alergi"] || symptomNames["Gatal-gatal"] || symptomNames["Ruam Kulit"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "allergy",
			Title:       "🚫 Tips untuk Alergi Makanan",
			Description: "Panduan menghindari alergen dan alternatif makanan",
			Foods:       []string{"Makanan segar non-olahan", "Sayuran dan buah lokal", "Nasi", "Daging ayam segar", "Ikan segar (jika tidak alergi)", "Minyak zaitun", "Air kelapa"},
			Avoid:       []string{"Makanan yang mengandung alergen Anda", "Makanan olahan (sering mengandung alergen tersembunyi)", "Saus dan bumbu kemasan", "Makanan restoran tanpa info alergen", "Susu (jika alergi susu)", "Kacang (jika alergi kacang)", "Seafood (jika alergi)"},
			Reason:      "Anda memiliki riwayat alergi. Selalu baca label makanan, bawa obat alergi, dan konsultasikan dengan dokter untuk tes alergi lengkap.",
		})
	}

	// BATUK
	if symptomNames["Batuk"] || symptomNames["Batuk Berdahak"] || symptomNames["Tenggorokan Gatal"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "cough",
			Title:       "🍵 Makanan Pereda Batuk",
			Description: "Makanan yang membantu meredakan batuk dan melegakan tenggorokan",
			Foods:       []string{"Madu hangat", "Jahe hangat", "Lemon hangat dengan madu", "Sup ayam", "Teh herbal", "Air hangat", "Nanas", "Bawang putih", "Kunyit"},
			Avoid:       []string{"Makanan dingin/es", "Makanan berminyak", "Gorengan", "Makanan pedas", "Susu (dapat memperbanyak lendir)", "Makanan manis berlebihan"},
			Reason:      "Anda mengalami batuk. Perbanyak minum air hangat, konsumsi madu untuk meredakan tenggorokan, dan istirahat yang cukup.",
		})
	}

	// DIARE
	if symptomNames["Diare"] || symptomNames["Mencret"] || symptomNames["Sakit Perut"] {
		recommendations = append(recommendations, models.FoodRecommendation{
			Category:    "diarrhea",
			Title:       "🍌 Makanan untuk Diare (BRAT Diet)",
			Description: "Makanan yang mudah dicerna untuk memulihkan pencernaan",
			Foods:       []string{"Pisang", "Nasi putih", "Roti tawar", "Apel (tanpa kulit)", "Oralit/larutan gula garam", "Air kelapa", "Kentang rebus", "Wortel rebus", "Bubur"},
			Avoid:       []string{"Susu dan produk susu", "Makanan pedas", "Makanan berminyak", "Sayuran mentah", "Buah-buahan asam", "Kafein", "Alkohol", "Makanan tinggi serat"},
			Reason:      "Anda mengalami diare. Fokus pada rehidrasi dengan oralit/air kelapa, makan BRAT diet (Banana, Rice, Applesauce, Toast), dan hindari makanan yang merangsang usus.",
		})
	}

	return recommendations
}


func generateExerciseRecommendations(user models.User, health models.HealthData, symptoms []models.Symptom) []models.ExerciseRecommendation {
	var recommendations []models.ExerciseRecommendation

	activityLevel := health.ActivityLevel
	if activityLevel == "" {
		activityLevel = user.ActivityLevel
	}

	bmiCategory := models.GetBMICategory(health.BMI)

	// Activity level based
	switch activityLevel {
	case "sedentary":
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "beginner",
			Title:       "Mulai dengan Aktivitas Ringan",
			Description: "Bangun kebiasaan olahraga secara bertahap",
			Exercises:   []string{"Jalan kaki 15-30 menit", "Stretching pagi", "Yoga pemula", "Berenang santai"},
			Duration:    "15-30 menit",
			Frequency:   "3-4 kali/minggu",
			Intensity:   "Ringan",
			Reason:      "Tingkat aktivitas Anda rendah, mulai perlahan",
		})
	case "light":
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "intermediate_light",
			Title:       "Tingkatkan Intensitas Olahraga",
			Description: "Tambah variasi dan durasi latihan",
			Exercises:   []string{"Jogging ringan", "Bersepeda santai", "Senam aerobik", "Pilates"},
			Duration:    "30-45 menit",
			Frequency:   "4-5 kali/minggu",
			Intensity:   "Ringan-Sedang",
			Reason:      "Anda sudah aktif ringan, tingkatkan intensitas",
		})
	case "moderate":
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "intermediate",
			Title:       "Variasikan Latihan Anda",
			Description: "Kombinasi kardio dan latihan kekuatan",
			Exercises:   []string{"Lari 5K", "HIIT workout", "Angkat beban", "Berenang lap", "Bulu tangkis"},
			Duration:    "45-60 menit",
			Frequency:   "5 kali/minggu",
			Intensity:   "Sedang",
			Reason:      "Tingkat aktivitas sedang, tambah variasi",
		})
	case "active":
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "advanced",
			Title:       "Pertahankan Performa",
			Description: "Jaga konsistensi dan hindari overtraining",
			Exercises:   []string{"Lari jarak jauh", "CrossFit", "Latihan interval", "Olahraga kompetitif"},
			Duration:    "60+ menit",
			Frequency:   "5-6 kali/minggu dengan 1 hari istirahat",
			Intensity:   "Tinggi",
			Reason:      "Anda sangat aktif, jaga keseimbangan",
		})
	}

	// BMI-specific
	if bmiCategory == "Overweight" || bmiCategory == "Obese" {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "weight_loss",
			Title:       "Olahraga untuk Menurunkan Berat",
			Description: "Kombinasi kardio untuk membakar kalori",
			Exercises:   []string{"Jalan cepat", "Berenang", "Sepeda statis", "Eliptical trainer", "Zumba"},
			Duration:    "45-60 menit",
			Frequency:   "5-6 kali/minggu",
			Intensity:   "Sedang",
			Reason:      "Fokus pada pembakaran kalori untuk penurunan berat badan",
		})
	}

	// Check for physical symptoms that affect exercise
	symptomNames := make(map[string]bool)
	for _, s := range symptoms {
		symptomNames[s.SymptomName] = true
		if s.SymptomName == "Nyeri Sendi" || s.SymptomName == "Nyeri Otot" {
			recommendations = append(recommendations, models.ExerciseRecommendation{
				Category:    "low_impact",
				Title:       "🦴 Olahraga Rendah Dampak",
				Description: "Aktivitas yang tidak membebani sendi dan otot",
				Exercises:   []string{"Berenang", "Yoga", "Tai Chi", "Bersepeda statis", "Water aerobics"},
				Duration:    "20-30 menit",
				Frequency:   "3-4 kali/minggu",
				Intensity:   "Ringan",
				Reason:      "Anda mengalami nyeri sendi/otot, pilih olahraga yang lembut",
			})
		}
	}

	// DEMAM DAN FLU - Istirahat, tidak olahraga berat
	if symptomNames["Demam"] || symptomNames["Flu"] || symptomNames["Pilek"] {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "recovery",
			Title:       "🤒 Istirahat saat Demam/Flu",
			Description: "Fokus pada pemulihan saat sakit",
			Exercises:   []string{"Istirahat total", "Stretching ringan di tempat tidur", "Pernapasan dalam", "Jalan pelan di dalam rumah"},
			Duration:    "5-10 menit",
			Frequency:   "Sesuai kemampuan",
			Intensity:   "Sangat Ringan",
			Reason:      "Saat demam/flu, prioritaskan istirahat. Olahraga berat dapat memperburuk kondisi. Mulai kembali olahraga secara bertahap setelah pulih.",
		})
	}

	// TEKANAN DARAH TINGGI
	if symptomNames["Tekanan Darah Tinggi"] || symptomNames["Hipertensi"] {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "blood_pressure",
			Title:       "❤️ Olahraga untuk Tekanan Darah",
			Description: "Aktivitas yang membantu mengontrol tekanan darah",
			Exercises:   []string{"Jalan kaki santai", "Berenang", "Bersepeda santai", "Yoga", "Tai Chi", "Senam ringan"},
			Duration:    "30-40 menit",
			Frequency:   "5 kali/minggu",
			Intensity:   "Ringan-Sedang",
			Reason:      "Olahraga teratur membantu menurunkan tekanan darah. Hindari angkat beban berat dan olahraga intensitas tinggi.",
		})
	}

	// STRES DAN KECEMASAN
	if symptomNames["Stres"] || symptomNames["Cemas"] || symptomNames["Kecemasan"] {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "stress_relief",
			Title:       "🧘 Olahraga Pereda Stres",
			Description: "Aktivitas fisik untuk mengurangi stres dan kecemasan",
			Exercises:   []string{"Yoga", "Tai Chi", "Jalan santai di alam", "Berenang", "Stretching/peregangan", "Pilates"},
			Duration:    "30-45 menit",
			Frequency:   "4-5 kali/minggu",
			Intensity:   "Ringan-Sedang",
			Reason:      "Olahraga melepaskan endorfin yang membantu mengurangi stres dan kecemasan. Fokus pada pernapasan dan gerakan mindful.",
		})
	}

	// GANGGUAN TIDUR
	if symptomNames["Gangguan Tidur"] || symptomNames["Insomnia"] || symptomNames["Sulit Tidur"] {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "sleep_improvement",
			Title:       "😴 Olahraga untuk Kualitas Tidur",
			Description: "Aktivitas yang membantu meningkatkan kualitas tidur",
			Exercises:   []string{"Yoga sebelum tidur", "Stretching malam", "Jalan kaki sore", "Tai Chi", "Pernapasan 4-7-8"},
			Duration:    "20-30 menit",
			Frequency:   "Setiap hari (hindari 3 jam sebelum tidur)",
			Intensity:   "Ringan",
			Reason:      "Olahraga teratur meningkatkan kualitas tidur, tetapi hindari olahraga intensif 3 jam sebelum tidur.",
		})
	}

	// KELELAHAN DAN BURNOUT
	if symptomNames["Kelelahan Fisik"] || symptomNames["Burnout"] || symptomNames["Lemas"] || symptomNames["Kelelahan Emosional (Burnout)"] {
		recommendations = append(recommendations, models.ExerciseRecommendation{
			Category:    "energy_boost",
			Title:       "⚡ Olahraga untuk Memulihkan Energi",
			Description: "Aktivitas ringan untuk mengembalikan energi tanpa menambah kelelahan",
			Exercises:   []string{"Jalan santai di luar ruangan", "Yoga restoratif", "Stretching pagi", "Berenang santai", "Berkebun"},
			Duration:    "15-30 menit",
			Frequency:   "3-4 kali/minggu, sesuai kemampuan",
			Intensity:   "Ringan",
			Reason:      "Saat burnout, olahraga ringan di luar ruangan dapat membantu memulihkan energi. Jangan memaksakan diri, dengarkan tubuh Anda.",
		})
	}

	return recommendations
}

func generateEmotionalRecommendations(emotionalState string, mentalSymptoms []models.Symptom) []models.EmotionalRecommendation {
	var recommendations []models.EmotionalRecommendation

	// Based on emotional state
	switch emotionalState {
	case "stressed":
		recommendations = append(recommendations, models.EmotionalRecommendation{
			EmotionalState: "stressed",
			Title:          "Kelola Stres Anda",
			Description:    "Teknik relaksasi untuk mengurangi stres",
			Activities:     []string{"Meditasi 10 menit", "Pernapasan dalam (4-7-8)", "Jalan santai di alam", "Mendengarkan musik menenangkan", "Journaling"},
			Tips:           []string{"Tidur cukup 7-8 jam", "Batasi screen time", "Luangkan waktu untuk diri sendiri", "Bicara dengan orang terdekat"},
			Reason:         "Anda sedang mengalami stres",
		})
	case "anxious":
		recommendations = append(recommendations, models.EmotionalRecommendation{
			EmotionalState: "anxious",
			Title:          "Atasi Kecemasan",
			Description:    "Aktivitas untuk menenangkan pikiran cemas",
			Activities:     []string{"Grounding technique (5-4-3-2-1)", "Progressive muscle relaxation", "Yoga restoratif", "Mewarnai mandala", "Merajut/craft"},
			Tips:           []string{"Hindari kafein berlebihan", "Batasi berita negatif", "Tetap terhubung dengan orang tersayang", "Fokus pada yang bisa dikontrol"},
			Reason:         "Anda sedang merasa cemas",
		})
	case "sad":
		recommendations = append(recommendations, models.EmotionalRecommendation{
			EmotionalState: "sad",
			Title:          "Tingkatkan Mood Anda",
			Description:    "Aktivitas untuk mengangkat suasana hati",
			Activities:     []string{"Olahraga ringan (endorfin)", "Bertemu teman", "Menonton film favorit", "Memasak makanan kesukaan", "Berkebun"},
			Tips:           []string{"Jangan isolasi diri", "Tetap jaga rutinitas", "Terpapar sinar matahari pagi", "Jika berlanjut, pertimbangkan konseling"},
			Reason:         "Anda sedang merasa sedih",
		})
	case "happy":
		recommendations = append(recommendations, models.EmotionalRecommendation{
			EmotionalState: "happy",
			Title:          "Pertahankan Kebahagiaan",
			Description:    "Aktivitas untuk menjaga mood positif",
			Activities:     []string{"Berbagi kebahagiaan", "Gratitude journal", "Lakukan hobi", "Quality time dengan keluarga", "Olahraga yang menyenangkan"},
			Tips:           []string{"Rayakan pencapaian kecil", "Bantu orang lain", "Simpan momen bahagia", "Tetap bersyukur"},
			Reason:         "Mood Anda sedang baik, pertahankan!",
		})
	case "neutral":
		recommendations = append(recommendations, models.EmotionalRecommendation{
			EmotionalState: "neutral",
			Title:          "Jaga Keseimbangan Emosi",
			Description:    "Aktivitas untuk kesejahteraan mental",
			Activities:     []string{"Mindfulness harian", "Olahraga rutin", "Hobi kreatif", "Sosialisasi sehat", "Belajar hal baru"},
			Tips:           []string{"Tetap jaga rutinitas sehat", "Check-in perasaan secara rutin", "Istirahat yang cukup"},
			Reason:         "Jaga keseimbangan emosional Anda",
		})
	}

	// Based on mental symptoms
	symptomActivities := map[string]models.EmotionalRecommendation{
		"Gangguan Tidur": {
			EmotionalState: "sleep_issue",
			Title:          "😴 Perbaiki Kualitas Tidur",
			Description:    "Tips dan aktivitas untuk tidur lebih berkualitas",
			Activities:     []string{"Rutinitas tidur tetap (jam sama setiap hari)", "Hindari gadget 1 jam sebelum tidur", "Mandi air hangat", "Aromatherapy lavender", "Membaca buku fisik", "Teknik relaksasi otot progresif"},
			Tips:           []string{"Jaga suhu kamar sejuk (18-22°C)", "Hindari kafein setelah jam 2 siang", "Olahraga pagi, hindari malam", "Konsisten jam tidur & bangun", "Gunakan masker mata jika perlu"},
			Reason:         "Anda mengalami gangguan tidur. Hindari kafein malam hari dan lakukan relaksasi ringan sebelum tidur.",
		},
		"Insomnia": {
			EmotionalState: "insomnia",
			Title:          "🌙 Atasi Insomnia",
			Description:    "Langkah-langkah untuk mengatasi kesulitan tidur",
			Activities:     []string{"Teknik pernapasan 4-7-8", "Body scan meditation", "White noise atau musik alam", "Journaling sebelum tidur", "Stretching ringan"},
			Tips:           []string{"Gunakan tempat tidur hanya untuk tidur", "Jangan lihat jam saat sulit tidur", "Bangun jika tidak bisa tidur 20 menit", "Hindari tidur siang terlalu lama"},
			Reason:         "Anda mengalami insomnia. Ciptakan rutinitas tidur yang konsisten dan lingkungan tidur yang nyaman.",
		},
		"Stres": {
			EmotionalState: "stress",
			Title:          "🧘 Kelola Stres dengan Efektif",
			Description:    "Teknik dan aktivitas untuk mengurangi stres",
			Activities:     []string{"Meditasi mindfulness 10 menit", "Teknik pernapasan dalam (4-7-8)", "Jalan santai di alam/taman", "Mendengarkan musik tenang", "Menulis jurnal perasaan", "Progressive muscle relaxation"},
			Tips:           []string{"Tidur cukup 7-8 jam", "Batasi screen time", "Luangkan waktu untuk diri sendiri", "Bicara dengan orang terdekat", "Batasi konsumsi berita negatif"},
			Reason:         "Anda mengalami stres. Disarankan meditasi 10 menit setiap hari dan menulis jurnal untuk mengekspresikan perasaan.",
		},
		"Cemas": {
			EmotionalState: "anxiety",
			Title:          "💆 Atasi Kecemasan",
			Description:    "Aktivitas untuk menenangkan pikiran yang cemas",
			Activities:     []string{"Grounding technique (5-4-3-2-1)", "Box breathing (4-4-4-4)", "Yoga restoratif", "Mewarnai mandala", "Merajut/craft", "Berjalan tanpa alas kaki di rumput"},
			Tips:           []string{"Hindari kafein berlebihan", "Batasi berita negatif", "Fokus pada hal yang bisa dikontrol", "Tetap terhubung dengan orang tersayang", "Rutin berolahraga ringan"},
			Reason:         "Anda merasa cemas. Latihan pernapasan dan grounding technique dapat membantu menenangkan pikiran.",
		},
		"Burnout": {
			EmotionalState: "burnout",
			Title:          "⚡ Pulihkan Diri dari Burnout",
			Description:    "Langkah pemulihan dari kelelahan emosional dan fisik",
			Activities:     []string{"Ambil cuti/istirahat", "Aktivitas luar ruangan (hiking, piknik)", "Digital detox", "Reconnect dengan hobi lama", "Quality time bersama keluarga", "Spa/self-care day"},
			Tips:           []string{"Set boundaries dengan jelas", "Belajar bilang 'tidak'", "Prioritaskan kesehatan", "Luangkan waktu berkualitas dengan keluarga", "Pertimbangkan konseling profesional"},
			Reason:         "Anda mengalami burnout. Istirahat, aktivitas luar ruangan, dan waktu bersama keluarga dapat membantu pemulihan.",
		},
		"Kelelahan Emosional (Burnout)": {
			EmotionalState: "emotional_exhaustion",
			Title:          "🌿 Pulihkan Energi Emosional",
			Description:    "Tips untuk memulihkan dari kelelahan emosional",
			Activities:     []string{"Aktivitas di alam terbuka", "Meditasi berjalan", "Hobi kreatif tanpa tekanan", "Waktu tenang sendirian", "Berkebun", "Bermain dengan hewan peliharaan"},
			Tips:           []string{"Kurangi tanggung jawab sementara", "Jangan merasa bersalah untuk istirahat", "Minta bantuan orang terdekat", "Hindari overthinking", "Fokus pada momen sekarang"},
			Reason:         "Anda mengalami kelelahan emosional. Penting untuk meluangkan waktu bersama keluarga dan melakukan aktivitas yang menyegarkan.",
		},
		"Kesepian Sosial": {
			EmotionalState: "lonely",
			Title:          "👨‍👩‍👧‍👦 Bangun Koneksi Sosial",
			Description:    "Aktivitas untuk mengurangi kesepian dan membangun hubungan",
			Activities:     []string{"Hubungi teman lama", "Ikut komunitas hobi", "Volunteer/sukarelawan", "Adopsi hewan peliharaan", "Ikut kelas/workshop", "Video call dengan keluarga jauh"},
			Tips:           []string{"Kualitas > kuantitas hubungan", "Jangan takut memulai percakapan", "Online community juga valid", "Jadi pendengar yang baik", "Rutin berkumpul dengan keluarga"},
			Reason:         "Anda merasa kesepian. Membangun koneksi dengan keluarga dan komunitas dapat membantu kesejahteraan mental.",
		},
	}

	for _, symptom := range mentalSymptoms {
		if rec, exists := symptomActivities[symptom.SymptomName]; exists {
			// Check if not already added
			exists := false
			for _, r := range recommendations {
				if r.EmotionalState == rec.EmotionalState {
					exists = true
					break
				}
			}
			if !exists {
				recommendations = append(recommendations, rec)
			}
		}
	}

	return recommendations
}
