package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Menggabungkan semua kemungkinan field detail dari SRS 3.2.1
type AchievementDetails struct {
	// Competition [cite: 116-119]
	CompetitionName  string `bson:"competitionName,omitempty" json:"competitionName,omitempty"`
	CompetitionLevel string `bson:"competitionLevel,omitempty" json:"competitionLevel,omitempty"` // international, national, etc
	Rank             int    `bson:"rank,omitempty" json:"rank,omitempty"`
	MedalType        string `bson:"medalType,omitempty" json:"medalType,omitempty"` // gold, silver, bronze

	// Publication [cite: 121-125]
	PublicationType  string   `bson:"publicationType,omitempty" json:"publicationType,omitempty"`
	PublicationTitle string   `bson:"publicationTitle,omitempty" json:"publicationTitle,omitempty"`
	Authors          []string `bson:"authors,omitempty" json:"authors,omitempty"`
	Publisher        string   `bson:"publisher,omitempty" json:"publisher,omitempty"`
	ISSN             string   `bson:"issn,omitempty" json:"issn,omitempty"`

	// Organization [cite: 127-131]
	OrganizationName string    `bson:"organizationName,omitempty" json:"organizationName,omitempty"`
	Position         string    `bson:"position,omitempty" json:"position,omitempty"`
	StartDate        time.Time `bson:"startDate,omitempty" json:"startDate,omitempty"`
	EndDate          time.Time `bson:"endDate,omitempty" json:"endDate,omitempty"`

	// Certification [cite: 133-136]
	CertificationName   string    `bson:"certificationName,omitempty" json:"certificationName,omitempty"`
	IssuedBy            string    `bson:"issuedBy,omitempty" json:"issuedBy,omitempty"`
	CertificationNumber string    `bson:"certificationNumber,omitempty" json:"certificationNumber,omitempty"`
	ValidUntil          time.Time `bson:"validUntil,omitempty" json:"validUntil,omitempty"`

	// General Fields [cite: 138-141]
	EventDate time.Time `bson:"eventDate,omitempty" json:"eventDate,omitempty"`
	Location  string    `bson:"location,omitempty" json:"location,omitempty"`
	Organizer string    `bson:"organizer,omitempty" json:"organizer,omitempty"`
	Score     float64   `bson:"score,omitempty" json:"score,omitempty"`
}

type Attachment struct {
	FileName   string    `bson:"fileName" json:"fileName"`
	FileURL    string    `bson:"fileUrl" json:"fileUrl"`
	FileType   string    `bson:"fileType" json:"fileType"`
	UploadedAt time.Time `bson:"uploadedAt" json:"uploadedAt"`
}

type Achievement struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID       string             `bson:"studentId" json:"studentId"` // Disimpan sebagai string UUID dari Postgres
	AchievementType string             `bson:"achievementType" json:"achievementType"` // academic, competition, etc
	Title           string             `bson:"title" json:"title"`
	Description     string             `bson:"description" json:"description"`
	Details         AchievementDetails `bson:"details" json:"details"`
	CustomFields    map[string]interface{} `bson:"customFields,omitempty" json:"customFields,omitempty"`
	Attachments     []Attachment       `bson:"attachments" json:"attachments"`
	Tags            []string           `bson:"tags" json:"tags"`
	Points          int                `bson:"points" json:"points"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}