package history

import (
	"reflect"
	"testing"
)

func TestHistoryRebuild(t *testing.T) {
	h := NewHistory()
	h.Add("foo")
	h.Rebuild("", false)
	expected := &History{
		histories: []string{"foo"},
		cache: map[string]int{
			"foo": 1,
		},
		tmp:      []string{"foo", ""},
		selected: 1,
		buf:      "",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}

	h.Add("fob")
	h.Rebuild("f", false)
	expected = &History{
		histories: []string{"foo", "fob"},
		cache: map[string]int{
			"foo": 1,
			"fob": 1,
		},
		tmp:      []string{"foo", "fob", ""},
		selected: 2,
		buf:      "f",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}

	h.Rebuild("foo", false)
	expected = &History{
		histories: []string{"foo", "fob"},
		cache: map[string]int{
			"foo": 1,
			"fob": 1,
		},
		tmp:      []string{"foo", ""},
		selected: 1,
		buf:      "foo",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}

	h.Add("fxb")
	expected = &History{
		histories: []string{"foo", "fob", "fxb"},
		cache: map[string]int{
			"foo": 1,
			"fob": 1,
			"fxb": 1,
		},
		tmp:      []string{"foo", "fob", "fxb", ""},
		selected: 3,
		buf:      "",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}
	h.Rebuild("", false)
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}

	h.Add("fob")
	expected = &History{
		histories: []string{"foo", "fob", "fxb"},
		cache: map[string]int{
			"foo": 1,
			"fob": 2,
			"fxb": 1,
		},
		tmp:      []string{"foo", "fob", "fxb", ""},
		selected: 3,
		buf:      "",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}
	// t.Logf("before - %#v\n", h)
	h.Remove("foo")

	expected = &History{
		histories: []string{"fob", "fxb"},
		cache: map[string]int{
			"fob": 2,
			"fxb": 1,
		},
		tmp:      []string{"fob", "fxb", ""},
		selected: 2,
		buf:      "",
	}
	if !reflect.DeepEqual(expected, h) {
		t.Errorf("Should be %#v, but got %#v", expected, h)
	}
	t.Logf("after  - %#v\n", h)
}

func TestHistoryOlder(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")

	// Prepare buffer
	buf := "echo 2"

	// [1 time] Call Older function
	buf1, changed := h.Older(buf)
	if !changed {
		t.Error("Should be changed history but not changed.")
	}
	if buf1 != "echo 1" {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf1)
	}

	// [2 times] Call Older function
	buf = "echo 1"
	buf2, changed := h.Older(buf)
	if changed {
		t.Error("Should be not changed history but changed.")
	}
	if !reflect.DeepEqual("echo 1", buf2) {
		t.Errorf("Should be %#v, but got %#v", "echo 1", buf2)
	}
}
