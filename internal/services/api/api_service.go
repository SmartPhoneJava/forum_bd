package api

import (
	"escapade/internal/models"
	"net/http"
)

// GetStatus get status
func (h *Handler) GetStatus(rw http.ResponseWriter, r *http.Request) {
	const place = "GetStatus"
	var (
		status models.Status
		err    error
	)

	rw.Header().Set("Content-Type", "application/json")

	if status, err = h.DB.GetStatus(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusConflict, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, status, place)
	printResult(err, http.StatusOK, place)
	return
}

// Clear clear database
func (h *Handler) Clear(rw http.ResponseWriter, r *http.Request) {
	const place = "Clear"
	var (
		status models.Status
		err    error
	)

	rw.Header().Set("Content-Type", "application/json")

	if status, err = h.DB.Clear(); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, status, place)
	printResult(nil, http.StatusOK, place)
	return
}
