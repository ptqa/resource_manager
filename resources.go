package main

import (
	"errors"
	"github.com/dgryski/go-jump"
	"strconv"
)

// struct for sending request
// to workers
type Message struct {
	data Resource  // new resource state
	ch   chan bool // channel for getting result
}

// Basic resource to store/access/change
type Resource struct {
	Id    int
	Free  bool   // is resource owned?
	Owner string // name of owner, '' if free
}

type Resources struct {
	members  []Resource     // collection of Resources
	input    []chan Message // input channels for change requests
	freeList chan int       // freeList chan of free resources
}

// Simple and fast hashring
func chooseWorker(i int, n int) int {
	i64 := uint64(i)
	place := jump.Hash(i64, n)
	return int(place)
}

// Worker that proccess change requests
func (a *Resources) worker(c <-chan Message) {
	for msg := range c {
		if len(a.members) < msg.data.Id || msg.data.Id < 0 {
			// resource is alredy in use or wrong id
			msg.ch <- false
		} else {
			// set new resource state
			current := a.members[msg.data.Id].Free
			if msg.data.Free != current {
				a.members[msg.data.Id] = msg.data
				if msg.data.Free == true {
					// new state is free? add it to free queue
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

// Take free resource from queue and send new resource state (free)
// to worker, then check result from channel
func (a *Resources) tryAllocate(name string, workers int) (int, error) {
	select {
	case i := <-a.freeList:
		output := make(chan bool)
		res := Resource{Id: i, Free: false, Owner: name}
		msg := Message{data: res, ch: output}
		place := chooseWorker(i, workers)
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

// Send resuurce state (busy) to worker and check result from channel
func (a *Resources) tryDeallocate(id int, workers int) error {
	output := make(chan bool)
	res := Resource{Id: id - 1, Free: true, Owner: ""}
	msg := Message{data: res, ch: output}
	place := chooseWorker(id-1, workers)
	a.input[place] <- msg
	result := <-output
	if result == false {
		return errors.New("Failed")
	}
	return nil
}

// shows all resources as JSON
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

// Create queue of free resources and create workers
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

// Search for resources with name == s
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

// Send request to free all resources
func (a *Resources) Reset(workers int) {
	for i := range a.members {
		output := make(chan bool)
		res := Resource{Id: i, Free: true, Owner: ""}
		msg := Message{data: res, ch: output}
		place := chooseWorker(i, workers)
		a.input[place] <- msg
		<-output
	}
}
