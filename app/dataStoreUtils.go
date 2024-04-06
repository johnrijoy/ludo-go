package app

import (
	"time"

	"github.com/ostafen/clover/v2/document"
)

type AudioListCriteria uint8

const (
	RecentlyPlayed AudioListCriteria = iota
	MostPlayed
	MostLikes
)

type audioDoc struct {
	AudioBasic
	Likes     int
	PlayCount int
	LastPlay  time.Time
}

func NewaudioDoc(aud AudioBasic) audioDoc {
	return audioDoc{
		AudioBasic: aud,
		LastPlay:   time.Now(),
		PlayCount:  1,
	}
}

func (audDoc *audioDoc) getDocument() *document.Document {
	return document.NewDocumentOf(audDoc)
}

func GetaudioDoc(doc *document.Document, err error) (*audioDoc, error) {
	if err != nil {
		return nil, err
	}
	audDoc := &audioDoc{}
	err = doc.Unmarshal(&audDoc)
	return audDoc, err
}

func GetaudioDocList(docs []*document.Document, err error) ([]*audioDoc, error) {
	if err != nil {
		return nil, err
	}
	audioDocs := make([]*audioDoc, len(docs))
	for i, doc := range docs {
		audDoc := &audioDoc{}
		err = doc.Unmarshal(audDoc)
		if err != nil {
			return nil, err
		}
		audioDocs[i] = audDoc
	}
	return audioDocs, err
}
