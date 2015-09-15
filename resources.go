package main

import (
	//"log"
	//"github.com/bradfitz/slice"
	"fmt"
)

type Resource struct {
	id    int
	free  bool
	owner string
}

type Resources struct {
	members  []Resource
	input    []chan Message
	freeList chan int
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
		if len(a.members) < msg.data.id || msg.data.id < 0 {
			msg.ch <- false
		} else {
			current := a.members[msg.data.id].free
			if msg.data.free != current {
				a.members[msg.data.id] = msg.data
				if msg.data.free == true {
					a.freeList <- msg.data.id
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
		res := Resource{id: i, free: false, owner: name}
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
	res := Resource{id: id, free: true, owner: ""}
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
	return fmt.Sprintf("%s", a.members)
}
