### go-bench-tweet-bot

bench-tweet is a small wrapper on functionality provided by
github.ChimeraCoder.anaconda, used to allow daily posting of
nightly Go development benchmarking results.

Installation:
```go get github.com/dr2chase/go-ench-tweet-bot/cmd/bench-tweet```

Usage:
bench-tweet -i tweet.txt -m media1 -m media2 -m media3 -m media4

or 

bench-tweet -t "Some tweet text" -m media1 -m media2 -m media3 -m media4

tweet.txt and mediaN must conform to the usual twitter rules,
in this case, 280 or fewer characters, and either png or static
gif and less than 5MB in size.

If successful, returns the JSON from the posted tweet, otherwise
panics in some way with a message that may or may not be helpful.