package main

import (
	"errors"
	"github.com/dgryski/go-jump"
	"strconv"
)

type Resource struct {
	Id    int
	Free  bool
	Owner string
}

type Resources struct {
	members  []Resource
	input    []chan Message
	freeList chan int
}

// Simple and fast hashring
func choose_worker(i int, n int) int {
	i64 := uint64(i)
	place := jump.Hash(i64, n)
	return int(place)
}

func (a *Resources) worker(c <-chan Message) {
	for msg := range c {
		if len(a.members) < msg.data.Id || msg.data.Id < 0 {
			msg.ch <- false
		} else {
			current := a.members[msg.data.Id].Free
			if msg.data.Free != current {
				a.members[msg.data.Id] = msg.data
				if msg.data.Free == true {
					a.freeList <- msg.data.Id
				}
				msg.ch <- true
			} else {
				msg.ch <- false
			}
		}
		close(msg.ch)
	}
}

func (a *Resources) try_allocate(name string, workers int) (int, error) {
	select {
	case i := <-a.freeList:
		output := make(chan bool)
		res := Resource{Id: i, Free: false, Owner: name}
		msg := Message{data: res, ch: output}
		place := choose_worker(i, workers)
		a.input[place] <- msg
		result := <-output
		if result == true {
			return i + 1, nil
		}
	default:
		return 0, errors.New("Failed")
	}
	return 0, errors.New("Failed")
}

func (a *Resources) try_deallocate(id int, workers int) error {
	id--
	output := make(chan bool)
	res := Resource{Id: id, Free: true, Owner: ""}
	msg := Message{data: res, ch: output}
	place := choose_worker(id, workers)
	a.input[place] <- msg
	result := <-output
	if result == false {
		return errors.New("Failed")
	}
	return nil
}

func (a *Resources) List() string {
	// My own json generator, yeah
	allocated := "{"
	deallocated := "["
	for i := range a.members {
		if a.members[i].Free {
			if deallocated != "[" {
				deallocated += ","
			}
			deallocated += "\"r" + (strconv.Itoa(a.members[i].Id + 1)) + "\""
		} else {
			if allocated != "{" {
				allocated += ","
			}
			allocated += "\"r" + (strconv.Itoa(a.members[i].Id + 1)) + "\":\"" + a.members[i].Owner + "\""
		}
	}
	allocated += "}"
	deallocated += "]"
	if allocated == "{}" {
		allocated = "[]"
	}
	return "{\"allocated\":" + allocated + "," + "\"deallocated\":" + deallocated + "}\n"
}

func (a *Resources) Init(c Config) {

	a.freeList = make(chan int, c.Limit)

	for i := 0; i < c.Limit; i++ {
		r := Resource{i, true, ""}
		a.freeList <- i
		a.members = append(a.members, r)
	}

	for i := 0; i < c.Workers; i++ {
		ch := make(chan Message, 10)
		a.input = append(a.input, ch)
		go a.worker(ch)
	}
}

func (a *Resources) Search(s string) string {
	// My own json generator, yeah
	found := "["
	for i := range a.members {
		if !a.members[i].Free && a.members[i].Owner == s {
			if found != "[" {
				found += ","
			}
			found += "\"r" + (strconv.Itoa(a.members[i].Id + 1)) + "\""
		}
	}
	found += "]\n"
	return found
}

func (a *Resources) Reset(workers int) {
	for i := range a.members {
		output := make(chan bool)
		res := Resource{Id: i, Free: true, Owner: ""}
		msg := Message{data: res, ch: output}
		place := choose_worker(i, workers)
		a.input[place] <- msg
		<-output
	}
}
