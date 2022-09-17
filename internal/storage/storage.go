package storage

type UserDatabase interface {
	InsertNewUser(id, username, password string) error
	GetUsername(userID string) (string, error)
	GetUserID(username string) (string, error)
	GetLastUserID() (string, error)
	GetUserPassword(username string) (string, error)
	CheckUsernameInDB(username string) (bool, error)
	DeleteUser(username string) error
	GetAllUsernames() (chan string, error)
}

type DiaryDatabase interface {
	AddItem(date string, id int64, amount float64, od int64) error
}

type Storage interface {
	UserDatabase
	DiaryDatabase
}
