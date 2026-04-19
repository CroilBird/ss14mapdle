package models

type Session struct {
	ID          uint64    `json:"-"`
	Guid        string    `json:"guid"`
	ChallengeID uint64    `json:"-"`
	Challenge   Challenge `json:"-"`
	ZoomLevel   int       `json:"zoom_level" gorm:"default:1"`
	Correct     bool      `json:"correct"`
}

func init() {
	Models = append(Models, Session{})
}
