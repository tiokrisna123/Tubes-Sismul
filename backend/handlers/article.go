package handlers

import (
	"health-tracker/database"
	"health-tracker/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetArticles returns all articles with optional category filter
func GetArticles(c *gin.Context) {
	var articles []models.Article
	category := c.Query("category")

	query := database.DB.Order("created_at DESC")
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// GetArticle returns a single article by ID
func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

// GetArticleCategories returns all available categories
func GetArticleCategories(c *gin.Context) {
	categories := []map[string]string{
		{"id": "all", "name": "Semua", "icon": "📚"},
		{"id": "nutrisi", "name": "Nutrisi", "icon": "🥗"},
		{"id": "olahraga", "name": "Olahraga", "icon": "🏃"},
		{"id": "mental", "name": "Kesehatan Mental", "icon": "🧠"},
		{"id": "tidur", "name": "Tidur", "icon": "😴"},
		{"id": "umum", "name": "Umum", "icon": "❤️"},
	}
	c.JSON(http.StatusOK, categories)
}

// SearchArticles searches articles by keyword
func SearchArticles(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search keyword required"})
		return
	}

	var articles []models.Article
	searchPattern := "%" + keyword + "%"
	
	if err := database.DB.Where("title LIKE ? OR content LIKE ?", searchPattern, searchPattern).
		Order("created_at DESC").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search articles"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// Pagination helper
func GetPaginatedArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	category := c.Query("category")

	offset := (page - 1) * limit

	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{})
	if category != "" && category != "all" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"total":    total,
		"page":     page,
		"limit":    limit,
		"pages":    (total + int64(limit) - 1) / int64(limit),
	})
}
