package main

import (
	//"log"
	//"github.com/bradfitz/slice"
	"fmt"
	"github.com/dgryski/go-jump"
	"strconv"
)

type Resource struct {
	Id    int    `json:"id"`
	Free  bool   `json:"free"`
	Owner string `json:"owner"`
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

/*
func (a *aresource) Less(i, j int) bool {
	if a.members[i].free == true && a.members[j].free == false {
		return true
	}
	return false
}
*/

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

func (a *Resources) try_allocate(name string, workers int) (int, string) {
	var i, httpStatus int
	var httpMsg string
	select {
	case i = <-a.freeList:
		output := make(chan bool)
		res := Resource{Id: i, Free: false, Owner: name}
		msg := Message{data: res, ch: output}
		place := choose_worker(i, workers)
		a.input[place] <- msg
		result := <-output
		if result == true {
			httpStatus = 200
			httpMsg = fmt.Sprintf("r%d\n", i+1)
		}
	default:
		httpStatus = 503
		httpMsg = "Out of resources.\n"
	}
	return httpStatus, httpMsg
}

func (a *Resources) try_deallocate(id int, workers int) (int, string) {
	id--
	var httpStatus int
	var httpMsg string
	output := make(chan bool)
	res := Resource{Id: id, Free: true, Owner: ""}
	msg := Message{data: res, ch: output}
	place := choose_worker(id, workers)
	a.input[place] <- msg
	result := <-output
	if result == true {
		httpStatus = 204
		httpMsg = ""
	} else {
		httpMsg = "Not allocated\n"
		httpStatus = 404
	}
	return httpStatus, httpMsg
}

func (a *Resources) list() string {
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

func (a *Resources) search(s string) string {
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
