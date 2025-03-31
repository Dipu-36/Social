package store

import (
	"context"
	"database/sql"
	"errors"
)

var(
	ErrorNotFound = errors.New("record not found")
)
//This struct is the aggregator of the interfaces that we have 
type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetByID(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface{
		Create(context.Context, *Users) error
	}
	Comments interface{
		GetByPostID(context.Context, int64) ([]Comments, error)
	}
}

//NewStorage initializes the Storage struct therefore it is called constructor it takes the database connection as input parameter and it passes this connection to the PostsStore and the UsersStore, 
// it then assigns PostsStore to the Posts interface and the UsersStore for the Users interface
//where PostsStore is the struct that implements all the methods that satisy the Posts interface and same goes for the Users and Commments interface
func NewStorage(db *sql.DB) Storage {
	return Storage {
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
