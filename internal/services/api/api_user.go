package api

import (
	data "escapade/internal/database"
	database "escapade/internal/database"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
)

// Handler is struct
type Handler struct {
	DB                    database.DataBase
	PlayersAvatarsStorage string
	FileMode              int
}

// CreateUser create new user
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateUser"
	var (
		user  models.User
		err   error
		users *[]models.User
	)

	rw.Header().Set("Content-Type", "application/json")

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user.Nickname, err = getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	users, user, err = h.DB.CreateUser(&user)

	if len(*users) > 0 {
		rw.WriteHeader(http.StatusConflict)
		printResult(err, http.StatusConflict, place)
		sendSuccessJSON(rw, users, place)
		return
	}

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		sendErrorJSON(rw, err, place)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	sendSuccessJSON(rw, user, place)
	printResult(err, http.StatusCreated, place)
	return
}

//GetUsers get users
func (h *Handler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUsers"
	var (
		users      []models.User
		slug       string
		limit      int
		since      string
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

	existSince, since = getNickNameMin(r)
	desc = getDesc(r)

	qgc.InitUser(existSince, since, existLimit, limit, desc)

	if users, err = h.DB.GetUsers(slug, qgc); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, users, place)
	printResult(err, http.StatusOK, place)
	return
}

// GetProfile returns model UserPublicInfo
func (h *Handler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "GetProfile"

	var (
		err      error
		nickname string
		user     models.User
	)

	rw.Header().Set("Content-Type", "application/json")

	if nickname, err = getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, err = h.DB.GetUser(nickname); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	sendSuccessJSON(rw, user, place)
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}

// UpdateProfile updates profile
func (h *Handler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateProfile"

	var (
		user models.User
		err  error
	)

	rw.Header().Set("Content-Type", "application/json")

	if user, err = getUser(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user.Nickname, err = getNickname(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		sendErrorJSON(rw, err, place)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	user.Print()

	if user, err = h.DB.UpdateUser(user); err != nil {
		if err.Error() == re.ErrorEmailIstaken().Error() {
			rw.WriteHeader(http.StatusConflict)
			printResult(err, http.StatusConflict, place)
		} else {
			rw.WriteHeader(http.StatusNotFound)
			printResult(err, http.StatusNotFound, place)
		}
		sendErrorJSON(rw, err, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	sendSuccessJSON(rw, user, place)
	printResult(err, http.StatusOK, place)
	return
}
