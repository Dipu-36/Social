package store

import (
	"context"
	"database/sql"

	
)

type Comments struct {
	ID int64 `json:"id"`
	PostID int64 `json:"post_id"`
	UserID int64 `json:"user_id"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	User  Users `json:"users"`

}
type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comments, error){
	query := `
		SELECT c.id, c.post_id, c.user_id, c.created_at, users.username, users.id FROM comments c
		JOIN users ON users.id = c.user_id 
		WHERE c.post_id = $1 
		ORDER BY c.created_at DESC;
	`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err!=nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comments{}
	for rows.Next(){
		var c Comments
		c.User = Users{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err!=nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	
	return comments, nil
}