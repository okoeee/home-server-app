package stock

import (
	"api/domain/model"
	"context"
	"reflect"
	"testing"
)

type fakeRepository struct {
	stocks []model.Stock
}

func (f *fakeRepository) Upsert(_ context.Context, stock model.Stock) error {

	f.stocks = append(f.stocks, stock)

	return nil
}

func TestImporterImportsStocksInSymbolOrder(t *testing.T) {
	repository := &fakeRepository{}
	importer := NewStockImporter(repository)

	stocks := []model.Stock{
		{Symbol: "86970"},
		{Symbol: "72030"},
	}

	if err := importer.Import(context.Background(), stocks); err != nil {
		t.Fatal(err)
	}

	got := []string{
		repository.stocks[0].Symbol,
		repository.stocks[1].Symbol,
	}
	want := []string{"72030", "86970"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("symbols = %v, want %v", got, want)
	}
}
