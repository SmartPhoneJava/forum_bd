package api

import (
	"net/http"

	"github.com/SmartPhoneJava/forum_bd/internal/models"
	re "github.com/SmartPhoneJava/forum_bd/internal/return_errors"
)

// CreateThread create thread
func (h *Handler) CreateThread(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateThread"
	var (
		tthread models.Thread
		thread  *models.Thread
		err     error
	)

	rw.Header().Set("Content-Type", "application/json")

	if tthread, err = getThread(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	thread = &tthread

	threadChan := make(chan *models.Thread, 1)
	errChan := make(chan error, 1)
	go h.DB.CreateThread(thread, threadChan, errChan)
	err = <-errChan
	thread = <-threadChan

	if thread == nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusNotFound, place)
		return
	} else if err != nil && err.Error() == re.ErrorThreadConflict().Error() {
		rw.WriteHeader(http.StatusConflict)
		sendSuccessJSON(rw, thread, place)
		printResult(err, http.StatusConflict, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusCreated, place)
	return
}

// UpdateThread update thread
func (h *Handler) UpdateThread(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateThread"
	var (
		thread models.Thread
		slug   string
		err    error
	)

	rw.Header().Set("Content-Type", "application/json")

	if thread, err = getThread(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, err = h.DB.UpdateThread(&thread, slug); err != nil {
		if err.Error() == re.ErrorThreadConflict().Error() {
			rw.WriteHeader(http.StatusConflict)
			sendSuccessJSON(rw, thread, place)
			printResult(err, http.StatusConflict, place)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			sendErrorJSON(rw, err, place)
			printResult(err, http.StatusNotFound, place)
		}
		return
	}

	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusOK, place)
	return
}

// GetThreadDetails get thread details
func (h *Handler) GetThreadDetails(rw http.ResponseWriter, r *http.Request) {
	const place = "GetThreadDetails"
	var (
		thread models.Thread
		slug   string
		err    error
	)

	rw.Header().Set("Content-Type", "application/json")

	if slug, err = getSlug(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if thread, err = h.DB.GetThread(slug); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, thread, place)
	printResult(err, http.StatusOK, place)
	return
}

// GetThreads get thread
func (h *Handler) GetThreads(rw http.ResponseWriter, r *http.Request) {
	const place = "GetThreads"
	var (
		threads    []models.Thread
		slug       string
		limit      int
		t          string
		err        error
		existLimit bool
		existTime  bool
		desc       bool
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

	if existTime, t, err = getTime(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	desc = getDesc(r)

	if threads, err = h.DB.GetThreads(slug, limit, existLimit, t, existTime, desc); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	sendSuccessJSON(rw, threads, place)
	printResult(err, http.StatusOK, place)
	return
}
