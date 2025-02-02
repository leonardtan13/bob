package orm

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stephenafamo/bob"
)

func TestHooks(t *testing.T) {
	var H Hooks[*string]

	// Test Adding Hooks
	for i := 0; i < 5; i++ {
		initial := len(H.hooks)
		f := func(ctx context.Context, _ bob.Executor, s *string) (context.Context, error) {
			*s = *s + fmt.Sprintf("%d", initial+1)
			return ctx, nil
		}
		H.Add(f)
		if len(H.hooks) != initial+1 {
			t.Fatalf("Did not add hook number %d", i+1)
		}
	}

	s := ""
	if _, err := H.Do(context.Background(), nil, &s); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff("12345", s); diff != "" {
		t.Fatal(diff)
	}
}
