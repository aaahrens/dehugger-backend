package main

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/coreos/bbolt"
	"github.com/graph-gophers/graphql-go"
	"strconv"
)

type HugsCache struct {
	hugsWriter chan []*Website
	hugsReader chan []*Website
}

func (hugs *HugsCache) SetHugs(arrayHugs []*Website) {
	hugs.hugsWriter <- arrayHugs
}

func (hugs *HugsCache) GetHugs(first int32, last int32) []*Website {
	return <-hugs.hugsReader
}

func (hugs HugsCache) New() *HugsCache {
	h := HugsCache{
		hugsWriter: make(chan []*Website),
		hugsReader: make(chan []*Website),
	}
	go h.Handler()
	return &h
}

func (hugs *HugsCache) Handler() {

	var hugsList []*Website

	for {
		select {

		case hugsList = <-hugs.hugsWriter: // set the current value.
		case hugs.hugsReader <- hugsList: // send the current value.
		}
	}
}

func WriteDatabaseThread() {

	for {
		toInsert := <-writeChannel
		fmt.Println(toInsert.URL)

		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(WEBSITEBUCKET))

			saved := b.Get([]byte(toInsert.URL))
			savedSite := Website{}

			if err := json.Unmarshal(saved, &savedSite); err != nil {
				id, _ := b.NextSequence()
				newId := int(id)
				toInsert.Counter = 0
				toInsert.IsDown = true
				hello := graphql.ID(strconv.Itoa(newId))
				toInsert.ID = hello
				jsonString, _ := json.Marshal(toInsert)
				err := b.Put([]byte(toInsert.URL), []byte(jsonString))
				if err != nil {
					fmt.Errorf("error is %s", err)
					return err
				}
				return nil

			} else {
				savedSite.Counter++
				jsonString, _ := json.Marshal(savedSite)
				err := b.Put([]byte(toInsert.URL), []byte(jsonString))
				if err != nil {
					fmt.Errorf("error is %s", err)
					return err
				}
				return nil
			}
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func WriteUserThread(){
	for {



	}
}


func CacheReloader() {
	for {
		var newCacheArray []*Website

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(WEBSITEBUCKET))

			c := b.Cursor()

			for k, v := c.First(); k != nil; k, v = c.Next() {
				unserialized := Website{}

				if err := json.Unmarshal(v, &unserialized); err != nil {
					fmt.Println("curruption?")
					continue
				}
				newCacheArray = append(newCacheArray, &unserialized)
			}
			return nil
		})

		HugsCacher.SetHugs(newCacheArray)

		time.Sleep(time.Second * 5)
	}
}
