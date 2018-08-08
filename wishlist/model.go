package wishlist

import (
	"github.com/dukfaar/goUtils/relay"
	"gopkg.in/mgo.v2/bson"
)

type WishlistEntry struct {
	ItemID bson.ObjectId `json:"itemId,omitempty" bson:"itemId,omitempty"`
}

type Model struct {
	ID      bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Entries []WishlistEntry `json:"entries,omitempty"`
}

var GraphQLType = `
type WishlistEntry {
	itemId: ID
}

type Wishlist {
	_id: ID
	entries: [WishlistEntry]
}
` +
	relay.GenerateConnectionTypes("Wishlist")
