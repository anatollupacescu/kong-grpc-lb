package main

import (
	context "context"
	"fmt"
	"net/url"

	"github.com/golang/protobuf/ptypes"

	product "github.com/anatollupacescu/atlant/internal"
	"github.com/anatollupacescu/atlant/proto"
)

type ProductServiceServer struct {
	proto.UnimplementedProductServiceServer
	product.App
}

func (p *ProductServiceServer) List(ctx context.Context, in *proto.ListReq) (*proto.ListRes, error) {
	page := toPageWithDefaults(in)

	products, err := p.App.ListProductPrices(ctx, page)
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

func toPageWithDefaults(in *proto.ListReq) product.Page {
	page := product.Page{
		Field: "Name",
		Limit: 10,
		Order: 1,
	}

	if in.SortBy != "" {
		page.Field = in.SortBy
	}

	if in.Limit > 0 {
		page.Limit = int64(in.Limit)
	}

	if in.SortDesc {
		page.Order = -1
	}

	page.Page = int64(1 + in.Page)

	return page
}
