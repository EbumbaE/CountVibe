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
	GetProduct(productID string) (*map[string]string, error)
	SetProduct(insertMap map[string]string) error
	GetPortions(diaryID, date, mealOrder string) (*[]map[string]string, error)
	SetPortion(insertMap map[string]string) error
}

type Storage interface {
	UserDatabase
	DiaryDatabase
}
