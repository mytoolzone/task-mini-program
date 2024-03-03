package entity

// UserSetting 用户的基本信息
// 年龄 ,学历 , 工作, 个人简介, 个人标签, 个人头像, 个人背景图 , 保险信息, 个人身份证信息, 体重, 身高
//type UserSetting struct {
//	Birthday   string `gorm:"column:birthday" json:"birthday"`
//	Education  string `gorm:"column:education" json:"education"`
//	Job        string `gorm:"column:job" json:"job"`
//	Describe   string `gorm:"column:describe" json:"describe"`
//	Tags       string `gorm:"column:tags" json:"tags"`
//	Avatar     string `gorm:"column:avatar" json:"avatar"`
//	Background string `gorm:"column:background" json:"background"`
//	Insurance  string `gorm:"column:insurance" json:"insurance"`
//	Card       string `gorm:"column:card" json:"card"`
//	Weight     string `gorm:"column:weight" json:"weight"`
//	Height     string `gorm:"column:height" json:"height"`
//	Phone      string `gorm:"column:phone" json:"phone"`
//	Email      string `gorm:"column:email" json:"email"`
//}

// UserSetting 用户的基本信息
type UserSetting struct {
	UserID                string      `json:"user_id"`
	IntroUserID           string      `json:"intro_user_id"`
	Name                  string      `json:"username"`
	LoginName             string      `json:"login_name"`
	Sex                   string      `json:"sex"`
	Phone                 string      `json:"phone"`
	Birthday              string      `json:"birthday"`
	WechatName            string      `json:"wechat_name"`
	Married               string      `json:"married"`
	IDCard                string      `json:"idcard"`
	Education             string      `json:"education"`
	Email                 string      `json:"email"`
	Mingzu                string      `json:"mingzu"`
	Region                string      `json:"region"`
	Address               string      `json:"address"`
	EmergencyContact      string      `json:"emergency_contact"`
	EmergencyPhone        string      `json:"emergency_phone"`
	EmergencyRelationship string      `json:"emergency_relationship"`
	InsuranceStart        string      `json:"insurance_start"`
	InsuranceEnd          string      `json:"insurance_end"`
	Intro                 string      `json:"intro"`
	WorkCategory          string      `json:"work_category"`
	InsuranceName         string      `json:"insurance_name"`
	InsurancePhoto        string      `json:"insurance_photo"`
	Insurances            []Insurance `json:"insurances"`
}
