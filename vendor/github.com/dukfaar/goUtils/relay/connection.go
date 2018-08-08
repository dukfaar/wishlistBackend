package relay

import graphql "github.com/graph-gophers/graphql-go"

type Connection struct {
	Total           int32
	From            string
	To              string
	HasNextPage     bool
	HasPreviousPage bool
}

type ConnectionResolver struct {
	Connection Connection
}

func (r *ConnectionResolver) TotalCount() *int32 {
	return &r.Connection.Total
}

func (r *ConnectionResolver) PageInfo() *PageInfoResolver {
	p := PageInfo{
		StartPage:       graphql.ID(r.Connection.From),
		EndPage:         graphql.ID(r.Connection.To),
		HasNextPage:     r.Connection.HasNextPage,
		HasPreviousPage: r.Connection.HasPreviousPage,
	}
	return &PageInfoResolver{PageInfo: &p}
}

func GenerateConnectionTypes(baseType string) string {
	return `
	type ` + baseType + `Connection {
		edges: [` + baseType + `Edge]
		totalCount: Int
		pageInfo: PageInfo!
	}

	type ` + baseType + `Edge {
		node: ` + baseType + `
		cursor: ID!
	}
	`
}
