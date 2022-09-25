package session

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-redis/redis/v7"

	"github.com/EbumbaE/CountVibe/internal/logger"
	"github.com/EbumbaE/CountVibe/internal/storage"
)

type JwtKey struct {
	access  []byte
	refresh []byte
}

type Session struct {
	pages        map[string]string
	paths        map[string]string
	formatsPages map[string]string
	jwtKey       JwtKey

	idGenerator IDGenerator
	Logger      logger.Logger
	tokensDB    *redis.Client
	db          storage.Storage
}

func NewSession(c Config, confpages map[string]string, db storage.Storage, logger logger.Logger) *Session {
	return &Session{
		pages:        confpages,
		paths:        c.Paths,
		formatsPages: c.FormatsPages,
		jwtKey: JwtKey{
			access:  []byte(c.JwtKey.Access),
			refresh: []byte(c.JwtKey.Refresh),
		},

		idGenerator: IDGenerator{id: 0},
		Logger:      logger,
		db:          db,
	}
}

func (s *Session) setupDefaultHandlers() {
	http.HandleFunc(s.pages["login"], s.loginHandler)
	http.HandleFunc(s.pages["registration"], s.registrationHandler)
	http.HandleFunc(s.pages["refresh"], s.refreshHandler)
}

func (s *Session) restoreUserHandlers() {

	usernamesChan, err := s.db.GetAllUsernames()
	if err != nil {
		s.Logger.Error(err, "get all usernames")
		return
	}

	for username := range usernamesChan {
		urlProfile := fmt.Sprintf(s.formatsPages["profile"], username)
		urlDiary := fmt.Sprintf(s.formatsPages["diary"], username)

		http.HandleFunc(urlProfile, s.userHandler)
		http.HandleFunc(urlDiary, s.diaryHandler)
	}
}

func (s *Session) setupSessionHandlers() {
	s.setupDefaultHandlers()
	s.restoreUserHandlers()
}

func (s *Session) setupIdGenerator() {
	strLastID, err := s.db.GetLastUserID()
	if err != nil {
		s.Logger.Process(err, "Get last userID")
		return
	}
	lastID, convErr := strconv.ParseInt(strLastID, 10, 64)
	if convErr != nil {
		s.Logger.Error(convErr, "convertation last userID")
		return
	}
	s.idGenerator.setID(lastID)
}

func (s *Session) parseUsernameFromURL(r *http.Request) string {

	url := r.URL.Path
	pattern := `/[a-zA-Z0-9]+/`
	re, _ := regexp.Compile(pattern)
	res := re.FindAllString(url, -1)
	firstRes := res[0]
	lenUsername := len(firstRes)

	return firstRes[1 : lenUsername-1]
}

func (s *Session) Run() {
	s.newTokensDB()
	if err := s.checkHealthTokensDB(); err != nil {
		s.Logger.Error(err, "check health TokensDB")
		return
	}
	s.setupSessionHandlers()
	s.setupIdGenerator()
}
