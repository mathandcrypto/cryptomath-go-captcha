package models

type CaptchaTask struct {
	Uuid	string	`gorm:"column:uuid;primary_key;not_null"`
	Index	int	`gorm:"column:index;index;not_null"`
	Difficulty int16	`gorm:"column:difficulty;not_null"`
	Math	string	`gorm:"column:math;not_null"`
	Answer	int	`gorm:"column:answer;not_null"`
}

func (CaptchaTask) TableName() string {
	return "captcha_tasks"
}
