package stock

import (
	"api/domain/model"
	"context"
	"errors"
	"reflect"
	"testing"
)

type fakeRepository struct {
	stocks []model.Stock
	err    error
}

func (f *fakeRepository) Upsert(_ context.Context, stock model.Stock) error {

	f.stocks = append(f.stocks, stock)

	return f.err
}

func TestImport_PropagatesRepositoryError(t *testing.T) {
	// Arrange
	wantErr := errors.New("save failed")
	repository := &fakeRepository{err: wantErr}
	importer := NewStockImporter(repository)
	stocks := []model.Stock{{Symbol: "72030"}}

	// Act
	err := importer.Import(context.Background(), stocks)

	// Assert
	if !errors.Is(err, wantErr) {
		t.Errorf("Import() error = %v, want wrapped %v", err, wantErr)
	}
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
