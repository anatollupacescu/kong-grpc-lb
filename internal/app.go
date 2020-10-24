package product

import (
	"context"
	"net/url"
	"time"

	pagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

type App struct {
	ProductDB *mongo.Collection
}

type productDTO struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	LastUpdated time.Time     `bson:"LastUpdated,omitempty"`
	UpdateCount int           `bson:"UpdateCount,omitempty"`
	Name        string        `bson:"Name,omitempty"`
	Price       string        `bson:"Price,omitempty"`
}

type CSVRow struct {
	Name  string
	Price string
}

type CSVRows []CSVRow

type Page struct {
	Field       string
	Order       int
	Limit, Page int64
}

func (a App) ListProductPrices(ctx context.Context, page Page) ([]productDTO, error) {
	var all bson.M
	paginatedData, err := pagination.New(a.ProductDB).
		Limit(page.Limit).
		Page(page.Page).
		Sort(page.Field, page.Order).
		Filter(all).
		Find()

	if err != nil {
		return nil, err
	}

	var products []productDTO
	for _, raw := range paginatedData.Data {
		var mapped productDTO
		if marshallErr := bson.Unmarshal(raw, &mapped); marshallErr == nil {
			products = append(products, mapped)
		}
	}

	return products, nil
}

func (a App) StoreProductPrices(ctx context.Context, url *url.URL) error {
	now := time.Now()

	reader, err := fetchPriceFile(url)
	if err != nil {
		return err
	}

	data, err := parseCSV(reader)
	if err != nil {
		return err
	}

	products, err := mapToProducts(data)
	if err != nil {
		return err
	}

	records := transform(products, now)

	writeModels := toWriteModels(records)
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(true) //default value is true anyway but it's good to have it explicit

	if _, err := a.ProductDB.BulkWrite(ctx, writeModels, &bulkOption); err != nil {
		errStatus := status.Errorf(codes.Internal, "insert records: %v", err)

		return errStatus
	}

	return nil
}

func toWriteModels(rr []productDTO) (models []mongo.WriteModel) {
	for _, r := range rr {
		m := mongo.NewUpdateOneModel()
		m.SetFilter(bson.M{"Name": r.Name})
		m.SetUpdate(bson.M{"$set": bson.M{
			"Price":       r.Price,
			"LastUpdated": r.LastUpdated,
			"UpdateCount": r.UpdateCount}})
		m.SetUpsert(true)

		models = append(models, m)
	}

	return
}
