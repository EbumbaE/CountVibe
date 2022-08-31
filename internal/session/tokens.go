package session

import (
    "fmt"
    "time"
    "net/http"
    "strconv"
    "errors"

    "github.com/dgrijalva/jwt-go"
    "github.com/go-redis/redis/v7"
    "github.com/twinj/uuid"
)

type Details struct{
    token string
    uuid string
    expires int64
}

type Tokens struct {
    access Details
    refresh Details
}

func (s *Session) newTokensDB(){
    ro := &redis.Options{
            Addr: "localhost:6379",
            Password: "",
            DB: 0, 
          }
    s.tokensDB = redis.NewClient(ro)
}

func (s *Session)checkHealthTokensDB()(error){
    _, err := s.tokensDB.Ping().Result()
    return err
}

func (s *Session)newAuthorization(userID int64, t *Tokens)(error){

    access := t.access
    refresh := t.refresh

	at := time.Unix(access.expires, 0)
	rt := time.Unix(refresh.expires, 0)
	now := time.Now()

	strUserID := strconv.Itoa(int(userID))
	if err := s.tokensDB.Set(access.uuid, strUserID, at.Sub(now)).Err(); err != nil {
		return err
	}

	if err := s.tokensDB.Set(refresh.uuid, strUserID, rt.Sub(now)).Err(); err != nil{
		return err
	}

	return nil
}

func (s *Session)deleteAuthorization(givenUuid string)(int64, error){

    deleted, err := s.tokensDB.Del(givenUuid).Result()
    if err != nil {
        return 0, err
    }

    return deleted, nil
}

func (s *Session)deleteTokens(ad *AuthDetails)(error){

    refreshUuid := fmt.Sprintf("%s::%d", ad.accessUuid, ad.userID)
    
    deletedAccess, errDA := s.tokensDB.Del(ad.accessUuid).Result()
    if errDA != nil {
        return errDA
    }

    deletedRefresh, errDR := s.tokensDB.Del(refreshUuid).Result()
    if errDR != nil {
        return errDR
    }

    if deletedAccess != 1 || deletedRefresh != 1 {
        return errors.New("delete tokens went wrong")
    }
    return nil

}

func (t *Tokens)newAcessDetails(userID int64, jwtKey []byte)(error){
    claims := jwt.MapClaims{}
    access := t.access

    claims["authorized"] = true
    claims["access_uuid"] = access.uuid
    claims["user_id"] = userID
    claims["expires"] = access.expires
    at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err := at.SignedString(jwtKey)

    t.access.token = token 
    return err
}

func (t *Tokens)newRefreshDetails(userID int64, jwtKey []byte)(error){
    claims := jwt.MapClaims{}
    refresh := t.refresh 
    
    claims["refresh_uuid"] = refresh.uuid
    claims["user_id"] = userID
    claims["expires"] = refresh.expires
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err := rt.SignedString(jwtKey)

    t.refresh.token = token
    return err
}

func newTokens(userID int64, jwtKey map[string][]byte)(*Tokens, error){

    access := Details{
        expires: time.Now().Add(time.Minute * 5).Unix(),
        uuid: uuid.NewV4().String(),
    }

    week := time.Hour * 24 * 7
    refresh := Details{
        expires: time.Now().Add(week).Unix(),
        uuid: access.uuid + "::" + strconv.Itoa(int(userID)),
    }

    t := &Tokens{
        access: access,
        refresh: refresh,
    }
    if err := t.newAcessDetails(userID, jwtKey["access"]); err != nil{
        return nil, err
    }
    if err := t.newRefreshDetails(userID, jwtKey["refresh"]); err != nil{
        return nil, err
    }

    return t, nil
}

func saveTokensInCookie(w http.ResponseWriter, t *Tokens){    
        
    access := (*t).access
    refresh := (*t).refresh

    http.SetCookie(w, &http.Cookie{
        Name:    "access_token",
        Value:   access.token,
        Expires: time.Unix(access.expires, 0),
        HttpOnly: true,
    })
    http.SetCookie(w, &http.Cookie{
        Name:    "refresh_token",
        Value:   refresh.token,
        Expires: time.Unix(refresh.expires, 0),
        HttpOnly: true,
    })
}

func GetjwtToken(r *http.Request, jwtKey map[string][]byte)(*jwt.Token, error){
    c, err := r.Cookie("access_token")
    if err != nil {
        return nil, err
    }

    tokenString := c.Value
    
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey["access"], nil
    })
    if err != nil {
        return nil, err
    }
    return token, nil

}