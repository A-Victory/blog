package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/A-Victory/blog/models"
	"github.com/go-chi/chi/v5"
)

func (httpConfig *HttpHandler) Post(w http.ResponseWriter, r *http.Request) {

	user, err := httpConfig.getUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "failed to retrieve user's details: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == "POST" {

		newPost := models.Post{}

		if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		newPost.AuthorID = user.ID

		id, err := httpConfig.db.CreatePost(newPost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := customResponse{Status: http.StatusOK, Message: "successfully created post", Data: map[string]interface{}{"post_id": id}}
		json.NewEncoder(w).Encode(response)

		// save the post to the database and return a successful response
	}

	if r.Method == "PUT" {
		id := chi.URLParam(r, "id")
		if id != "" {

			postID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error parsing URL param..."}}
				json.NewEncoder(w).Encode(response)
				return
			}
			log.Println(postID)

			postToUpdate := models.Post{}

			if err := json.NewDecoder(r.Body).Decode(&postToUpdate); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}

			post, err := httpConfig.db.GetPostByID(postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
			if post == (models.Post{}) {
				w.WriteHeader(http.StatusNotFound)
				response := customResponse{Status: http.StatusNotFound, Message: "post not found", Data: map[string]interface{}{"msg": fmt.Sprintf("no post found with id %d", postID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			if post.AuthorID != user.ID {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "invalid request", Data: map[string]interface{}{"msg": fmt.Sprintf("not the author of post with id: %d", postID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			postToUpdate.AuthorID = post.AuthorID
			postToUpdate.ID = postID

			id, err := httpConfig.db.UpdatePost(postToUpdate)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}

			if id == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "database error", Data: map[string]interface{}{"msg": fmt.Sprintf("failed to update post with id: %d", post.ID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "successfully updated post", Data: map[string]interface{}{"post_id": postID}}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"msg": "no id provided"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		// update the post in the database after which you return the id of the post updated
	}

	if r.Method == "GET" {

		id := chi.URLParam(r, "id")

		if id != "" {
			// return the specific post from the database
			postID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error converting string..."}}
				json.NewEncoder(w).Encode(response)
				return
			}
			post, err := httpConfig.db.GetPostByID(postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"post": post}}
			json.NewEncoder(w).Encode(response)

		} else {
			query := r.URL.Query()

			// Get pagination parameters
			page, _ := strconv.Atoi(query.Get("page"))
			if page < 1 {
				page = 1
			}
			limit, _ := strconv.Atoi(query.Get("limit"))
			if limit < 1 {
				limit = 10
			}
			offset := (page - 1) * limit

			// Get search parameter
			search := query.Get("search")
			// return all the posts taking into account the pagination and search parameters
			posts, err := httpConfig.db.GetPosts(&limit, &offset, &search)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"posts": posts}}
			json.NewEncoder(w).Encode(response)
		}
	}

	if r.Method == "DELETE" {

		id := chi.URLParam(r, "id")

		if id != "" {
			// return an error specificing that the user needs to include the id of the post
			postID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error converting string..."}}
				json.NewEncoder(w).Encode(response)
				return
			}

			post, err := httpConfig.db.GetPostByID(postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
			if post == (models.Post{}) {
				w.WriteHeader(http.StatusNotFound)
				response := customResponse{Status: http.StatusNotFound, Message: "post not found", Data: map[string]interface{}{"msg": fmt.Sprintf("no post found with id %d", postID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			if post.AuthorID != user.ID {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "invalid request", Data: map[string]interface{}{"msg": fmt.Sprintf("not the author of post with id: %d", postID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			id, err := httpConfig.db.DeletePost(postID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}

			if id == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "database error", Data: map[string]interface{}{"msg": fmt.Sprintf("failed to update post with id: %d", post.ID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "successfully deleted post", Data: map[string]interface{}{"post_id": id}}
			json.NewEncoder(w).Encode(response)

		} else {
			// include response for when the user does not indicate id
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"msg": "no id provided"}}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
