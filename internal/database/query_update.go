package database

import (
	"strings"

	//
	_ "github.com/lib/pq"
)

func queryUpdateString(name string, value string) string {

	if strings.Trim(value, " ") == "" {
		return ""
	}
	return ` ` + name + ` = '` + value + `' `
}

//
func queryUpdateMessage(value string) string {

	return queryUpdateString("message", value)
}

func queryUpdateTitle(value string) string {

	return queryUpdateString("title", value)
}

func queryUpdateFullname(value string) string {

	return queryUpdateString("fullname", value)
}

func queryUpdateEmail(value string) string {

	return queryUpdateString("email", value)
}

func queryUpdateAbout(value string) string {

	return queryUpdateString("about", value)
}

func queryUpdateThread(message string, title string) (query string) {
	str1 := queryUpdateMessage(message)
	str2 := queryUpdateTitle(title)
	if str1 == "" && str2 == "" {
		return
	}
	if str1 != "" && str2 != "" {
		str1 += ","
	}
	query = `	UPDATE Thread set ` + str1 + str2
	return
}

func queryUpdatePost(message string) (query string) {
	str1 := queryUpdateMessage(message)
	if str1 == "" {
		return
	}
	query = `	UPDATE Post set ` + str1 + `, isEdited=true`
	return
}

func queryUpdateUser(fullname string, email string,
	about string) (query string) {
	str1 := queryUpdateFullname(fullname)
	str2 := queryUpdateEmail(email)
	str3 := queryUpdateAbout(about)
	if str1 == "" && str2 == "" && str3 == "" {
		return
	}

	query = `	UPDATE UserForum set `
	count := 0
	insertStr(&query, str1, &count)
	insertStr(&query, str2, &count)
	insertStr(&query, str3, &count)
	return
}

func insertStr(query *string, add string, count *int) {
	if add == "" {
		return
	}
	if *count > 0 {
		*query += "," + add
	} else {
		*query += add
	}
	*count++
}
