package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/gulmix/Social-Network/internal/models"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(group *models.Group) error {
	group.ID = uuid.New().String()
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	query := `INSERT INTO groups (id, name, description, avatar_url, creator_id, is_private, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(query, group.ID, group.Name, group.Description, group.AvatarURL, group.CreatorID, group.IsPrivate, group.CreatedAt, group.UpdatedAt)
	return err
}

func (r *GroupRepository) GetByID(id string) (*models.Group, error) {
	query := `SELECT id, name, description, avatar_url, creator_id, is_private, created_at, updated_at FROM groups WHERE id = $1`
	group := &models.Group{}
	err := r.db.QueryRow(query, id).Scan(&group.ID, &group.Name, &group.Description, &group.AvatarURL, &group.CreatorID, &group.IsPrivate, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) Update(group *models.Group) error {
	group.UpdatedAt = time.Now()
	query := `UPDATE groups SET name = $1, description = $2, avatar_url = $3, is_private = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(query, group.Name, group.Description, group.AvatarURL, group.IsPrivate, group.UpdatedAt, group.ID)
	return err
}

func (r *GroupRepository) Delete(id string) error {
	query := `DELETE FROM groups WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *GroupRepository) GetAll(limit, offset int) ([]*models.Group, error) {
	query := `SELECT id, name, description, avatar_url, creator_id, is_private, created_at, updated_at FROM groups WHERE is_private = false ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		g := &models.Group{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.AvatarURL, &g.CreatorID, &g.IsPrivate, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) GetUserGroups(userID string) ([]*models.Group, error) {
	query := `
		SELECT g.id, g.name, g.description, g.avatar_url, g.creator_id, g.is_private, g.created_at, g.updated_at
		FROM groups g
		JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.user_id = $1
		ORDER BY g.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		g := &models.Group{}
		if err := rows.Scan(&g.ID, &g.Name, &g.Description, &g.AvatarURL, &g.CreatorID, &g.IsPrivate, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

func (r *GroupRepository) AddMember(groupID, userID, role string) error {
	query := `INSERT INTO group_members (id, group_id, user_id, role, joined_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(query, uuid.New().String(), groupID, userID, role, time.Now())
	return err
}

func (r *GroupRepository) RemoveMember(groupID, userID string) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, groupID, userID)
	return err
}

func (r *GroupRepository) GetMembers(groupID string) ([]*models.GroupMember, error) {
	query := `SELECT id, group_id, user_id, role, joined_at FROM group_members WHERE group_id = $1 ORDER BY joined_at`
	rows, err := r.db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.GroupMember
	for rows.Next() {
		m := &models.GroupMember{}
		if err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.Role, &m.JoinedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *GroupRepository) IsMember(groupID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, groupID, userID).Scan(&exists)
	return exists, err
}

func (r *GroupRepository) IsAdmin(groupID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2 AND role = 'admin')`
	var exists bool
	err := r.db.QueryRow(query, groupID, userID).Scan(&exists)
	return exists, err
}

func (r *GroupRepository) IsCreator(groupID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM groups WHERE id = $1 AND creator_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, groupID, userID).Scan(&exists)
	return exists, err
}

func (r *GroupRepository) GetMemberCount(groupID string) (int, error) {
	query := `SELECT COUNT(*) FROM group_members WHERE group_id = $1`
	var count int
	err := r.db.QueryRow(query, groupID).Scan(&count)
	return count, err
}

func (r *GroupRepository) CreatePost(post *models.GroupPost) error {
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	query := `INSERT INTO group_posts (id, group_id, user_id, content, image_urls, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, post.ID, post.GroupID, post.UserID, post.Content, post.ImageURLs, post.CreatedAt, post.UpdatedAt)
	return err
}

func (r *GroupRepository) GetPosts(groupID string, limit, offset int) ([]*models.GroupPost, error) {
	query := `SELECT id, group_id, user_id, content, image_urls, created_at, updated_at FROM group_posts WHERE group_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, groupID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.GroupPost
	for rows.Next() {
		p := &models.GroupPost{}
		if err := rows.Scan(&p.ID, &p.GroupID, &p.UserID, &p.Content, &p.ImageURLs, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *GroupRepository) DeletePost(postID string) error {
	query := `DELETE FROM group_posts WHERE id = $1`
	_, err := r.db.Exec(query, postID)
	return err
}

func (r *GroupRepository) IsPostOwner(postID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM group_posts WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(query, postID, userID).Scan(&exists)
	return exists, err
}
