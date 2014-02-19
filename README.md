# twitch-plays

Inspired by [TwitchPlaysPokemon](http://www.twitch.tv/twitchplayspokemon), twitch-plays sends input from the TwitchTV chat to the [DeSmuME](http://desmume.org/) window.

![input](http://i.imgur.com/aatUPTc.png)

## Installation

* Install [Go](http://golang.org/)
* Run the following command `go get github.com/vevix/twitch-plays`
* Navigate to `$GOPATH\src\github.com\vevix\twitch-plays` and rename `config/app.conf.example` to `config/app.conf` and fill in your own details
* Make sure you have `DeSmuME` open and ready to go
* run `go run main.go` or alternatively run `go build` then `./twitch-plays.exe`
