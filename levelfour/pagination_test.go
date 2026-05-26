package levelfour

import (
	"context"
	"testing"

	"github.com/LevelFourAI/levelfour-go/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectAll_SinglePage(t *testing.T) {
	page := &core.Page[int, string, any]{
		Results: []string{"a", "b", "c"},
		NextPageFunc: func(ctx context.Context) (*core.Page[int, string, any], error) {
			return nil, core.ErrNoPages
		},
	}

	results, err := CollectAll(context.Background(), page)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, results)
}

func TestCollectAll_MultiplePages(t *testing.T) {
	page2 := &core.Page[int, int, any]{
		Results: []int{3, 4},
		NextPageFunc: func(ctx context.Context) (*core.Page[int, int, any], error) {
			return nil, core.ErrNoPages
		},
	}
	page1 := &core.Page[int, int, any]{
		Results: []int{1, 2},
		NextPageFunc: func(ctx context.Context) (*core.Page[int, int, any], error) {
			return page2, nil
		},
	}

	results, err := CollectAll(context.Background(), page1)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, results)
}

func TestCollectAll_EmptyPage(t *testing.T) {
	page := &core.Page[int, string, any]{
		Results: []string{},
		NextPageFunc: func(ctx context.Context) (*core.Page[int, string, any], error) {
			return nil, core.ErrNoPages
		},
	}

	results, err := CollectAll(context.Background(), page)
	require.NoError(t, err)
	assert.Nil(t, results)
}
