package pad

import (
	"../shared"
	"flag"
	"fmt"
	"io"
	"net/http"
)

type indexer func() ([]shared.Entry, error)
type parser func(string, io.Reader) []shared.Entry
type downloader func(string) (*http.Response, error)

type config struct {
	indexer indexer
	name    string
	source  string
}

func readConfig() config {
	return config{}
}

func main() {

	var (
		indexerName = flag.String("indexer", "", "Indexer type for this pack (scala, java)")
		packName    = flag.String("name", "", "Name for this pack")
		source      = flag.String("source", "", "Source directory for this pack")
	)

	flag.Parse()

	fmt.Printf(
		"indexerName %v, packName %v, source %v\n",
		indexerName,
		packName,
		source,
	)

	// 0. copy all into temp location
	// 1. index
	// 2. serialize conf
	// 3. serialize entries
	// 4. zip it all up

}
