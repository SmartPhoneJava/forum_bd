package api

import (
	data "escapade/internal/database"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
)

// CreatePosts create posts
func (h *Handler) CreatePosts(rw http.ResponseWriter, r *http.Request) {
	const place = "CreatePosts"
	var (
		posts []models.Post
		err   error
		slug  string
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if posts, err = getPosts(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if posts, err = h.DB.CreatePost(posts, slug); err != nil {
		if err.Error() != re.ErrorPostConflict().Error() {
			rw.WriteHeader(http.StatusNotFound)
			sendErrorJSON(rw, err, place)
		} else {
			rw.WriteHeader(http.StatusConflict)
			sendErrorJSON(rw, err, place)
		}
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, posts, place)
	printResult(err, http.StatusCreated, place)
	return
}

// GetPosts get posts
func (h *Handler) GetPosts(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPosts"
	var (
		posts      []models.Post
		slug       string
		sort       string
		limit      int
		since      int
		err        error
		existLimit bool
		existSince bool
		desc       bool
		qgc        data.QueryGetConditions
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existLimit, limit, err = getLimit(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if existSince, since, err = getIDmin(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	desc = getDesc(r)
	sort = getSort(r)

	qgc.InitPost(existSince, since, existLimit, limit, desc)

	if posts, err = h.DB.GetPosts(slug, qgc, sort); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, posts, place)
	printResult(err, http.StatusOK, place)
	return
}

// UpdatePost update post
func (h *Handler) UpdatePost(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdatePost"
	var (
		post models.Post
		err  error
		id   int
	)

	rw.Header().Set("Content-Type", "application/json")

	if post, err = getPost(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if id, err = getPostID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if post, err = h.DB.UpdatePost(post, id); err != nil {
		//if err.Error() == re.ErrorForumUserNotExist().Error() {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		// } else {
		// 	rw.WriteHeader(http.StatusConflict)
		// 	sendSuccessJSON(rw, forum, place)
		// }
		printResult(err, http.StatusNotFound, place)
		return
	}

	post.Print()
	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, post, place)
	printResult(err, http.StatusOK, place)
	return
}
