package relay

import (
	"github.com/dukfaar/goUtils/service"
	"github.com/graph-gophers/graphql-go"
)

type PageInfo struct {
	StartPage       graphql.ID
	EndPage         graphql.ID
	HasNextPage     bool
	HasPreviousPage bool
}

type PageInfoResolver struct {
	PageInfo *PageInfo
}

func (r *PageInfoResolver) StartCursor() *graphql.ID {
	return &r.PageInfo.StartPage
}

func (r *PageInfoResolver) EndCursor() *graphql.ID {
	return &r.PageInfo.EndPage
}

func (r *PageInfoResolver) HasNextPage() *bool {
	return &r.PageInfo.HasNextPage
}

func (r *PageInfoResolver) HasPreviousPage() *bool {
	return &r.PageInfo.HasPreviousPage
}

const PageInfoGraphQLString string = `
type PageInfo {
	startCursor: ID
	endCursor: ID
	hasNextPage: Boolean
	hasPreviousPage: Boolean
}
`

func GetHasPreviousAndNextPage(resultLength int, firstItemID string, lastItemID string, service service.DBService) (chan bool, chan bool) {
	hasPreviousPage := make(chan bool)
	hasNextPage := make(chan bool)

	go func() {
		if resultLength > 0 {
			result, _ := service.HasElementBeforeID(firstItemID)
			hasPreviousPage <- result
		} else {
			hasPreviousPage <- false
		}
	}()

	go func() {
		if resultLength > 0 {
			result, _ := service.HasElementAfterID(lastItemID)
			hasNextPage <- result
		} else {
			hasNextPage <- false
		}
	}()

	return hasPreviousPage, hasNextPage
}
