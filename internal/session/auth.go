package session

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (s *Session) verifyUserPass(username, password string) (bool, error) {

	if hasUser, err := s.db.CheckUsernameInDB(username); err != nil || !hasUser {
		return false, err
	}

	rightPass, err := s.db.GetUserPassword(username)
	if err != nil {
		return false, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(rightPass), []byte(password)); err == nil {
		return true, nil
	}

	return false, nil
}

func (s *Session) loginVerification(r *http.Request, jwtKey JwtKey) (bool, error) {

	token, err := GetjwtToken(r, jwtKey.access, "access_token")
	if err != nil {
		return false, err
	}

	ok, err := s.confirmToken(token, "access_uuid")

	return ok, err
}

func (s *Session) loginHandler(w http.ResponseWriter, r *http.Request) {

	ok, err := s.loginVerification(r, s.jwtKey)
	if err != nil {
		s.Logger.Error(err, "Login verification")
	}
	if ok {
		ad, err := s.getAuthDetails(r)
		if err != nil {
			s.Logger.Error(err, "get auth details")
			return
		}

		strUserID := strconv.FormatInt(ad.userID, 10)
		username, err := s.db.GetUsername(strUserID)
		if err == nil {
			url := "/" + username
			http.Redirect(w, r, url, http.StatusFound)
			return
		}
	}

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, s.paths["login"])
	case "POST":

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parse form err: %v", err)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		ok, err := s.verifyUserPass(username, password)
		if err != nil {
			s.Logger.Error(err, "verify user password")
		}
		if !ok {
			fmt.Fprintf(w, "Incorrect login or password")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		strUserID, err := s.db.GetUserID(username)
		if err != nil {
			s.Logger.Error(err, "get userID")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := strconv.ParseInt(strUserID, 10, 64)
		if err != nil {
			s.Logger.Error(err, "convertation strUserID to userID")
			w.WriteHeader(http.StatusInternalServerError)
			return

		}

		tokens, err := newTokens(userID, s.jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.newAuthorization(userID, tokens); err != nil {
			w.WriteHeader(http.StatusForbidden)
			s.Logger.Error(err, "create auth")
			return
		}

		saveTokensCookie(w, tokens)

		urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
		http.Redirect(w, r, urlProfile, http.StatusFound)
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *Session) registrationHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		http.ServeFile(w, r, s.paths["registration"])
	case "POST":

		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(w, "parse form err: %v", err)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		hasUser, err := s.db.CheckUsernameInDB(username)
		if err != nil {
			s.Logger.Error(err, "check user in db")
		}
		if hasUser {
			fmt.Fprintf(w, "User already exist")
			// todo notisfaction
			return
		}

		hash, err := hashPassword(password)
		if err != nil {
			s.Logger.Error(err, "hash password")
			return
		}

		userID := s.idGenerator.newID()
		strUserID := strconv.FormatInt(userID, 10)

		if err := s.db.InsertNewUser(strUserID, username, hash); err != nil {
			s.Logger.Error(err, "insert new user in db")
			return
		}

		tokens, err := newTokens(userID, s.jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.newAuthorization(userID, tokens); err != nil {
			w.WriteHeader(http.StatusForbidden)
			s.Logger.Error(err, "create auth")
			return
		}

		saveTokensCookie(w, tokens)

		urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
		urlDiary := fmt.Sprintf(s.formatsPages["diary"], username)
		http.HandleFunc(urlProfile, s.userHandler)
		http.HandleFunc(urlDiary, s.diaryHandler)

		http.Redirect(w, r, urlProfile, http.StatusFound)
	}
}

func (s *Session) refreshHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetjwtToken(r, s.jwtKey.refresh, "refresh_token")
	if err != nil {
		s.Logger.Error(err, "get jwt token")
	}

	ok, err := s.confirmToken(token, "refresh_uuid")
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		s.Logger.Error(err, "Confirm token while refresh")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	strUserID := fmt.Sprintf("%.f", claims["user_id"])
	userID, convErr := strconv.ParseInt(strUserID, 10, 64)
	if convErr != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		s.Logger.Error(convErr, "convertation userID")
		return
	}

	deleted, delErr := s.deleteAuthorization(refreshUuid)
	if delErr != nil || deleted == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		s.Logger.Error(delErr, "delete auth")
		return
	}

	tokens, createErr := newTokens(userID, s.jwtKey)
	if createErr != nil {
		w.WriteHeader(http.StatusForbidden)
		s.Logger.Error(createErr, "new tokens")
		return
	}

	if err := s.newAuthorization(userID, tokens); err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.Logger.Error(err, "create auth")
		return
	}

	saveTokensCookie(w, tokens)

}
