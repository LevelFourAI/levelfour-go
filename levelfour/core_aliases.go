package levelfour

import "github.com/LevelFourAI/levelfour-go/core"

// Page re-exports [core.Page] so users never need to import the core package.
type Page[Cursor comparable, T any, R any] = core.Page[Cursor, T, R]

// PageIterator re-exports [core.PageIterator] so users never need to import the core package.
type PageIterator[Cursor comparable, T any, R any] = core.PageIterator[Cursor, T, R]

// ErrNoPages is a sentinel error returned by [Page.GetNextPage] when no pages remain.
var ErrNoPages = core.ErrNoPages
