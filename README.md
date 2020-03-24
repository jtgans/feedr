# FeedR

FeedR (pronounced feeder) is a royally stupid RSS generation and storage
microservice designed for the author's weird use case of scraping web snippets
via kibitzr.

Basically, kibitzr would generate useful data from the web and notify the author
of updates via RSS feed and her normal, non-interruptus reading flow.

## How's it work?

Basically FeedR provides a drop-dead simple JSON-over-HTTP endpoint that any
client can POST to to store an RSS item. The JSON must have at least a title
field and a description field, and roughly corresponds to the following Go
struct:

```go
type FeedItem struct {
	Title       string
	Link        string
	AuthorName  string
	AuthorEmail string
	Description string
	Content     string
}
```

Field names are expected to be lowercase. As per the RSS 2.0 standard, either
`title` or `description` should be provided. 

To get to the RSS feed, simply issue a `GET /` to the root of the server.

The feed by default holds onto the last 100 items `POST`ed, and tosses anything
greater than that in FIFO order.

## How do I configure it?

Via flags, of course. The `-link` and `-description` flags are required, and
specify the source URL of the feed as well as a short description of it. Other
flags can be adjusted to change the metadata for the whole feed as well. Consult
`-help` for more information.

## How do I install it?

You don't.

Really.

This microservice is actually designed to be run from inside of a K8S or Docker
swarm environment. Instead, you either build the docker container using the
included Dockerfile and then write your own K8S or swarm YAML.

## Contributing

For now, I don't have time to write this for everybody's needs -- it was written
to fit my needs. Additionally, the code is pretty shoddy and non-idiomatic, as
it's the first service I've written using Go in anger.

Feel free to fork, send pull requests, and file issues. Can't guarantee I'll be
able to get to them in reasonable time, but please be patient. I usually like
contributions. :D

## Code of Conduct

Don't be a jerk, be excellent to each other.
