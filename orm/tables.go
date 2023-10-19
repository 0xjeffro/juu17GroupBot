package orm

type GroupMessage struct {
	MessageID     int    `gorm:"index"`
	UserID        int64  `gorm:"index"`
	UserName      string `gorm:"index"`
	UserFirstName string `gorm:"index"`
	UserLastName  string `gorm:"index"`
	Text          string
	Date          int `gorm:"index"`
}

func CreateTable() {
	err := Init()
	if err != nil {
		panic(err)
	}

	db := GetConn()
	err = db.AutoMigrate(&GroupMessage{})
	if err != nil {
		panic(err)
	}
}

func main() {
	CreateTable()
}
