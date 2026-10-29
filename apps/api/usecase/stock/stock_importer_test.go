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

func TestImport_SavesStocksInSymbolOrder(t *testing.T) {
	// Arrange
	repository := &fakeRepository{}
	importer := NewStockImporter(repository)

	stocks := []model.Stock{
		{Symbol: "86970"},
		{Symbol: "72030"},
	}
	want := []string{"72030", "86970"}

	// Act
	err := importer.Import(context.Background(), stocks)

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if len(repository.stocks) != len(want) {
		t.Fatalf("saved stocks = %d, want %d", len(repository.stocks), len(want))
	}

	got := []string{
		repository.stocks[0].Symbol,
		repository.stocks[1].Symbol,
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("symbols = %v, want %v", got, want)
	}
}
