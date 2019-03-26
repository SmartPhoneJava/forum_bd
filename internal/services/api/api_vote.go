package api

import (
	"escapade/internal/models"
	"net/http"
)

// Vote vote
func (h *Handler) Vote(rw http.ResponseWriter, r *http.Request) {
	const place = "Vote"
	var (
		thread models.Thread
		vote   models.Vote
		err    error
		slug   string
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if vote, err = getVote(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, err = h.DB.CreateVote(vote, slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusOK, place)
	return
}
