// datastore.go

package naprrql

import (
	"bytes"
	"log"
	"os"

	"github.com/dgraph-io/badger"
	"github.com/nsip/nias2/naprr"
)

var kv *badger.KV
var ge = naprr.GobEncoder{}

func GetKV() *badger.KV {
	return kv
}

func init() {
	kv = openDB()
}

//
// open the badger kv store
//
func openDB() *badger.KV {

	log.Println("creating data-store directory...")
	err := os.Mkdir("kvs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Println("Error trying to create datastore working directory", err)
	}

	log.Println("opening database...")
	opt := badger.DefaultOptions
	opt.Dir = "kvs"
	opt.ValueDir = "kvs"
	kv, dbErr := badger.NewKV(&opt)
	if dbErr != nil {
		log.Fatalln("DB Create error: ", dbErr)
	}
	return kv
}

//
// Given a key or key-prefix, returns the reference ids that
// can be used in a Get operation to retreive the
// desired object
//
func getObjectReferences(keyPrefix string) [][]byte {

	itrOpt := badger.IteratorOptions{
		PrefetchSize: 1000,
		FetchValues:  true,
		Reverse:      false,
	}
	itr := kv.NewIterator(itrOpt)

	searchKey := []byte(keyPrefix + ":")
	objIDs := make([][]byte, 0)
	for itr.Seek(searchKey); itr.Valid(); itr.Next() {
		if !bytes.Contains(itr.Item().Key(), searchKey) {
			break
		}
		objIDs = append(objIDs, itr.Item().Value())
	}
	return objIDs
}

func getObjectForKeys(objIDs [][]byte) ([]interface{}, error) {
	summaries := []interface{}{}

	for _, objID := range objIDs {
		// fmt.Printf("%s - %+v\n", key, val)
		var item badger.KVItem
		var summary interface{}
		if err := kv.Get(objID, &item); err != nil {
			log.Println("Cannot find object with key: ", string(objID))
			return summaries, err
		}
		if err := ge.Decode(item.Value(), &summary); err != nil {
			log.Println("Cannot decode object with key: ", objID, err)
			return summaries, err
		}
		summaries = append(summaries, summary)

	}

	return summaries, nil

}

//
//
//
