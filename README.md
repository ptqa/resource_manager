# Resource manager
Implemented for fun. Original task text can be found [here](http://machinezone.ru/challenges/resource-manager) or check [local copy](https://github.com/ptqa/resource_manager/blob/master/TASK.txt).

## Installation
    go get github.com/ptqa/resource_manager
    go get github.com/ptqa/resource_manager
    
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

TODO
