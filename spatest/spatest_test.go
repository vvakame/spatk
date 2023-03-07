package spatest_test

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/vvakame/spatk/spatest"
	"google.golang.org/api/iterator"
)

func TestSetupDedicatedDB(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	dbName, closeFn, err := spatest.SetupDedicatedDB(
		ctx,
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		err := closeFn()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Logf("spanner db: %s", dbName)

	spCli, err := spanner.NewClient(ctx, dbName)
	if err != nil {
		t.Fatal(err)
	}

	iter := spCli.Single().Query(ctx, spanner.NewStatement("SELECT 1"))
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			t.Fatal(err)
		}

		var i int64
		err = row.Column(0, &i)
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("value: %d", i)
	}
}
