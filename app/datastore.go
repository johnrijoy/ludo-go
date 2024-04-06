package app

import (
	"io"
	"log"
	"time"

	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

var dbLog = log.New(io.Discard, "audioDb: ", log.LstdFlags|log.Lmsgprefix)

type AudioDatastore struct {
	db *clover.DB
}

const (
	audioDocCollection = "audioDocs"
	playListCollection = "playlists"
)

func (adb *AudioDatastore) InitDb(path string) error {
	db, err := clover.Open(path)
	if err != nil {
		return err
	}
	adb.db = db

	if ok, _ := db.HasCollection(audioDocCollection); !ok {
		db.CreateCollection(audioDocCollection)
	}

	return nil
}

func (adb *AudioDatastore) CloseDb() error {
	if adb.db == nil {
		return nil
	}
	return adb.db.Close()
}

// Audio Doc collection

func (adb *AudioDatastore) GetaudioDoc(ytId string) (*audioDoc, error) {
	audDoc, err := GetaudioDoc(adb.db.FindFirst(query.NewQuery(audioDocCollection).Where(query.Field("YtId").Eq(ytId))))
	if err != nil {
		return nil, err
	}
	return audDoc, nil
}

func (adb *AudioDatastore) IsAudioDocExist(ytId string) (bool, error) {
	ok, err := adb.db.Exists(query.NewQuery(audioDocCollection).Where(query.Field("YtId").Eq(ytId)))
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (adb *AudioDatastore) SaveAudioDoc(aud AudioBasic) error {
	audioDoc := NewaudioDoc(aud)
	_, err := adb.db.InsertOne(audioDocCollection, audioDoc.getDocument())
	return err
}

func (adb *AudioDatastore) UpdateListened(ytId string) error {
	doc, err := adb.db.FindFirst(query.NewQuery(audioDocCollection).Where(query.Field("YtId").Eq(ytId)))
	if err != nil {
		return err
	}

	id := doc.ObjectId()

	currPlayed, ok := doc.Get("PlayCount").(int64)
	if !ok {
		currPlayed = 0
	}
	played := currPlayed + 1

	err = adb.db.UpdateById(audioDocCollection, id, func(doc *document.Document) *document.Document {
		dbLog.Println("Before:", doc.AsMap())
		doc.Set("PlayCount", played)
		doc.Set("LastPlay", time.Now())
		dbLog.Println("After:", doc.AsMap())
		return doc
	})

	dbLog.Println("error:", err)

	return err
}

func (adb *AudioDatastore) SaveOrIncrementAudioDoc(aud AudioBasic) error {
	ok, err := adb.IsAudioDocExist(aud.YtId)
	if err != nil {
		return err
	}
	if ok {
		return adb.UpdateListened(aud.YtId)
	}
	return adb.SaveAudioDoc(aud)
}

func (adb *AudioDatastore) UpdateLikes(ytId string) error {
	doc, err := adb.db.FindFirst(query.NewQuery(audioDocCollection).Where(query.Field("YtId").Eq(ytId)))
	if err != nil {
		return err
	}

	id := doc.ObjectId()
	currLikes, ok := doc.Get("Likes").(int64)
	if !ok {
		currLikes = 0
	}
	likes := currLikes + 1

	err = adb.db.UpdateById(audioDocCollection, id, func(doc *document.Document) *document.Document {
		doc.Set("Likes", likes)
		return doc
	})

	return err
}

func (adb *AudioDatastore) GetAudioList(crit AudioListCriteria, offset int, limit int) ([]*audioDoc, error) {
	q := query.NewQuery(audioDocCollection)

	if crit == RecentlyPlayed {
		q = q.Sort(query.SortOption{
			Field: "LastPlay", Direction: -1,
		})
	} else if crit == MostPlayed {
		q = q.Sort(query.SortOption{
			Field: "PlayCount", Direction: -1,
		})
	} else if crit == MostLikes {
		q = q.Where(query.Field("Likes").Gt(0)).Sort(query.SortOption{
			Field: "LastPlay", Direction: -1,
		})
	}

	q = q.Skip(offset).Limit(limit)

	audDocs, err := GetaudioDocList(adb.db.FindAll(q))
	if err != nil {
		return nil, err
	}
	return audDocs, nil
}
