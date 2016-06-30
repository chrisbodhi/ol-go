#Businesses API

###Giving my first Golang project a **go**.

![take your pick -- imgs based on work by Renee French](https://blog.golang.org/gopher/usergroups.png)

_Possible downside of Go's name: lazy jokes are too easy!_

###Setup -- _untested_
- Clone this repo somewhere within your Gopath directory
- Run `go get` to install the dependencies
- Add db.sqlite3 to the root directory [copy it from another location, like the output of the `rake setup` task from chrisbodhi/ol]
- After that has installed, run `go run server.go` to start the server on port 8080
- Visit [localhost:8080/api/businesses](http://localhost:8080/api/businesses) to get the first 50 results returned in JSON format.
- Visit [localhost:8080/api/businesses?page=3](http://localhost:8080/api/businesses?page=3) to get the third 50 results returned in JSON format.
- Visit [localhost:8080/api/businesses/42](http://localhost:8080/api/businesses/42) to get that business's entry returned.

###todo
- ~~add a db and a way to talk to it _[oh, is that all?]_~~ _completed Wed Jun 29 17:05_
- ~~display a single result from the db~~ _completed Thurs Jun 30 0:33_
- ~~paginate the results~~ _completed Thurs Jun 30 1:54_
- add metadata
- add function to create the db from the input CSV

![one day, i'll be as well-balanced with go](https://media.giphy.com/media/67rLnLxPbC7VS/giphy.gif)
