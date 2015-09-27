package main

import (
	"../shared"
	"testing"
)

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	v := "scala.reflect.macros.contexts.Parsers@notify():Unit"
	expected := shared.Entry{
		Namespace: []string{"scala", "reflect", "macros", "contexts"},
		Name:      "Parsers",
		Members:   []shared.Member{{Name: "notify", Signature: "():Unit", Target: "/" + target + v}},
	}

	actual, _ := parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit"
	expected = shared.Entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Name:      "SymDef",
		Members:   []shared.Member{{Name: "notifyAll", Signature: "():Unit", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]"
	expected = shared.Entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Name:      "SymDef",
		Members: []shared.Member{{Name: "unapply",
			Signature: "(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]",
			Target:    "/" + target + v,
		}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.AnyRef@notify():Unit"
	expected = shared.Entry{
		Namespace: []string{"scala"},
		Name:      "AnyRef",
		Members:   []shared.Member{{Name: "notify", Signature: "():Unit", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.tools.cmd.Spec$@InfoextendsAnyRef"
	expected = shared.Entry{
		Namespace: []string{"scala", "tools", "cmd"},
		Name:      "Spec",
		Members:   []shared.Member{{Name: "Info", Signature: "extendsAnyRef", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.tools.ant.FastScalac"
	expected = shared.Entry{
		Namespace: []string{"scala", "tools", "ant"},
		Name:      "FastScalac",
		Members:   []shared.Member{{Name: "", Signature: "", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]"
	expected = shared.Entry{
		Namespace: []string{"scala", "collection", "MapLike"},
		Name:      "FilteredKeys",
		Members:   []shared.Member{{Name: "andThen", Signature: "[C](k:B=>C):PartialFunction[A,C]", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.util.Success@isFailure:Boolean"
	expected = shared.Entry{
		Namespace: []string{"scala", "util"},
		Name:      "Success",
		Members:   []shared.Member{{Name: "isFailure", Signature: ":Boolean", Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "package"
	expected = shared.Entry{
		Namespace: []string{},
		Name:      "package",
		Members:   []shared.Member{{Target: "/" + target + v}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.concurrent.TrieMap$@Coll=CC[_,_]"
	expected = shared.Entry{
		Namespace: []string{"scala", "collection", "concurrent"},
		Name:      "TrieMap",
		Members:   []shared.Member{{Name: "Coll", Target: "/" + target + v, Signature: "=CC[_,_]"}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.MapLike@FilteredKeysextendsAbstractMap[A,B]withDefaultMap[A,B]"
	expected = shared.Entry{
		Namespace: []string{"scala", "collection"},
		Name:      "MapLike",
		Members:   []shared.Member{{Name: "FilteredKeys", Target: "/" + target + v, Signature: "extendsAbstractMap[A,B]withDefaultMap[A,B]"}},
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}
}
