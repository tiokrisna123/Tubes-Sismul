package models

import "time"

// Article represents a health article/tip
type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	Summary   string    `json:"summary" gorm:"size:500"`
	Category  string    `json:"category" gorm:"size:50;not null"` // nutrisi, olahraga, mental, tidur, umum
	ImageURL  string    `json:"image_url" gorm:"size:500"`
	ReadTime  int       `json:"read_time"` // in minutes
	CreatedAt time.Time `json:"created_at"`
}

// GetSampleArticles returns sample health articles for seeding
func GetSampleArticles() []Article {
	return []Article{
		{
			Title:    "10 Makanan Super untuk Meningkatkan Imunitas Tubuh",
			Summary:  "Temukan makanan-makanan yang dapat membantu memperkuat sistem kekebalan tubuh Anda secara alami.",
			Content:  "Sistem imunitas yang kuat adalah kunci untuk menjaga kesehatan tubuh. Berikut adalah 10 makanan super yang dapat membantu meningkatkan imunitas:\n\n1. **Jeruk dan Buah Sitrus** - Kaya vitamin C yang meningkatkan produksi sel darah putih.\n\n2. **Brokoli** - Mengandung vitamin A, C, E serta antioksidan dan serat.\n\n3. **Bawang Putih** - Memiliki sifat antimikroba dan meningkatkan respons imun.\n\n4. **Jahe** - Membantu mengurangi peradangan dan meredakan mual.\n\n5. **Bayam** - Kaya vitamin C, antioksidan, dan beta karoten.\n\n6. **Yogurt** - Mengandung probiotik yang baik untuk kesehatan usus.\n\n7. **Almond** - Sumber vitamin E yang merupakan antioksidan kuat.\n\n8. **Kunyit** - Mengandung kurkumin dengan sifat anti-inflamasi.\n\n9. **Teh Hijau** - Kaya antioksidan EGCG yang meningkatkan fungsi imun.\n\n10. **Pepaya** - Mengandung papain dan vitamin C tinggi.",
			Category: "nutrisi",
			ImageURL: "https://images.unsplash.com/photo-1490645935967-10de6ba17061?w=800",
			ReadTime: 5,
		},
		{
			Title:    "Panduan Olahraga untuk Pemula: Mulai dari Mana?",
			Summary:  "Tips praktis untuk memulai rutinitas olahraga yang sehat dan berkelanjutan.",
			Content:  "Memulai rutinitas olahraga bisa terasa menakutkan, tapi dengan pendekatan yang tepat, Anda bisa membangun kebiasaan sehat yang bertahan lama.\n\n## Langkah Awal\n\n1. **Mulai dengan Jalan Kaki** - 15-20 menit per hari adalah awal yang baik.\n\n2. **Tetapkan Jadwal** - Pilih waktu tetap setiap hari untuk berolahraga.\n\n3. **Pemanasan Wajib** - Selalu lakukan pemanasan 5-10 menit sebelum olahraga.\n\n## Jenis Olahraga untuk Pemula\n\n- **Jalan cepat** - Low impact, bisa dilakukan di mana saja\n- **Berenang** - Baik untuk sendi dan otot seluruh tubuh\n- **Yoga** - Meningkatkan fleksibilitas dan keseimbangan\n- **Bersepeda** - Menyenangkan dan baik untuk kardio\n\n## Tips Penting\n\n- Jangan terlalu memaksakan diri di awal\n- Dengarkan tubuh Anda\n- Istirahat cukup antara sesi olahraga\n- Minum air yang cukup",
			Category: "olahraga",
			ImageURL: "https://images.unsplash.com/photo-1571019614242-c5c5dee9f50b?w=800",
			ReadTime: 4,
		},
		{
			Title:    "Mengelola Stres di Era Modern: Teknik yang Efektif",
			Summary:  "Pelajari cara-cara mengelola stres sehari-hari untuk kesehatan mental yang lebih baik.",
			Content:  "Stres adalah bagian dari kehidupan modern, tapi kita bisa belajar mengelolanya dengan efektif.\n\n## Teknik Relaksasi\n\n1. **Pernapasan Dalam** - Tarik napas dalam selama 4 detik, tahan 4 detik, hembuskan 4 detik.\n\n2. **Meditasi Mindfulness** - Fokus pada momen sekarang, 5-10 menit per hari.\n\n3. **Progressive Muscle Relaxation** - Tegangkan lalu rilekskan otot secara bergantian.\n\n## Gaya Hidup Anti-Stres\n\n- **Tidur cukup** - 7-9 jam per malam\n- **Olahraga teratur** - Minimal 30 menit, 3x seminggu\n- **Batasi kafein** - Terutama setelah jam 2 siang\n- **Sosial yang sehat** - Luangkan waktu dengan orang tersayang\n\n## Kapan Harus Mencari Bantuan?\n\nJika stres mengganggu aktivitas sehari-hari, tidur, atau nafsu makan selama lebih dari 2 minggu, pertimbangkan untuk berkonsultasi dengan profesional.",
			Category: "mental",
			ImageURL: "https://images.unsplash.com/photo-1506126613408-eca07ce68773?w=800",
			ReadTime: 6,
		},
		{
			Title:    "Kualitas Tidur: Rahasia Kesehatan yang Sering Diabaikan",
			Summary:  "Mengapa tidur berkualitas penting dan bagaimana cara mendapatkannya.",
			Content:  "Tidur berkualitas sama pentingnya dengan makan sehat dan olahraga teratur.\n\n## Manfaat Tidur Berkualitas\n\n- Meningkatkan daya ingat dan fokus\n- Memperkuat sistem imun\n- Mengatur berat badan\n- Memperbaiki mood\n- Mengurangi risiko penyakit kronis\n\n## Tips Tidur Lebih Baik\n\n1. **Jadwal Konsisten** - Tidur dan bangun pada jam yang sama setiap hari.\n\n2. **Lingkungan Nyaman** - Kamar gelap, sejuk, dan tenang.\n\n3. **Batasi Layar** - Hindari gadget 1 jam sebelum tidur.\n\n4. **Hindari Makan Berat** - Jangan makan besar 2-3 jam sebelum tidur.\n\n5. **Ritual Sebelum Tidur** - Mandi hangat, membaca buku, atau meditasi ringan.\n\n## Sleep Hygiene Checklist\n\n- [ ] Kamar tidur hanya untuk tidur\n- [ ] Suhu ruangan 18-22°C\n- [ ] Tidak ada TV di kamar tidur\n- [ ] Kasur dan bantal yang nyaman",
			Category: "tidur",
			ImageURL: "https://images.unsplash.com/photo-1541781774459-bb2af2f05b55?w=800",
			ReadTime: 5,
		},
		{
			Title:    "Hidrasi yang Tepat: Lebih dari Sekadar Minum Air",
			Summary:  "Panduan lengkap tentang kebutuhan cairan tubuh dan cara memenuhinya.",
			Content:  "Air adalah komponen vital tubuh yang membentuk sekitar 60% dari berat badan kita.\n\n## Mengapa Hidrasi Penting?\n\n- Mengatur suhu tubuh\n- Melancarkan pencernaan\n- Membuang racun dari tubuh\n- Menjaga kesehatan kulit\n- Meningkatkan energi dan fokus\n\n## Berapa Banyak Air yang Dibutuhkan?\n\nRata-rata kebutuhan air:\n- **Pria**: 3.7 liter per hari\n- **Wanita**: 2.7 liter per hari\n\nKebutuhan bisa bertambah jika:\n- Cuaca panas\n- Olahraga intensif\n- Sedang sakit\n\n## Tanda-tanda Dehidrasi\n\n- Urine berwarna gelap\n- Mulut kering\n- Pusing atau lelah\n- Kulit kering\n\n## Tips Minum Air Lebih Banyak\n\n1. Bawa botol air ke mana-mana\n2. Set reminder di handphone\n3. Minum segelas air sebelum makan\n4. Tambahkan irisan buah untuk rasa",
			Category: "nutrisi",
			ImageURL: "https://images.unsplash.com/photo-1548839140-29a749e1cf4d?w=800",
			ReadTime: 4,
		},
		{
			Title:    "Kesehatan Jantung: Langkah Preventif yang Bisa Anda Lakukan",
			Summary:  "Pelajari cara menjaga kesehatan jantung dengan perubahan gaya hidup sederhana.",
			Content:  "Penyakit jantung adalah penyebab kematian nomor satu di dunia, tapi sebagian besar kasus bisa dicegah.\n\n## Faktor Risiko yang Bisa Dikontrol\n\n- Tekanan darah tinggi\n- Kolesterol tinggi\n- Merokok\n- Kurang gerak\n- Kelebihan berat badan\n- Diabetes\n\n## Makanan untuk Jantung Sehat\n\n1. **Ikan berlemak** - Salmon, makarel, sarden (omega-3)\n2. **Kacang-kacangan** - Almond, walnut, kenari\n3. **Sayuran hijau** - Bayam, brokoli, kale\n4. **Buah beri** - Blueberry, strawberry\n5. **Oatmeal** - Menurunkan kolesterol\n\n## Olahraga untuk Jantung\n\n- Jalan cepat 30 menit/hari\n- Berenang\n- Bersepeda\n- Aerobik\n\n## Pemeriksaan Rutin\n\n- Tekanan darah setiap tahun\n- Kolesterol setiap 5 tahun\n- Gula darah setiap 3 tahun",
			Category: "umum",
			ImageURL: "https://images.unsplash.com/photo-1628348070889-cb656235b4eb?w=800",
			ReadTime: 6,
		},
	}
}
