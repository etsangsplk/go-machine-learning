package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// Example of using Bolt DB to cache data on local disk.

func main() {
	db, err := bolt.Open("embedded.db", 0600, nil)
	if err != nil {
		log.Printf("Error opening Db %s\n", err.Error())
		return
	}

	defer db.Close()

	// Create bucket in db file
	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	}); err != nil {
		log.Printf("Error creating bucket %s\n", err.Error())
	}

	// Put key-values in bucket
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("myKey"), []byte("myValue"))
		return err
	}); err != nil {
		log.Printf("Error adding key-values to bucket %s\n", err.Error())
		return
	}

	// Read values from Db
	// Put key-values in bucket
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key: %s, value: %s\n", k, v)
		}
		return nil
	}); err != nil {
		log.Printf("Error adding key-values to bucket %s\n", err.Error())
		return
	}
}
