package entity

// UserSetting 用户的基本信息
// 年龄 ,学历 , 工作, 个人简介, 个人标签, 个人头像, 个人背景图 , 保险信息, 个人身份证信息, 体重, 身高
type UserSetting struct {
	Birthday   string `gorm:"column:birthday" json:"birthday"`
	Education  string `gorm:"column:education" json:"education"`
	Job        string `gorm:"column:job" json:"job"`
	Describe   string `gorm:"column:describe" json:"describe"`
	Tags       string `gorm:"column:tags" json:"tags"`
	Avatar     string `gorm:"column:avatar" json:"avatar"`
	Background string `gorm:"column:background" json:"background"`
	Insurance  string `gorm:"column:insurance" json:"insurance"`
	Card       string `gorm:"column:card" json:"card"`
	Weight     string `gorm:"column:weight" json:"weight"`
	Height     string `gorm:"column:height" json:"height"`
	Phone      string `gorm:"column:phone" json:"phone"`
	Email      string `gorm:"column:email" json:"email"`
}
