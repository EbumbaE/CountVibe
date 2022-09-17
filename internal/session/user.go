package session

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
)

type AuthDetails struct {
	accessUuid string
	userID     int64
}

type User struct {
	id       int64  `yaml:"ID"`
	username string `yaml:"username"`
	password string `yaml:"password"`
	isOnline bool   `yaml:"isOnline"`
}

func (s *Session) getAuthDetails(r *http.Request) (*AuthDetails, error) {
	token, err := GetjwtToken(r, s.jwtKey["access"], "access_token")
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}

		strUserID := fmt.Sprintf("%.f", claims["user_id"])
		userID, err := strconv.ParseInt(strUserID, 10, 64)
		if err != nil {
			return nil, err
		}

		ad := &AuthDetails{
			accessUuid: accessUuid,
			userID:     userID,
		}
		return ad, nil
	}
	return nil, err
}

func (s *Session) userLogout(w http.ResponseWriter, r *http.Request) {
	deleteTokensCookie(w)

	ad, err := s.getAuthDetails(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error(err, "get authorization details")
		return
	}
	delErr := s.deleteTokens(ad)
	if delErr != nil {
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error(err, "delete tokens")
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Session) verificationUserID(authD *AuthDetails) (bool, error) {
	strUserID, err := s.tokensDB.Get(authD.accessUuid).Result()
	if err != nil {
		return false, err
	}

	uUserID, _ := strconv.ParseUint(strUserID, 10, 64)
	userID := int64(uUserID)

	if authD.userID != userID {
		return false, errors.New("userID do not mutch")
	}
	return true, nil
}

func (s *Session) newTemplate(w http.ResponseWriter, data any, paths []string) error {

	ts, err := template.ParseFiles(paths...)
	if err != nil {
		return err
	}

	err = ts.Execute(w, data)
	return err
}

func (s *Session) compareLogin(r *http.Request) (bool, error) {

	url := r.URL.Path

	pattern := `/[a-zA-Z0-9]+/`
	re, _ := regexp.Compile(pattern)
	res := re.FindAllString(url, -1)
	firstRes := res[0]
	lenUsername := len(firstRes)
	username := firstRes[1 : lenUsername-1]

	strPathUserID, err := s.db.GetUserID(username)
	if err != nil {
		return false, err
	}
	pathUserID, err := strconv.ParseInt(strPathUserID, 10, 64)

	token, err := GetjwtToken(r, s.jwtKey["access"], "access_token")
	if err != nil {
		return false, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	strUserID := fmt.Sprintf("%.f", claims["user_id"])
	userID, convErr := strconv.ParseInt(strUserID, 10, 64)
	if convErr != nil {
		return false, convErr
	}

	if userID == pathUserID {
		return true, nil
	}
	return false, nil

}

func (s *Session) userHandler(w http.ResponseWriter, r *http.Request) {

	isLogin, err := s.compareLogin(r)
	if err != nil {
		s.Logger.Error(err, "Login verification")
	}

	switch r.Method {
	case "GET":
		data := ViewUserData{
			IsLogin: isLogin,
		}
		paths := []string{
			s.paths["user"],
		}
		if err := s.newTemplate(w, data, paths); err != nil {
			s.Logger.Error(err, "new Template")
			http.Error(w, "error in create Template", http.StatusInternalServerError)
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parse form err: %v", err)
			return
		}
		button := r.FormValue("button")
		if button == "diary" {
			url := r.URL.Path + s.pages["diary"]
			http.Redirect(w, r, url, http.StatusFound)
		}
		if isLogin {
			if button == "exit" {
				s.userLogout(w, r)
				http.Redirect(w, r, "/", http.StatusFound)
			}
		}
	}
}
