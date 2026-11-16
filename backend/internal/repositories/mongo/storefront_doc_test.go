package mongorepo

import (
	"reflect"
	"testing"

	"shop24/internal/domain/storefront"
)

func TestStorefrontRoundTrip(t *testing.T) {
	cfg := storefront.Default("North", "Wool and linen")
	if got := storefrontFromDoc(storefrontToDoc(cfg)); !reflect.DeepEqual(cfg, got) {
		t.Fatalf("round trip mismatch\n got %#v\nwant %#v", got, cfg)
	}
}
