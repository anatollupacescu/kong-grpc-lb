package main

import (
	context "context"
	"fmt"
	"net/url"

	product "github.com/anatollupacescu/atlant/internal"
	"github.com/anatollupacescu/atlant/proto"
	"github.com/golang/protobuf/ptypes"
)

type ProductServiceServer struct {
	proto.UnimplementedProductServiceServer
	product.App
}

func (p *ProductServiceServer) List(ctx context.Context, in *proto.ListReq) (*proto.ListRes, error) {
	sort := toSortWithDefaults(in)

	products, err := p.App.ListProductPrices(ctx, sort)
	if err != nil {
		return nil, err
	}

	var protoProducts []*proto.Product
	for _, product := range products {
		lastUpdated, err := ptypes.TimestampProto(product.LastUpdated)
		if err != nil {
			return nil, fmt.Errorf("product %v has bad timestamp: %w", product.ID, err)
		}

		protoProducts = append(protoProducts, &proto.Product{
			Id:            product.ID.Hex(),
			Name:          product.Name,
			Price:         product.Price,
			LastUpdatedOn: lastUpdated,
			UpdateCount:   uint64(product.UpdateCount),
		})
	}

	response := proto.ListRes{
		Products: protoProducts,
	}

	return &response, nil
}

func (p *ProductServiceServer) Fetch(ctx context.Context, in *proto.FetchReq) (*proto.FetchRes, error) {
	url, err := url.Parse(in.Url)
	if err != nil {
		return nil, err
	}

	if err := p.App.StoreProductPrices(ctx, url); err != nil {
		return nil, err
	}

	response := proto.FetchRes{
		Success: true,
	}

	return &response, nil
}

func toSortWithDefaults(in *proto.ListReq) product.Sort {
	sort := product.Sort{
		Field: "Name",
		Limit: 10,
		Order: 1,
		Page:  1,
	}

	if in.SortBy != "" {
		sort.Field = in.SortBy
	}

	if in.Limit > 0 {
		sort.Limit = int64(in.Limit)
	}

	if in.SortDesc {
		sort.Order = -1
	}

	if sort.Page > 0 { //one-based array indexing
		sort.Page = int64(in.Page)
	}

	return sort
}
