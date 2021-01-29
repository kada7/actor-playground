/*
 * @author: Haoyuan Liu
 * @date: 2021/1/27
 */

package persis_provider

import (
	"actor-playground/core/persistence"
	"actor-playground/util"
	boltdb "github.com/artyomturkin/protoactor-go-persistence-boltdb"
	"github.com/boltdb/bolt"
	"os"
)

type BoltProviderState struct {
	filename string
	db       *bolt.DB
	persistence.ProviderState
}

func newBoltProvider(filename string) *BoltProviderState {
	db, err := bolt.Open(filename, 0666, nil)
	util.Must(err)

	return &BoltProviderState{
		filename:      filename,
		db:            db,
		ProviderState: boltdb.NewBoltProvider(3, db),
	}
}

func (b BoltProviderState) ActorNameList() []string {
	l := make([]string, 0)
	err := b.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			l = append(l, string(name))
			return nil
		})
	})
	util.Must(err)
	return l
}

func (b BoltProviderState) Close() {
	err := b.db.Close()
	util.Must(err)
}

func (b BoltProviderState) DropDB() {
	util.Must(os.Remove(b.filename))
}
