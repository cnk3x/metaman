package douban

import (
	"testing"

	"github.com/cnk3x/metaman/pkg/strs"
)

func TestChart(t *testing.T) {
	s := MovieSearcher()
	items, err := s.Chart()
	if err != nil {
		t.Fatal(err)
	}
	for i, it := range items {
		t.Logf("%-2d %s", i+1, strs.Json(it))
	}
}

func TestSearch(t *testing.T) {
	s := MovieSearcher()
	items, err := s.Search(`惊天营救`, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i, it := range items {
		t.Logf("%-2d %s", i+1, strs.Json(it))
	}
}

func TestSubject(t *testing.T) {
	s := MovieSearcher()
	item, err := s.Subject(`35056376`)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", strs.Json(item))
}
