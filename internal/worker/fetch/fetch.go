// Package fetch gets automatically new demo files
package fetch

type Fetcher interface {
	Fetch() ([]byte, error)
}

// TODO: Implement for steam
// TODO: Implement for faceit
