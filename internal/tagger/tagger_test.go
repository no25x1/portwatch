package tagger

import (
	"testing"
)

func TestTags_CloneIsIndependent(t *testing.T) {
	orig := Tags{"env": "prod", "region": "us-east"}
	clone := orig.Clone()
	clone["env"] = "staging"
	if orig["env"] != "prod" {
		t.Fatalf("expected original to be unmodified, got %q", orig["env"])
	}
}

func TestTags_Has(t *testing.T) {
	tags := Tags{"env": "prod"}
	if !tags.Has("env") {
		t.Fatal("expected Has to return true for existing key")
	}
	if tags.Has("missing") {
		t.Fatal("expected Has to return false for missing key")
	}
}

func TestTags_Get(t *testing.T) {
	tags := Tags{"role": "web"}
	v, ok := tags.Get("role")
	if !ok || v != "web" {
		t.Fatalf("expected (web, true), got (%q, %v)", v, ok)
	}
	_, ok = tags.Get("nope")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestRegistry_SetAndGet(t *testing.T) {
	r := New()
	r.Set("host1:80", Tags{"env": "prod"})
	tags, ok := r.Get("host1:80")
	if !ok {
		t.Fatal("expected target to be found")
	}
	if tags["env"] != "prod" {
		t.Fatalf("expected env=prod, got %q", tags["env"])
	}
}

func TestRegistry_Get_Missing(t *testing.T) {
	r := New()
	_, ok := r.Get("ghost:9999")
	if ok {
		t.Fatal("expected missing target to return false")
	}
}

func TestRegistry_GetIsACopy(t *testing.T) {
	r := New()
	r.Set("h:80", Tags{"k": "v"})
	tags, _ := r.Get("h:80")
	tags["k"] = "mutated"
	again, _ := r.Get("h:80")
	if again["k"] != "v" {
		t.Fatal("registry value was mutated through returned map")
	}
}

func TestRegistry_Delete(t *testing.T) {
	r := New()
	r.Set("h:443", Tags{"tls": "true"})
	r.Delete("h:443")
	_, ok := r.Get("h:443")
	if ok {
		t.Fatal("expected deleted target to be absent")
	}
}

func TestRegistry_All(t *testing.T) {
	r := New()
	r.Set("a:80", Tags{"x": "1"})
	r.Set("b:80", Tags{"x": "2"})
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}
