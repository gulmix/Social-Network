package graph

import (
	"database/sql"
	"strings"
	"time"

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

func nullTimeToTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func toGraphQLConversation(conv *models.Conversation) *model.Conversation {
	return &model.Conversation{
		ID:        conv.ID,
		IsGroup:   conv.IsGroup,
		Name:      nullStringToStringPtr(conv.Name),
		CreatedAt: conv.CreatedAt,
		UpdatedAt: conv.UpdatedAt,
	}
}

func toGraphQLMessage(msg *models.Message) *model.Message {
	return &model.Message{
		ID:        msg.ID,
		Content:   nullStringToStringPtr(msg.Content),
		MediaURL:  nullStringToStringPtr(msg.MediaURL),
		CreatedAt: msg.CreatedAt,
	}
}

func toGraphQLMessages(msgs []*models.Message) []*model.Message {
	result := make([]*model.Message, len(msgs))
	for i, msg := range msgs {
		result[i] = toGraphQLMessage(msg)
	}
	return result
}

func toGraphQLGroup(group *models.Group) *model.Group {
	return &model.Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: nullStringToStringPtr(group.Description),
		AvatarURL:   nullStringToStringPtr(group.AvatarURL),
		IsPrivate:   group.IsPrivate,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
	}
}

func toGraphQLGroups(groups []*models.Group) []*model.Group {
	result := make([]*model.Group, len(groups))
	for i, group := range groups {
		result[i] = toGraphQLGroup(group)
	}
	return result
}

func toGraphQLGroupPost(post *models.GroupPost) *model.GroupPost {
	return &model.GroupPost{
		ID:        post.ID,
		Content:   post.Content,
		ImageUrls: parseImageUrls(post.ImageURLs),
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
}

func toGraphQLGroupPosts(posts []*models.GroupPost) []*model.GroupPost {
	result := make([]*model.GroupPost, len(posts))
	for i, post := range posts {
		result[i] = toGraphQLGroupPost(post)
	}
	return result
}

func toGraphQLEvent(event *models.Event) *model.Event {
	return &model.Event{
		ID:          event.ID,
		Title:       event.Title,
		Description: nullStringToStringPtr(event.Description),
		Location:    nullStringToStringPtr(event.Location),
		StartTime:   event.StartTime,
		EndTime:     nullTimeToTimePtr(event.EndTime),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func toGraphQLEvents(events []*models.Event) []*model.Event {
	result := make([]*model.Event, len(events))
	for i, event := range events {
		result[i] = toGraphQLEvent(event)
	}
	return result
}

func toGraphQLNotification(n *models.Notification) *model.Notification {
	return &model.Notification{
		ID:          n.ID,
		Type:        n.Type,
		ReferenceID: nullStringToStringPtr(n.ReferenceID),
		Content:     n.Content,
		IsRead:      n.IsRead,
		CreatedAt:   n.CreatedAt,
	}
}

func toGraphQLNotifications(notifs []*models.Notification) []*model.Notification {
	result := make([]*model.Notification, len(notifs))
	for i, n := range notifs {
		result[i] = toGraphQLNotification(n)
	}
	return result
}

func toGraphQLStory(story *models.Story) *model.Story {
	return &model.Story{
		ID:        story.ID,
		MediaURL:  story.MediaURL,
		Content:   nullStringToStringPtr(story.Content),
		ExpiresAt: story.ExpiresAt,
		CreatedAt: story.CreatedAt,
	}
}

func toGraphQLStories(stories []*models.Story) []*model.Story {
	result := make([]*model.Story, len(stories))
	for i, story := range stories {
		result[i] = toGraphQLStory(story)
	}
	return result
}

func toGraphQLStoryView(view *models.StoryView) *model.StoryView {
	return &model.StoryView{
		ID:       view.ID,
		ViewedAt: view.ViewedAt,
	}
}

func toGraphQLStoryViews(views []*models.StoryView) []*model.StoryView {
	result := make([]*model.StoryView, len(views))
	for i, view := range views {
		result[i] = toGraphQLStoryView(view)
	}
	return result
}
