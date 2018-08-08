package main

import (
	"context"

	"github.com/dukfaar/goUtils/relay"
	"github.com/dukfaar/wishlistBackend/wishlist"
	graphql "github.com/graph-gophers/graphql-go"
	"gopkg.in/mgo.v2/bson"
)

type Resolver struct {
}

func (r *Resolver) Wishlists(ctx context.Context, args struct {
	First  *int32
	Last   *int32
	Before *string
	After  *string
}) (*wishlist.ConnectionResolver, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	var totalChannel = make(chan int)
	go func() {
		var total, _ = wishlistService.Count()
		totalChannel <- total
	}()

	var wishlistChannel = make(chan []wishlist.Model)
	go func() {
		result, _ := wishlistService.List(args.First, args.Last, args.Before, args.After)
		wishlistChannel <- result
	}()

	var (
		start string
		end   string
	)

	var wishlists = <-wishlistChannel

	if len(wishlists) == 0 {
		start, end = "", ""
	} else {
		start, end = wishlists[0].ID.Hex(), wishlists[len(wishlists)-1].ID.Hex()
	}

	hasPreviousPageChannel, hasNextPageChannel := relay.GetHasPreviousAndNextPage(len(wishlists), start, end, wishlistService)

	return &wishlist.ConnectionResolver{
		Models: wishlists,
		ConnectionResolver: relay.ConnectionResolver{
			relay.Connection{
				Total:           int32(<-totalChannel),
				From:            start,
				To:              end,
				HasNextPage:     <-hasNextPageChannel,
				HasPreviousPage: <-hasPreviousPageChannel,
			},
		},
	}, nil
}

func (r *Resolver) CreateWishlist(ctx context.Context, /*args struct {
}*/) (*wishlist.Resolver, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	inputModel := wishlist.Model{}

	newModel, err := wishlistService.Create(&inputModel)

	if err == nil {
		return &wishlist.Resolver{
			Model: newModel,
		}, nil
	}

	return nil, err
}

func (r *Resolver) AddItemToWishlist(ctx context.Context, args struct {
	WishlistId string
	ItemId     string
}) (*wishlist.Resolver, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	currentWishlist, err := wishlistService.FindByID(args.WishlistId)

	if err != nil {
		return nil, err
	}

	currentWishlist.Entries = append(currentWishlist.Entries, wishlist.WishlistEntry{ItemID: bson.ObjectIdHex(args.ItemId)})
	updatedWishlist, err := wishlistService.Update(args.WishlistId, &currentWishlist)

	if err == nil {
		return &wishlist.Resolver{
			Model: updatedWishlist,
		}, nil
	}

	return nil, err
}

func (r *Resolver) RemoveItemFromWishlist(ctx context.Context, args struct {
	WishlistId string
	ItemId     string
}) (*wishlist.Resolver, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	currentWishlist, err := wishlistService.FindByID(args.WishlistId)

	if err != nil {
		return nil, err
	}

	var foundPosition = -1
	for i := range currentWishlist.Entries {
		if currentWishlist.Entries[i].ItemID.Hex() == args.ItemId {
			foundPosition = i
			break
		}
	}

	currentWishlist.Entries[foundPosition] = currentWishlist.Entries[len(currentWishlist.Entries)-1]
	currentWishlist.Entries = currentWishlist.Entries[:len(currentWishlist.Entries)-1]
	updatedWishlist, err := wishlistService.Update(args.WishlistId, &currentWishlist)

	if err == nil {
		return &wishlist.Resolver{
			Model: updatedWishlist,
		}, nil
	}

	return nil, err
}

func (r *Resolver) DeleteWishlist(ctx context.Context, args struct {
	Id string
}) (*graphql.ID, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	deletedID, err := wishlistService.DeleteByID(args.Id)
	result := graphql.ID(deletedID)

	if err == nil {
		return &result, nil
	}

	return nil, err
}

func (r *Resolver) Wishlist(ctx context.Context, args struct {
	Id string
}) (*wishlist.Resolver, error) {
	wishlistService := ctx.Value("wishlistService").(wishlist.Service)

	queryWishlist, err := wishlistService.FindByID(args.Id)

	if err == nil {
		return &wishlist.Resolver{
			Model: queryWishlist,
		}, nil
	}

	return nil, err
}
