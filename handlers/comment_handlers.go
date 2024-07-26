package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/A-Victory/blog/models"
	"github.com/go-chi/chi/v5"
)

func (httpConfig *HttpHandler) Comment(w http.ResponseWriter, r *http.Request) {

	user, err := httpConfig.getUser(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "failed to retrieve user's details: " + err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == "POST" {

		postid := chi.URLParam(r, "postId")
		postID, err := strconv.Atoi(postid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error parsing URL param..."}}
			json.NewEncoder(w).Encode(response)
			return
		}

		newComment := models.Comment{}

		if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		newComment.AuthorID = user.ID
		newComment.Postid = postID

		id, err := httpConfig.db.AddComment(newComment)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		if id == 0 {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "database error", Data: map[string]interface{}{"msg": fmt.Sprintf("failed to add comment to post with id: %d", postID)}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := customResponse{Status: http.StatusOK, Message: "successfully added comment", Data: map[string]interface{}{"msg": fmt.Sprintf("successfully added comment to post with id: %d", postID)}}
		json.NewEncoder(w).Encode(response)
	}

	if r.Method == "PUT" {

		id := chi.URLParam(r, "id")
		if id != "" {
			commentID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error converting string..."}}
				json.NewEncoder(w).Encode(response)
				return
			}
			updateComment := models.Comment{}

			if err := json.NewDecoder(r.Body).Decode(&updateComment); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"message": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
			updateComment.ID = commentID

			id, err := httpConfig.db.EditComment(updateComment)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
			if id == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response := customResponse{Status: http.StatusBadRequest, Message: "database error", Data: map[string]interface{}{"msg": fmt.Sprintf("failed to update comment with id: %d", commentID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "successfully updated comment", Data: map[string]interface{}{"msg": fmt.Sprintf("successfully updated comment with id %d", commentID)}}
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"msg": "no id provided"}}
			json.NewEncoder(w).Encode(response)
			return
		}

	}

	if r.Method == "GET" {

		id := chi.URLParam(r, "postId")
		postID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error converting string..."}}
			json.NewEncoder(w).Encode(response)
			return
		}

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

		comments, err := httpConfig.db.GetComments(postID, limit, offset)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := customResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"comments": comments}}
		json.NewEncoder(w).Encode(response)

		// here retrieve all comments associated to a particular post
	}

	if r.Method == "DELETE" {

		id := chi.URLParam(r, "id")

		if id != "" {

			commentID, err := strconv.Atoi(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "error converting string..."}}
				json.NewEncoder(w).Encode(response)
				return
			}

			id, err := httpConfig.db.DeleteComment(commentID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := customResponse{Status: http.StatusInternalServerError, Message: "server error", Data: map[string]interface{}{"msg": "database connection error: " + err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}

			if id <= 0 {
				w.WriteHeader(http.StatusOK)
				response := customResponse{Status: http.StatusOK, Message: "no post found", Data: map[string]interface{}{"msg": fmt.Sprintf("no comment found with id: %d", commentID)}}
				json.NewEncoder(w).Encode(response)
				return
			}

			w.WriteHeader(http.StatusOK)
			response := customResponse{Status: http.StatusOK, Message: "successfully deleted comment", Data: map[string]interface{}{"comment_id": commentID}}
			json.NewEncoder(w).Encode(response)

		} else {
			w.WriteHeader(http.StatusBadRequest)
			response := customResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"msg": "no id provided"}}
			json.NewEncoder(w).Encode(response)
			return
		}
	}
}
