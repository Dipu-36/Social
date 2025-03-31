package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

//This is the model for Posts
type Post struct {
	ID int64 			`json:"id"`
	Content string 		`json:"content"`
	Title string 		`json:"title"`
	UserID int64 		`json:"user_id"`
	Tags []string 		`json:"tags"`
	CreatedAt string 	`json:"created_at"`
	UpdatedAt string 	`json:"updated_at"`
	Comments []Comments `json:"comments"`
}

//This struct is used to interact with the Db connection
type PostsStore struct {
	db *sql.DB
}

//Create method constructs and executes an INSERT INTO SQL query to add a new post to the databse, //this provides context for the request like timeouts, deadlines, etc it is passed to make sure that the DB query honrs any external context like client request timeout
func (s * PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING	id, created_at, updated_at
	`
	//This is where actual query execution happens
	err := s.db.QueryRowContext(
		ctx, //The context is passed to the database query 
		query,//The SQL query is passed as a string
		post.Content,//
		post.Title,
		post.UserID,
		pq.Array(post.Tags),

	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err!=nil{
		return err
	}
	return nil
}

func (s *PostsStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	query := `
	SELECT id, user_id, title, content, created_at, updated_at, tags
	FROM posts
	WHERE id = $1`
	//Now after retreiving the data from the databaase we want to store the results into the struct temporarily so that we can serialize the data into json to send to the client, for this purpose we use pointers so that we can storer the actual data and not the copy of the data.
	var post Post
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
	)
	if err!=nil {
		switch{
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}
	return &post, nil 
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error{
	query := `DELETE FROM posts WHERE id = $1`

	res,err := s.db.ExecContext(ctx, query, id)
	if err!=nil{
		return err
	}
	rows, err := res.RowsAffected()
	if err!=nil{
		return err
	}
	if rows == 0 {
		return ErrorNotFound
	}
	return nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error{
	query := `
		UPDATE posts
		SET title = $1, content = $2
		WHERE  id = $3 
	`
	_, err := s.db.ExecContext(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return err
	}

	return nil
}