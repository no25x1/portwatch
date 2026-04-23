package tagger_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/user/portwatch/internal/tagger"
)

func TestRegistry_ConcurrentSetGet(t *testing.T) {
	r := tagger.New()
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers * 2)

	for i := 0; i < workers; i++ {
		target := fmt.Sprintf("host%d:80", i)
		go func(tgt string, idx int) {
			defer wg.Done()
			r.Set(tgt, tagger.Tags{"worker": fmt.Sprintf("%d", idx)})
		}(target, i)

		go func(tgt string) {
			defer wg.Done()
			r.Get(tgt) // may or may not find it; must not race
		}(target)
	}
	wg.Wait()

	all := r.All()
	if len(all) > workers {
		t.Fatalf("unexpected extra entries: %d", len(all))
	}
}

func TestRegistry_OverwritePreservesLatest(t *testing.T) {
	r := tagger.New()
	r.Set("svc:8080", tagger.Tags{"env": "staging"})
	r.Set("svc:8080", tagger.Tags{"env": "prod", "region": "eu"})

	tags, ok := r.Get("svc:8080")
	if !ok {
		t.Fatal("expected target to exist after overwrite")
	}
	if tags["env"] != "prod" {
		t.Fatalf("expected env=prod after overwrite, got %q", tags["env"])
	}
	if tags["region"] != "eu" {
		t.Fatalf("expected region=eu, got %q", tags["region"])
	}
}

func TestRegistry_AllReturnsSnapshot(t *testing.T) {
	r := tagger.New()
	r.Set("a:1", tagger.Tags{"k": "v"})
	snap := r.All()
	r.Set("b:2", tagger.Tags{"k": "v2"})

	if len(snap) != 1 {
		t.Fatalf("snapshot should not reflect post-snapshot mutations, len=%d", len(snap))
	}
}
