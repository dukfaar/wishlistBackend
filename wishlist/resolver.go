package wishlist

import graphql "github.com/graph-gophers/graphql-go"

type EntryResolver struct {
	Entry *WishlistEntry
}

type Resolver struct {
	Model *Model
}

func (r *Resolver) ID() *graphql.ID {
	id := graphql.ID(r.Model.ID.Hex())
	return &id
}

func (r *Resolver) Entries() *[]*EntryResolver {
	l := make([]*EntryResolver, len(r.Model.Entries))
	for i := range r.Model.Entries {
		l[i] = &EntryResolver{Entry: &r.Model.Entries[i]}
	}
	return &l
}

func (r *EntryResolver) ItemID() *graphql.ID {
	id := graphql.ID(r.Entry.ItemID.Hex())
	return &id
}
