# aleatory
_[a·le·a·to·ry] is a musical description, referring music composed with elements of random choice_

This project attempts to create music from events streaming in Twitter.  Inspiration for the project comes from work done at http://www.bitlisten.com/ and https://github.com/soulwire/sketch.js

Harp sounds created using http://www.audacityteam.org/
Bird sounds from http://www.orangefreesounds.com/free-birds-chirping-sound-effect/
Ambient pad from http://www.dl-sounds.com/royalty-free/summer-ambient/

## Getting Started
This project requires Go v1.6+

```
# get the source
cd /path/to/your/go/projects
go get github.com/mixmastermike/aleatory
# ensure the dependancies are present
cd src/github.com/mixmastermike/aleatory
glide install
```

```
# to run tests
go test $(go list ./... | grep -v '/vendor/')
```

```
# to run
go build -o aleatory github.com/mixmastermike/aleatory/app
./aleatory -consumer-key={aa} -consumer-secret={bb} -access-token={cc} - access-secret={dd}
```
