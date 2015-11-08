# rad - read api docs

Open source app to read API docs offline.

## Quick start:

Download a release from the [release section](https://github.com/fgeller/rad/releases), unpack it and start it. Or as a one-liner (for linux replace `darwin` with `linux`):

    curl http://geller.io/rad/sad-1.0.0-darwin-amd64.xz | tar Jx > sad && ./sad

Then move the binary to a location that's on your `PATH` and make it executable. For example, to add it to `~/bin`:

    mv sad ~/bin/ && chmod +x ~/bin/sad

Just make sure you download packs before you go offline :)

All releases including snapshots are available [here](http://geller.io/rad/).

Alternatively, if you just want to give it a try there is a sample instance
running [here](http://geller.io:3024). Keep in mind though that the latency to
that instance probably slows down the app's responsiveness. For example in my
case, I have a latency of about 500ms to the sample instance while most requests
render successfully in 40-100ms for an instance running on my tiny laptop.

Contents:

* [Pictures!](https://github.com/fgeller/rad#pictures)
* [Available packs](https://github.com/fgeller/rad#available-packs)
* [Integrations](https://github.com/fgeller/rad#integrations)
  * [Browser](https://github.com/fgeller/rad#browser)
  * [Alfred](https://github.com/fgeller/rad#alfred)
  * [HUBOT + Slack](https://github.com/fgeller/rad#hubot--slack)
* [Why?](https://github.com/fgeller/rad#why)
* [How?](https://github.com/fgeller/rad#how)
* [More packs and contributions](https://github.com/fgeller/rad#more-packs-and-contributions)
* [Others](https://github.com/fgeller/rad#others)

## Pictures!

Start screen:

![Blank](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/blank.png)

Searching Go's stdlib:

![Go](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/go.png)

Searching Scala's stdlib using regular expressions:

![Scala](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/scala.png)

Displaying installed and available packs:

![Packs](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/packs.png)

## Available packs

* [Clojure](https://clojuredocs.org)
* [Django](https://clojuredocs.org)
* [Go](https://golang.org)
  * [Standard packages](https://godoc.org/-/go)
  * [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
  * [github.com/pborman/uuid](https://github.com/pborman/uuid)
  * [golang.org/x/net/html](https://godoc.org/golang.org/x/net/html)
* [Java](http://www.oracle.com/technetwork/java/javase/overview/index.html)
* [jQuery](http://jquery.com)
* [lodash](http://lodash.com)
* [Man pages](http://man7.org/linux/man-pages/index.html)
* [MDN Javascript](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference)
* [node.js](http://nodejs.org)
* [Python](https://www.python.org)
* [react.js](https://facebook.github.io/react/)
* [Scala](http://www.scala-lang.org)
  * Standard library (split into compiler, library and reflect)
  * [akka](http://www.scala-lang.org)
  * [Gremlin-Scala](https://github.com/mpollmeier/gremlin-scala)
  * [Kafka](http://kafka.apache.org)
  * [Shapeless](https://github.com/milessabin/shapeless)
  * [Spray](http://spray.io)

## Integrations

While `sad` includes a UI, it might be faster or more convenient to query
through one of the following integrations:

### Browser

You can save some key strokes by assigning a custom search engine keyword to
`rad`. For example in Chrome:

![Chrome Search Engine](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/chrome-search-engine.png)

This will allow you to search for documentation by using a custom keyword
directly from the location bar. To speed things up even more, you can enable
live loading of documentation in `rad` settings:

![Live updates](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/live-updates.png)

This will allow you to jump directly to the first result after you've entered a
query in the location bar (`shift + →` and `shift + ←` will cycle through the
search results if the query wasn't unique).

### Alfred

You can query a running `rad` instance from [Alfred](https://www.alfredapp.com/)
if you install the included [workflow](https://raw.githubusercontent.com/fgeller/rad/master/integrations/alfred/rad.alfredworkflow).
Then you can start a query with the configured keyword (the default is simply `r`):

![Alfred](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/alfred.png)

Selecting a result will open the documention in your default browser via `rad`.

### HUBOT + Slack

[rad.coffee](https://github.com/fgeller/rad/blob/master/integrations/rad.coffee)
is a small [HUBOT](https://hubot.github.com) script that works with the Slack
adapter:

![HUBOT + Slack](https://raw.githubusercontent.com/fgeller/rad/master/screenshots/hubot-slack.png)

You start a query with the keyword `rad` and the results are links to a running
`sad` instance.

## Why?

I enjoy open source software and hacking on my projects while on the go or in
cafes. And this project gave me an excuse to learn more about
[Go](https://golang.org) and [react.js](https://facebook.github.io/react/) :)

![https://golang.org/doc/faq#Whats_the_origin_of_the_mascot](https://golang.org/doc/gopher/pencil/gopherswrench.jpg)

## How?

rad is split into multiple components:

* `sad`: web app to serve API docs including a UI.
* `pad`: app to package API docs.
* `sap`: web app to serve documentation packs.

All three apps are written in Go and the `sad` front-end uses react.js and plain
old HTML/Javascript. `sad` uses web-sockets and Go's channels for streaming
results immediately to the UI.

### Pack format

Packs are zip archives with the following structure:

    pack-name
     |
     |-- pack.json
     |-- data.json
     `-- html-contents
       |
       |-- first.html
       `-- second.html

They contain a top-level directory that has the same name as the pack. The pack
name is used to identify the pack for searching, in the above example the name
is `pack-name`. The top-level directory contains two data files that include
meta-data about the pack `pack.json` and the documentation entries `data.json`.

A `pack.json` file has at least the following information: Name, Type, Version, Created, Description. Consider the following example:

    {
      "Name": "pack-name",
      "Type": "go",
      "Version": "3.1.4",
      "Created": "2015-10-25T15:42:01.936261923+13:00",
      "Description": "This is an <i>HTML</i> string containing source and copyright information."
    }

The `data.json` file contains the entries that `sad` will index and search. An
entry is grouped under a path and contains member entries that are defined by a
name and a relative link to their location within the pack archive. Consider the
following example:

    [
      {
        "Path": "com.example.super",
        "Members": [
          {
            "Name": "Driver",
            "Target": "html-contents/first.html#Driver"
          }
        ]
      },
      {
        "Path": "com.example.super",
        "Members": [
          {
            "Name": "Runner",
            "Target": "html-contents/second.html#Runner"
          }
        ]
      }
    ]


## More packs and contributions

![Moar!](http://memeorama.com/wp-content/uploads/2012/01/MOAR.jpg)

The currently included packs are driven by my use and what friends are
interested in -- if you'd like to see another pack included please create an
[issue](https://github.com/fgeller/rad/issues/new) or even send me a wrapped
pack :)

## Others

* [Dash](https://kapeli.com) - Very mature OSX app with lots of packs available.
* [DevDocs](http://devdocs.io) - Web app that supports offline mode with lots of packs especially for front-end tech.

<br /><br />
<p align="center"><sup>Made with ❤ in Piha.</sup></p>
