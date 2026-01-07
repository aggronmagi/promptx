package history

import (
	"os"
	"reflect"
	"testing"
)

func TestHistoryRebuild(t *testing.T) {
	h := NewHistory(WithMaxSize(0), WithDeduplicate(true))
	h.Add("foo")
	h.Rebuild("", false)
	
	if !reflect.DeepEqual([]string{"foo"}, h.histories) {
		t.Errorf("Should be %v, but got %v", []string{"foo"}, h.histories)
	}

	h.Add("fob")
	h.Rebuild("f", false)
	if !reflect.DeepEqual([]string{"foo", "fob"}, h.histories) {
		t.Errorf("Should be %v, but got %v", []string{"foo", "fob"}, h.histories)
	}

	h.Rebuild("foo", false)
	if !reflect.DeepEqual([]string{"foo", ""}, h.tmp) {
		t.Errorf("Should be %v, but got %v", []string{"foo", ""}, h.tmp)
	}

	h.Add("fxb")
	if !reflect.DeepEqual([]string{"foo", "fob", "fxb"}, h.histories) {
		t.Errorf("Should be %v, but got %v", []string{"foo", "fob", "fxb"}, h.histories)
	}

	h.Add("fob")
	// With deduplicate, fob should move to the end
	expected := []string{"foo", "fxb", "fob"}
	if !reflect.DeepEqual(expected, h.histories) {
		t.Errorf("Should be %v, but got %v", expected, h.histories)
	}

	h.Remove("foo")
	if !reflect.DeepEqual([]string{"fxb", "fob"}, h.histories) {
		t.Errorf("Should be %v, but got %v", []string{"fxb", "fob"}, h.histories)
	}
}

func TestHistoryOlder(t *testing.T) {
	h := NewHistory()
	h.Add("echo 1")
	h.Add("echo 2")

	// Current state: histories=["echo 1", "echo 2"], tmp=["echo 1", "echo 2", ""], selected=2
	
	buf := "echo 3"
	
	// [1 time] Call Older function -> selected becomes 1, returns "echo 2"
	buf1, changed := h.Older(buf)
	if !changed {
		t.Error("Should be changed history but not changed.")
	}
	if buf1 != "echo 2" {
		t.Errorf("Should be %s, but got %s", "echo 2", buf1)
	}

	// [2 times] Call Older function -> selected becomes 0, returns "echo 1"
	buf2, changed := h.Older(buf1)
	if !changed {
		t.Error("Should be changed history but not changed.")
	}
	if buf2 != "echo 1" {
		t.Errorf("Should be %s, but got %s", "echo 1", buf2)
	}
	
	// [3 times] Call Older function -> already at 0, returns false
	_, changed = h.Older(buf2)
	if changed {
		t.Error("Should not be changed.")
	}
}

func TestHistoryDeduplicate(t *testing.T) {
	h := NewHistory(WithDeduplicate(true))
	h.Add("ls")
	h.Add("cd")
	h.Add("ls")
	
	expected := []string{"cd", "ls"}
	if !reflect.DeepEqual(expected, h.histories) {
		t.Errorf("Should be %v, but got %v", expected, h.histories)
	}
}

func TestHistoryMaxSize(t *testing.T) {
	h := NewHistory(WithMaxSize(2))
	h.Add("1")
	h.Add("2")
	h.Add("3")
	
	expected := []string{"2", "3"}
	if !reflect.DeepEqual(expected, h.histories) {
		t.Errorf("Should be %v, but got %v", expected, h.histories)
	}
}

func TestHistoryTimestamp(t *testing.T) {
	h := NewHistory(WithTimestamp(true))
	h.Add("ls")
	
	if len(h.histories) != 1 {
		t.Fatal("should have 1 history item")
	}
	if !reflect.DeepEqual("ls", h.extractCommand(h.histories[0])) {
		t.Errorf("Extracted command should be 'ls', but got %s", h.extractCommand(h.histories[0]))
	}
	
	h.Rebuild("", false)
	if h.tmp[0] != "ls" {
		t.Errorf("Suggestion should be 'ls', but got %s", h.tmp[0])
	}
}

func TestHistoryLoadTimestamp(t *testing.T) {
	h := NewHistory()
	content := ": 1641542400:0;ls\n: 1641542460:0;cd\n"
	
	tmpfile, err := os.CreateTemp("", "history_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	
	if err := h.Load(tmpfile.Name()); err != nil {
		t.Fatal(err)
	}
	
	if !h.timestamp {
		t.Error("timestamp mode should be auto-detected")
	}
	
	expected := []string{": 1641542400:0;ls", ": 1641542460:0;cd"}
	if !reflect.DeepEqual(expected, h.histories) {
		t.Errorf("Should be %v, but got %v", expected, h.histories)
	}
}
