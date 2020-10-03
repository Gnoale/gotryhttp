package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Database interface {
	GetItem(name string) (*Item, error)
	InsertItem(i *Item) error
	RemoveItem(i *Item) error
}

type FileDB struct {
	FileName string
}

type Item struct {
	Name string            `json:"name"`
	Map  map[string]string `json:"map"`
}

func NewItem() *Item {
	return &Item{
		"new",
		make(map[string]string),
	}
}

func (f *FileDB) Init() {
	// check if the file already exist
	file, err := os.Open(f.FileName)
	defer file.Close()
	if err == nil {
		return
	}
	if _, err := os.Create(f.FileName); err != nil {
		panic(err)
	}
}

func (f *FileDB) InsertItem(i *Item) error {
	if already := f.exist(i.Name); already {
		return errors.New("no insertion will be made, entry exist")
	}
	file, err := os.OpenFile(f.FileName, os.O_RDWR|os.O_APPEND, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "%v\n", i.Name); err != nil {
		return err
	}
	for k, v := range i.Map {
		if _, err := fmt.Fprintf(file, "  %v:%v\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *FileDB) GetItem(name string) (*Item, error) {
	if exist := f.exist(name); !exist {
		return nil, errors.New("the entry doesn't exist")
	}
	file, err := os.Open(f.FileName)
	defer file.Close()
	item := NewItem()
	if err != nil {
		return item, err
	}
	scanner := bufio.NewScanner(file)
	found := false
	for scanner.Scan() {
		if found && !strings.HasPrefix(scanner.Text(), "  ") {
			break
		}
		if found && strings.HasPrefix(scanner.Text(), "  ") {
			slice := strings.Fields(scanner.Text())
			slice = strings.Split(slice[0], ":")
			item.Map[slice[0]] = slice[1]
		}
		if strings.Contains(scanner.Text(), name) && !strings.HasPrefix(scanner.Text(), "  ") {
			found = true
			item.Name = name
		}
	}
	if err := scanner.Err(); err != nil {
		return item, err
	}
	return item, nil
}

func (f *FileDB) RemoveItem(i *Item) error {
	return errors.New("not implemented")
}

func (f *FileDB) exist(name string) bool {
	file, _ := os.Open(f.FileName)
	defer file.Close()
	fileInfo, _ := file.Stat()
	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)
	var builder strings.Builder
	builder.Write(buffer)
	if strings.Contains(builder.String(), name) {
		return true
	}
	return false
}

//func (f *FileDB) RemoveItem(i *Item) error {
//	file, err := os.Open(f.FileName)
//	if err != nil {
//		return nil, err
//	}
