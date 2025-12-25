package handlers

import (
	"health-tracker/database"
	"health-tracker/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPosts returns all forum posts
func GetPosts(c *gin.Context) {
	userID, _ := c.Get("userID")
	var posts []models.Post

	if err := database.DB.Preload("User").Order("created_at DESC").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	// Build response with user names and like status
	var response []models.PostResponse
	for _, post := range posts {
		isLiked := false
		if userID != nil {
			var like models.Like
			if err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&like).Error; err == nil {
				isLiked = true
			}
		}

		response = append(response, models.PostResponse{
			ID:            post.ID,
			UserID:        post.UserID,
			UserName:      post.User.Name,
			Title:         post.Title,
			Content:       post.Content,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			IsLiked:       isLiked,
			CreatedAt:     post.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
}

// CreatePost creates a new forum post
func CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := models.Post{
		UserID:    userID.(uint),
		Title:     input.Title,
		Content:   input.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Load user data
	database.DB.Preload("User").First(&post, post.ID)

	c.JSON(http.StatusCreated, models.PostResponse{
		ID:            post.ID,
		UserID:        post.UserID,
		UserName:      post.User.Name,
		Title:         post.Title,
		Content:       post.Content,
		LikesCount:    0,
		CommentsCount: 0,
		IsLiked:       false,
		CreatedAt:     post.CreatedAt,
	})
}

// GetPost returns a single post with comments
func GetPost(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	var post models.Post

	if err := database.DB.Preload("User").Preload("Comments.User").First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	isLiked := false
	if userID != nil {
		var like models.Like
		if err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&like).Error; err == nil {
			isLiked = true
		}
	}

	// Build comments response
	var comments []models.CommentResponse
	for _, comment := range post.Comments {
		comments = append(comments, models.CommentResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			UserName:  comment.User.Name,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"post": models.PostResponse{
			ID:            post.ID,
			UserID:        post.UserID,
			UserName:      post.User.Name,
			Title:         post.Title,
			Content:       post.Content,
			LikesCount:    post.LikesCount,
			CommentsCount: post.CommentsCount,
			IsLiked:       isLiked,
			CreatedAt:     post.CreatedAt,
		},
		"comments": comments,
	})
}

// AddComment adds a comment to a post
func AddComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := models.Comment{
		PostID:    post.ID,
		UserID:    userID.(uint),
		Content:   input.Content,
		CreatedAt: time.Now(),
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	// Update comments count
	database.DB.Model(&post).Update("comments_count", post.CommentsCount+1)

	// Load user data
	database.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, models.CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		UserName:  comment.User.Name,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	})
}

// ToggleLike toggles a like on a post
func ToggleLike(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var existingLike models.Like
	isLiked := false

	if err := database.DB.Where("post_id = ? AND user_id = ?", post.ID, userID).First(&existingLike).Error; err == nil {
		// Unlike - remove existing like
		database.DB.Delete(&existingLike)
		database.DB.Model(&post).Update("likes_count", post.LikesCount-1)
		isLiked = false
	} else {
		// Like - add new like
		like := models.Like{
			PostID:    post.ID,
			UserID:    userID.(uint),
			CreatedAt: time.Now(),
		}
		database.DB.Create(&like)
		database.DB.Model(&post).Update("likes_count", post.LikesCount+1)
		isLiked = true
	}

	// Get updated count
	database.DB.First(&post, postID)

	c.JSON(http.StatusOK, gin.H{
		"is_liked":    isLiked,
		"likes_count": post.LikesCount,
	})
}

// DeletePost deletes a post (only by owner)
func DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	postID := c.Param("id")
	var post models.Post
	if err := database.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to delete this post"})
		return
	}

	// Delete comments and likes first
	database.DB.Where("post_id = ?", post.ID).Delete(&models.Comment{})
	database.DB.Where("post_id = ?", post.ID).Delete(&models.Like{})
	database.DB.Delete(&post)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
