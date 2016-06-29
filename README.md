#Businesses API

###Giving my first Golang project a **go**.

![take your pick -- imgs based on work by Renee French](https://blog.golang.org/gopher/usergroups.png)

_Possible downside of Go's name: lazy jokes are too easy!_

###Setup -- _untested_
- Clone the repo somewhere within your Gopath directory
- Run `go get` to install the single [at this time] dependency
- After that has installed, run `go run server.go` to start the server on port 8080
- Visit [localhost:8080/api/businesses](http://localhost:8080/api/businesses) to get some static JSON returned.
- Visit [localhost:8080/api/businesses/42](http://localhost:8080/api/businesses/42) to get your param returned.

###todo
- add a db and a way to talk to it _[oh, it that all?]_

![one day, i'll be well-balanced with go](https://media.giphy.com/media/67rLnLxPbC7VS/giphy.gif)
