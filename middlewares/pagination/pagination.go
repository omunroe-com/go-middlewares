// Package pagination provides support pagination for
// jsonapi that is given in the query parameter of the
// request url.
// It supports the following query pattern
//		page[number]={num}&page[size]={size}
package pagination

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

// String output the details of context key
func (c contextKey) String() string {
	return "pagination context key " + string(c)
}

// ContextKeyPagination is the key used for stroing Prop struct in request
// context
var (
	ContextKeyPagination = contextKey("page")
)

// Prop represents various pagination properties
type Prop struct {
	// Total no of records that will be paginated
	Records int
	// No of entries to have per page
	Entries int
	// Current page no
	Current int
}

// MiddlewareFn parses the pagination query string and stores
// it in request context under  ContextKeyPagination variable as a Prop type
func MiddlewareFn(fn http.HandlerFunc) http.HandlerFunc {
	newFn := func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if num, ok := values["page[number]"]; ok {
			if size, ok := values["page[size]"]; ok {
				curr, _ := strconv.Atoi(num[0])
				entries, _ := strconv.Atoi(size[0])
				prop := &Prop{
					Current: curr,
					Entries: entries,
				}
				ctx := context.WithValue(r.Context(), ContextKeyPagination, prop)
				fn(w, r.WithContext(ctx))
			}
		} else {
			fn(w, r)
		}
	}
	return newFn
}