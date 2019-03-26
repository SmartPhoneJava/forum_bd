package api

import (
	"escapade/internal/models"
	"net/http"
)

// GetPostfull get postfull
func (h *Handler) GetPostfull(rw http.ResponseWriter, r *http.Request) {
	const place = "GetPost"
	var (
		post         models.Postfull
		err          error
		id           int
		related      string
		existRelated bool
	)

	rw.Header().Set("Content-Type", "application/json")

	if id, err = getPostID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	existRelated, related = getRelated(r)

	if post, err = h.DB.GetPostfull(existRelated, related, id); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusConflict, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, post, place)
	printResult(err, http.StatusOK, place)
	return
}
