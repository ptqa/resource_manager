# Resource manager [![Build Status](https://travis-ci.org/ptqa/resource_manager.svg?branch=master)](https://travis-ci.org/ptqa/resource_manager) [![Coverage](https://img.shields.io/coveralls/ptqa/resource_manager.svg)](https://coveralls.io/github/ptqa/resource_manager)
Implemented for fun. Original task text can be found [here](http://machinezone.ru/challenges/resource-manager) or check [local copy](https://github.com/ptqa/resource_manager/blob/master/TASK.txt). 

## Installation
    go get github.com/ptqa/resource_manager
    go install github.com/ptqa/resource_manager
    
## Running resource manager
As written in golang can be runned as single binary:

     resource_manager
    
## Configuration
Uses config.json in the same dir. Example:
```json
{
 "Port":3000,
 "Limit":5,
 "Workers":128
}
```
Where:
* Port -- HTTP port for running server
* Limit -- how many resources should we create
* Workers -- ammount of concurrent workers (for parallel allocate/deallocate)
 
## Runing tests
Resource manager uses [goconvey](http://goconvey.co/) for testing. Running and checking coverage is simple:

     go get github.com/smartystreets/goconvey
     $GOPATH/bin/goconvey
     
After opening a broser at http://127.1:8080/ you should be able to see somethin like that:
![goconvey example](goconvey.png?raw=true)

## Runing in docker
Interactive (logs will appear in stdout):

     docker run -it -e PORT=5000 -e WORKERS=3 -e LIMIT=100 -p 5000:5000 ptqa/resm

Daemon:

     docker run -d -e PORT=5000 -e WORKERS=3 -e LIMIT=100 -p 5000:5000 ptqa/resm

This will run resource manager with 100 resources, 3 workers and listening on port 5000. After that you should be able to make a query:

     $ curl 127.1:5000/list
	 {"allocated":[],"deallocated":["r1","r2","r3","r4","r5","r6","r7","r8","r9","r10","r11","r12","r13","r14","r15","r16","r17","r18","r19","r20","r21","r22","r23","r24","r25","r26","r27","r28","r29","r30","r31","r32","r33","r34","r35","r36","r37","r38","r39","r40","r41","r42","r43","r44","r45","r46","r47","r48","r49","r50","r51","r52","r53","r54","r55","r56","r57","r58","r59","r60","r61","r62","r63","r64","r65","r66","r67","r68","r69","r70","r71","r72","r73","r74","r75","r76","r77","r78","r79","r80","r81","r82","r83","r84","r85","r86","r87","r88","r89","r90","r91","r92","r93","r94","r95","r96","r97","r98","r99","r100"]}
