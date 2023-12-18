package models

type RequestP struct {
	URL					string				`json:"url"`
	CustomShort			string				`json:"short"`
}

type ResponseP struct {
	URL					string				`json:"url" gorm:"not null"`
	ShortURL			string				`json:"short" gorm:"unique;"`
}