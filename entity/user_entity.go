package entity

type User struct {
	ID       string `gorm:"column:id"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
	Role_id  int    `gorm:"column:role_id"`
}

func (User) TableName() string {
	return "Users"
}
