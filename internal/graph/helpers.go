package graph

import (
	"database/sql"
	"strings"

	"github.com/gulmix/Social-Network/internal/graph/model"
	"github.com/gulmix/Social-Network/internal/models"
)

func toGraphQLUser(user *models.User) *model.User {
	return &model.User{
		ID:            user.ID,
		Email:         user.Email,
		Username:      user.Username,
		FirstName:     nullStringToStringPtr(user.FirstName),
		LastName:      nullStringToStringPtr(user.LastName),
		Bio:           nullStringToStringPtr(user.Bio),
		AvatarURL:     nullStringToStringPtr(user.AvatarURL),
		EmailVerified: user.EmailVerified,
		OauthProvider: nullStringToStringPtr(user.OAuthProvider),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

func toGraphQLUsers(users []*models.User) []*model.User {
	result := make([]*model.User, len(users))
	for i, user := range users {
		result[i] = toGraphQLUser(user)
	}
	return result
}

func toGraphQLPost(post *models.Post) *model.Post {
	imageUrls := parseImageUrls(post.ImageURLs)
	return &model.Post{
		ID:        post.ID,
		Content:   post.Content,
		ImageUrls: imageUrls,
		IsPublic:  post.IsPublic,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func toGraphQLPosts(posts []*models.Post) []*model.Post {
	result := make([]*model.Post, len(posts))
	for i, post := range posts {
		result[i] = toGraphQLPost(post)
	}
	return result
}

func toGraphQLComment(comment *models.Comment) *model.Comment {
	return &model.Comment{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}

func toGraphQLComments(comments []*models.Comment) []*model.Comment {
	result := make([]*model.Comment, len(comments))
	for i, comment := range comments {
		result[i] = toGraphQLComment(comment)
	}
	return result
}

func toGraphQLLike(like *models.Like) *model.Like {
	return &model.Like{
		ID:        like.ID,
		CreatedAt: like.CreatedAt,
	}
}

func toGraphQLLikes(likes []*models.Like) []*model.Like {
	result := make([]*model.Like, len(likes))
	for i, like := range likes {
		result[i] = toGraphQLLike(like)
	}
	return result
}

func toGraphQLFollow(follow *models.Follow) *model.Follow {
	return &model.Follow{
		ID:        follow.ID,
		CreatedAt: follow.CreatedAt,
	}
}

func parseImageUrls(imageUrls sql.NullString) []string {
	if !imageUrls.Valid || imageUrls.String == "" {
		return []string{}
	}
	urls := strings.Trim(imageUrls.String, "{}")
	if urls == "" {
		return []string{}
	}
	return strings.Split(urls, ",")
}

func nullStringToStringPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
