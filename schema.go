package main

import (
	"github.com/dukfaar/goUtils/relay"

	"github.com/dukfaar/wishlistBackend/wishlist"
)

var Schema string = `
		schema {
			query: Query
			mutation: Mutation
		}

		type Query {
			wishlists(first: Int, last: Int, before: String, after: String): WishlistConnection!
			wishlist(id: ID!): Wishlist!
		}

		type Mutation {
			createWishlist(): Wishlist!
			deleteWishlist(id: ID!): ID

			addItemToWishlist(wishlistId: ID!, itemId: ID!): Wishlist
			removeItemFromWishlist(wishlistId: ID!, itemId: ID!): Wishlist
		}` +
	relay.PageInfoGraphQLString +
	wishlist.GraphQLType
