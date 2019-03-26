package api

import (
	"encoding/json"
	"escapade/internal/models"
	re "escapade/internal/return_errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getStringFromPath(r *http.Request, name string,
	expected error) (str string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)

	if str = vars[name]; str == "" {
		err = expected
		return
	}
	return
}

func getIntFromPath(r *http.Request, name string,
	expected error) (val int, err error) {
	var (
		vars map[string]string
		str  string
	)

	vars = mux.Vars(r)

	if str = vars[name]; str == "" {
		err = expected
		return
	}
	if val, err = strconv.Atoi(str); err != nil {
		err = expected
		return
	}
	if val < 0 {
		err = expected
		return
	}
	return
}

func getNickname(r *http.Request) (nickname string, err error) {
	return getStringFromPath(r, "name", re.ErrorInvalidName())
}

func getPostID(r *http.Request) (id int, err error) {
	return getIntFromPath(r, "id", re.ErrorInvalidID())
}

func getSlug(r *http.Request) (slug string, err error) {
	return getStringFromPath(r, "slug", re.ErrorForumSlugInvalid())
}

func getUser(r *http.Request) (user models.User, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()

		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}

func getLimit(r *http.Request) (exist bool, limit int, err error) {
	exist = true
	str := r.FormValue("limit")
	if str == "" {
		exist = false
		return
	}
	if limit, err = strconv.Atoi(str); err != nil {
		return
	}
	if limit < 0 {
		err = re.ErrorInvalidLimit()
		return
	}

	return
}

func getNickNameMin(r *http.Request) (exist bool, str string) {
	exist = true
	str = r.FormValue("since")
	if str == "" {
		exist = false
	}
	return
}

func getRelated(r *http.Request) (exist bool, str string) {
	exist = true
	str = r.FormValue("related")
	if str == "" {
		exist = false
	}
	return
}

func getIDmin(r *http.Request) (exist bool, since int, err error) {
	exist = true
	str := r.FormValue("since")
	if str == "" {
		exist = false
		return
	}

	if since, err = strconv.Atoi(str); err != nil {
		return
	}
	if since < 0 {
		err = re.ErrorInvalidDate()
		return
	}

	return
}

func getTime(r *http.Request) (exist bool, t time.Time, err error) {

	exist = true
	str := r.FormValue("since")
	if str == "" {
		exist = false
		return
	}

	var num int

	if t, err = time.Parse("2006-01-02T15:04:05.000+03:00", str); err != nil {
		if t, err = time.Parse("2006-01-02T15:04:05.000Z", str); err != nil {
			if num, err = strconv.Atoi(str); err != nil {
				return
			}
			if num < 0 {
				err = re.ErrorInvalidDate()
				return
			} else if num < 10000 {
				t = time.Date(num, 0, 0, 0, 0, 0, 0, time.UTC)
			} else {
				err = re.ErrorInvalidDate()
				return
			}

		}

	}

	return
}

func getDesc(r *http.Request) (desc bool) {

	str := r.FormValue("desc")
	if str == "" {
		desc = false
		return
	}
	if str == "true" {
		desc = true
	}
	return
}

func getSort(r *http.Request) string {

	return r.FormValue("sort")
}

func getForum(r *http.Request) (user models.Forum, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)
	return
}

func getThread(r *http.Request) (user models.Thread, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&user)

	return
}

func getPosts(r *http.Request) (posts []models.Post, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&posts)

	return
}

func getPost(r *http.Request) (post models.Post, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&post)

	return
}

func getVote(r *http.Request) (vote models.Vote, err error) {

	if r.Body == nil {
		err = re.ErrorNoBody()
		return
	}
	defer r.Body.Close()

	_ = json.NewDecoder(r.Body).Decode(&vote)

	return
}

// 250
