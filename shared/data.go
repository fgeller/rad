package shared

import (
	"fmt"
	"reflect"
)

type Pack struct {
	Name string
	Type string
}

type Member struct {
	Name      string
	Signature string
	Target    string
}

type Entry struct {
	Namespace []string
	Name      string
	Members   []Member
}

func (m Member) Eq(other Member) bool {
	return reflect.DeepEqual(m, other)
}

func (m Member) String() string {
	return fmt.Sprintf(
		"Member{Name: %v, Target: %v, Signature: %v}",
		m.Name,
		m.Target,
		m.Signature,
	)
}

func (e Entry) Eq(other Entry) bool {
	return reflect.DeepEqual(e, other)
}

func (e Entry) String() string {
	return fmt.Sprintf(
		"Entry{Name: %v, Namespace: %v, Members: %v}",
		e.Name,
		e.Namespace,
		e.Members,
	)
}

func IsSameEntry(a Entry, b Entry) bool {
	if a.Name != b.Name ||
		len(a.Namespace) != len(b.Namespace) {
		return false
	}
	for i := range a.Namespace {
		if a.Namespace[i] != b.Namespace[i] {
			return false
		}
	}

	return true
}

func MergeEntries(entries []Entry) []Entry {
	if len(entries) < 1 {
		return entries
	}

	unmerged := entries[1:]
	merged := []Entry{entries[0]}

merging:
	for ui := range unmerged {
		for mi := range merged {
			if IsSameEntry(unmerged[ui], merged[mi]) {
				merged[mi].Members = append(merged[mi].Members, unmerged[ui].Members...)
				continue merging
			}
		}
		merged = append(merged, unmerged[ui])
	}

	return merged
}
