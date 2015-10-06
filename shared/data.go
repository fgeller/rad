package shared

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Pack struct {
	Name    string
	Type    string
	Version string
	Created time.Time
}

type Member struct {
	Name   string
	Target string
}

type Namespace struct {
	Path    string
	Members []Member
}

func (m Member) Eq(other Member) bool {
	return reflect.DeepEqual(m, other)
}

func (m Member) String() string {
	return fmt.Sprintf("Member{Name: %v, Target: %v}", m.Name, m.Target)
}

func (n Namespace) Eq(other Namespace) bool {
	return reflect.DeepEqual(n, other)
}

func (n Namespace) String() string {
	return fmt.Sprintf("Namespace{Path: %v, Members: %v}", n.Path, n.Members)
}

func (n Namespace) Last() string {
	parts := strings.Split(n.Path, ".")
	return parts[len(parts)-1]
}

func Merge(ns []Namespace) []Namespace {
	if len(ns) < 1 {
		return ns
	}

	unmerged := ns[1:]
	merged := []Namespace{ns[0]}

merging:
	for ui := range unmerged {
		for mi := range merged {
			if unmerged[ui].Path == merged[mi].Path {
				// TODO: dedupe?
				merged[mi].Members = append(merged[mi].Members, unmerged[ui].Members...)
				continue merging
			}
		}
		merged = append(merged, unmerged[ui])
	}

	return merged
}
