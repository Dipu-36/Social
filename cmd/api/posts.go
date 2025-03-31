package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/Dipu-36/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string
const postCtx postKey = "post" 
//This is the http Payload for the Create method so that user can send or write only those data that he is permitted to write like he is not authorized to write/create his own id 
type CreatePostPayload struct{
	Title string `json:"title" validate:"required,max=100"`
	Content string `json:"content" validate:"required,max=1000"`
	Tags []string `json:"tags"`
}

func (app *application) createPosthandler(w http.ResponseWriter, r *http.Request) {
	//The below block of code is not reccomended as it accpets everything that is in the post and user can overwrite the data and corrupt the DB

	// var post store.Post
	// if err := readJSON(w, r, post); err!=nil{
	// 	writeJSON(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	//therefore we create a payload so that the user can send or write only those data that he has been permitted 
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err!=nil{
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err!=nil{
		app.badRequestError(w, r, err)
		return
	}
	post := &store.Post{
		Title: payload.Title,
		Content: payload.Content,
		Tags: payload.Tags,
		//TODO: Change after auth
		UserID: 1,
	}
	
	ctx := r.Context()//ctx variable returns the context.Context of the HTTP request, and thus carries the authentication or trace ID data, deadlines(timeouts), cancellation signals(if client disconnects or request is canceled)

	//we pass the ctx or the context.Context of the request inside the data layer to handle timeouts or cancellation thus it will cancel the DB read/write query also called zombie query

	if err:= app.store.Posts.Create(ctx, post); err!=nil{
		app.internalServerError(w, r, err)
		return 
	}

	if err:= writeJSON(w, http.StatusCreated, post); err!=nil{
		app.internalServerError(w, r, err)
		return 
	}
}

func(app *application) getPosthandler(w http.ResponseWriter, r *http.Request){
	post := getPostFromCtx(r)

	
	//fetching the comments of the post using the qury from the storage layer
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil{
		app.internalServerError(w, r, err)
		return
	}
	//assiging the comments to the Post struct field
	post.Comments = comments

	if err:= writeJSON(w, http.StatusOK, post); err!=nil{
		app.internalServerError(w, r, err)
		return 
	}
}

func (app *application) deletePosthandler(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err!=nil{
		app.internalServerError(w, r, err)
		return 
	}
	if err := app.store.Posts.Delete(ctx, id); err!=nil{
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

//This is http Payload for Update method 
type UpdatePostpayload struct {
	Title 	*string `json:"title" validate:"omitempty,max=100"`
	Content *string `json:"content" validate:"omitempty,max=1000"`

}

func (app *application) updatePosthandler(w http.ResponseWriter, r *http.Request){
	//first check if the post exists or not before updating therefore
	post := getPostFromCtx(r)

	var payload UpdatePostpayload
	//Unmarshalling it 
	if err := readJSON(w, r, &payload); err!= nil{
		app.badRequestError(w, r, err)
		return 
	}

	//Validating it 
	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	//Actual execution of query in database
	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	//Writing the response to the ResponseWriter
	 if err := writeJSON(w, http.StatusOK, post); err != nil{
		app.internalServerError(w, r, err)
	 }
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler{
	//using http.HandlerFunc to convert this into http.Handler and return it 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err!=nil{
			app.internalServerError(w, r, err)
			//writeJSONError(w, http.StatusInternalServerError, err.Error()) //a bad practice cause it can print the stack trace which can leak internal details
			return 
		}
		ctx := r.Context()
	
		post, err := app.store.Posts.GetByID(ctx, id)
		if err!=nil{
			switch{
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		//Never mutate a context, always create from scratch 
		//Therefore we are passing the parsed post to the request context
		//Here we are storing the context as key-value pair with the key being postCtx and value being *store.Post type
		ctx = context.WithValue(ctx, postCtx, post)
		//Now we pass the context to the next handler in the middleware chain, which is getPostFromCtx
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//Final handler of the postsContextMiddleware chain, which accesses the post from the request's context 
func getPostFromCtx(r *http.Request) *store.Post{
	//Here we are accessing the value associated with the context and we are type asserting it to *store.Post because.Value() returns an empty interface
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}