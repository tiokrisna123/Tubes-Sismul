package models

import "time"

// Post represents a forum post
type Post struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	Title         string    `json:"title" gorm:"size:200;not null"`
	Content       string    `json:"content" gorm:"type:text;not null"`
	LikesCount    int       `json:"likes_count" gorm:"default:0"`
	CommentsCount int       `json:"comments_count" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Comments      []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Likes         []Like    `json:"-" gorm:"foreignKey:PostID"`
}

// Comment represents a comment on a post
type Comment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at"`
}

// Like represents a like on a post
type Like struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PostID    uint      `json:"post_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// PostResponse is the response structure for a post
type PostResponse struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	UserName      string    `json:"user_name"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	IsLiked       bool      `json:"is_liked"`
	CreatedAt     time.Time `json:"created_at"`
}

// CommentResponse is the response structure for a comment
type CommentResponse struct {
	ID        uint      `json:"id"`
	PostID    uint      `json:"post_id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
