// Package stow is used to persist objects to a bolt.DB database.
package stow

import (
	"bytes"
	"errors"
	"sync"

	"github.com/boltdb/bolt"
)

var pool = &sync.Pool{
	New: func() interface{} { return bytes.NewBuffer(nil) },
}

// ErrNotFound indicates object is not in database.
var ErrNotFound = errors.New("not found")

// Store manages objects persistence.
type Store struct {
	db     *bolt.DB
	bucket bucketSpec
	codec  Codec
}

// NewStore creates a new Store, using the underlying
// bolt.DB "bucket" to persist objects.
// NewStore uses GobEncoding, your objects must be registered
// via gob.Register() for this encoding to work.
func NewStore(db *bolt.DB, bucket []byte) *Store {
	return NewCustomStore(db, bucket, GobCodec{})
}

// NewJSONStore creates a new Store, using the underlying
// bolt.DB "bucket" to persist objects as json.
func NewJSONStore(db *bolt.DB, bucket []byte) *Store {
	return NewCustomStore(db, bucket, JSONCodec{})
}

// NewXMLStore creates a new Store, using the underlying
// bolt.DB "bucket" to persist objects as xml.
func NewXMLStore(db *bolt.DB, bucket []byte) *Store {
	return NewCustomStore(db, bucket, XMLCodec{})
}

// NewCustomStore allows you to create a store with
// a custom underlying Encoding
func NewCustomStore(db *bolt.DB, bucket []byte, codec Codec) *Store {
	return &Store{db: db, bucket: bucketSpec{bucket}, codec: codec}
}

// NewNestedStore returns a new Store which is nested inside the current store's
// bucket. It inherits the original store's Codec, and will be deleted by the parent
// store's DeleteAll method. Also note that buckets are in the parents key-space so
// you cannot have a NestedStore whose "bucket" is the same as a parent's key.
func (s *Store) NewNestedStore(bucket []byte) *Store {
	return s.NewCustomNestedStore(bucket, s.codec)
}

// NewCustomNestedStore works the same as NewNestedStore except you can override the
// Codec used by the returned Store.
func (s *Store) NewCustomNestedStore(bucket []byte, codec Codec) *Store {
	return &Store{
		db:     s.db,
		bucket: append(s.bucket, bucket),
		codec:  codec,
	}
}

func (s *Store) marshal(val interface{}) (data []byte, err error) {
	buf := pool.Get().(*bytes.Buffer)
	err = s.codec.NewEncoder(buf).Encode(val)
	data = append(data, buf.Bytes()...)
	buf.Reset()
	pool.Put(buf)

	return data, err
}

func (s *Store) unmarshal(data []byte, val interface{}) (err error) {
	return s.codec.NewDecoder(bytes.NewReader(data)).Decode(val)
}

func (s *Store) toBytes(key interface{}) (keyBytes []byte, err error) {
	switch k := key.(type) {
	case string:
		return []byte(k), nil
	case []byte:
		return k, nil
	default:
		return s.marshal(key)
	}
}

// Put will store b with key "key". If key is []byte or string it uses the key
// directly. Otherwise, it marshals the given type into bytes using the stores Encoder.
func (s *Store) Put(key interface{}, b interface{}) error {
	keyBytes, err := s.toBytes(key)
	if err != nil {
		return err
	}
	return s.put(keyBytes, b)
}

// Put will store b with key "key". If key is []byte or string it uses the key
// directly. Otherwise, it marshals the given type into bytes using the stores Encoder.
func (s *Store) put(key []byte, b interface{}) (err error) {
	var data []byte
	data, err = s.marshal(b)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bolt.Tx) error {
		objects, err := s.bucket.createOrGet(tx)
		if err != nil {
			return err
		}
		return objects.Put(key, data)
	})
}

// Pull will retrieve b with key "key", and removes it from the store.
func (s *Store) Pull(key interface{}, b interface{}) error {
	keyBytes, err := s.toBytes(key)
	if err != nil {
		return err
	}
	return s.pull(keyBytes, b)
}

// Pull will retrieve b with key "key", and removes it from the store.
func (s *Store) pull(key []byte, b interface{}) error {
	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	err := s.db.Update(func(tx *bolt.Tx) error {
		objects := s.bucket.get(tx)
		if objects == nil {
			return ErrNotFound
		}

		data := objects.Get(key)
		if data == nil {
			return ErrNotFound
		}

		buf.Write(data)
		return objects.Delete(key)
	})

	if err != nil {
		return err
	}

	return s.unmarshal(buf.Bytes(), b)
}

// Get will retrieve b with key "key". If key is []byte or string it uses the key
// directly. Otherwise, it marshals the given type into bytes using the stores Encoder.
func (s *Store) Get(key interface{}, b interface{}) error {
	keyBytes, err := s.toBytes(key)
	if err != nil {
		return err
	}
	return s.get(keyBytes, b)
}

// Get will retrieve b with key "key"
func (s *Store) get(key []byte, b interface{}) error {
	buf := bytes.NewBuffer(nil)
	err := s.db.View(func(tx *bolt.Tx) error {
		objects := s.bucket.get(tx)
		if objects == nil {
			return ErrNotFound
		}
		data := objects.Get(key)
		if data == nil {
			return ErrNotFound
		}
		buf.Write(data)
		return nil
	})

	if err != nil {
		return err
	}

	return s.unmarshal(buf.Bytes(), b)
}

// ForEach will run do on each object in the store.
// do can be a function which takes either: 1 param which will take on each "value"
// or 2 params where the first param is the "key" and the second is the "value".
func (s *Store) ForEach(do interface{}) error {
	fc, err := newFuncCall(s, do)
	if err != nil {
		return err
	}

	return s.db.View(func(tx *bolt.Tx) error {
		objects := s.bucket.get(tx)
		if objects == nil {
			return nil
		}
		return objects.ForEach(fc.call)
	})
}

// DeleteAll empties the store
func (s *Store) DeleteAll() error {
	return s.db.Update(s.bucket.delete)
}

// Delete will remove the item with the specified key from the store.
// It returns nil if the item was not found (like BoltDB).
func (s *Store) Delete(key interface{}) error {
	keyBytes, err := s.toBytes(key)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		objects := s.bucket.get(tx)
		if objects == nil {
			return nil
		}
		return objects.Delete(keyBytes)
	})
}

type bucketSpec [][]byte

func (bs bucketSpec) get(tx *bolt.Tx) (bt *bolt.Bucket) {
	for _, b := range bs {
		if bt != nil {
			bt = bt.Bucket(b)
		} else {
			bt = tx.Bucket(b)
		}

		if bt == nil {
			break
		}
	}
	return bt
}

func (bs bucketSpec) createOrGet(tx *bolt.Tx) (bt *bolt.Bucket, err error) {
	for _, b := range bs {
		if bt != nil {
			bt, err = bt.CreateBucketIfNotExists(b)
		} else {
			bt, err = tx.CreateBucketIfNotExists(b)
		}

		if bt == nil || err != nil {
			break
		}
	}
	return bt, err
}

func (bs bucketSpec) delete(tx *bolt.Tx) (err error) {
	switch len(bs) {
	case 0:
		return nil
	case 1:
		return tx.DeleteBucket(bs[0])
	default:
		lastParentBucket := bs[:len(bs)-1]
		childBucketName := bs[len(bs)-1]
		return lastParentBucket.get(tx).DeleteBucket(childBucketName)
	}
}
