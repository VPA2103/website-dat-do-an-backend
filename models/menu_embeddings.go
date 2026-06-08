package models
type MenuEmbedding struct {
	ID string `gorm:"primaryKey"`

	Document string         `gorm:"type:text"`
	Metadata string         `gorm:"type:jsonb"`
	Embedding []float32     `gorm:"type:vector(3072)"`
}