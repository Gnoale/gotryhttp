package main

import (
	"fmt"
	"testing"
)

var (
	fsdb       *FileDB
	items      []*Item
	sampleKv   []map[string]string
	sampleName []string
)

func setupTesting() {
	fsdb = &FileDB{"./test/db.txt"}
	sampleKv = []map[string]string{
		{
			"forest":   "wood",
			"mountain": "rock",
			"fire":     "stone",
		},
		{
			"forest":    "limestone",
			"limestone": "climber",
			"climber":   "rope",
		},
		{
			"movie":    "pelicule",
			"aperture": "science",
			"bad":      "motherfucker",
		},
	}
	sampleName = []string{"primal", "natural", "pop"}
	for i, kv := range sampleKv {
		items = append(items, &Item{sampleName[i], kv})
	}
}

func TestInit(t *testing.T) {
	setupTesting()
	fsdb.Init()
}

func TestInsertItems(t *testing.T) {
	for _, item := range items {
		if got := fsdb.InsertItem(item); got != nil {
			t.Errorf("An error occured while inserting the item %v : %v", item.Name, got)
		}
	}
}

func TestGetItem(t *testing.T) {
	item, err := fsdb.GetItem("natural")
	if err != nil {
		t.Error("failed to get item")
	}
	fmt.Println(item.Name)
	for k, v := range item.Map {
		fmt.Printf("  %v:%v\n", k, v)
	}
	// empty name
	item, err = fsdb.GetItem("")
	if err == nil {
		t.Error("empty item name should return an error")
	}
}
