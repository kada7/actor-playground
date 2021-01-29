/*
 * @author: Haoyuan Liu
 * @date: 2021/1/27
 */

package persis_provider

import (
	"actor-playground/util"
	boltdb "github.com/artyomturkin/protoactor-go-persistence-boltdb"
	"github.com/boltdb/bolt"
	"os"
)

const boltFile = "my.db"

var (
	db *bolt.DB
)

func newBolt() *Provider {
	var err error
	db, err = bolt.Open(boltFile, 0666, nil)
	util.Must(err)
	return NewProvider(boltdb.NewBoltProvider(3, db))
}

func CloseDB() {
	err := db.Close()
	util.Must(err)
}

func CleanBolt() {
	util.Must(os.Remove(boltFile))
	util.Must(os.Remove(boltFile + ".lock"))
}

func GetAllKeys() ([]string, error) {
	l := make([]string, 0)
	err := db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			l = append(l, string(name))
			return nil
		})
	})
	return l, err
}
