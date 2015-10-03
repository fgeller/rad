package main

import (
	"../shared"
	"testing"
)

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	v := "scala.reflect.macros.contexts.Parsers@notify():Unit"
	expected := shared.Namespace{
		Path:    []string{"scala", "reflect", "macros", "contexts", "Parsers"},
		Members: []shared.Member{{Name: "notify", Target: "/" + target + v}},
	}

	actual, _ := parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit"
	expected = shared.Namespace{
		Path:    []string{"scala", "reflect", "reify", "utils", "Extractors", "SymDef"},
		Members: []shared.Member{{Name: "notifyAll", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]"
	expected = shared.Namespace{
		Path:    []string{"scala", "reflect", "reify", "utils", "Extractors", "SymDef"},
		Members: []shared.Member{{Name: "unapply", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.AnyRef@notify():Unit"
	expected = shared.Namespace{
		Path:    []string{"scala", "AnyRef"},
		Members: []shared.Member{{Name: "notify", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.tools.cmd.Spec$@InfoextendsAnyRef"
	expected = shared.Namespace{
		Path:    []string{"scala", "tools", "cmd", "Spec"},
		Members: []shared.Member{{Name: "Info", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.tools.ant.FastScalac"
	expected = shared.Namespace{
		Path:    []string{"scala", "tools", "ant", "FastScalac"},
		Members: []shared.Member{{Name: "", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]"
	expected = shared.Namespace{
		Path:    []string{"scala", "collection", "MapLike", "FilteredKeys"},
		Members: []shared.Member{{Name: "andThen", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.util.Success@isFailure:Boolean"
	expected = shared.Namespace{
		Path:    []string{"scala", "util", "Success"},
		Members: []shared.Member{{Name: "isFailure", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "package"
	expected = shared.Namespace{
		Path:    []string{"package"},
		Members: []shared.Member{{Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.concurrent.TrieMap$@Coll=CC[_,_]"
	expected = shared.Namespace{
		Path:    []string{"scala", "collection", "concurrent", "TrieMap"},
		Members: []shared.Member{{Name: "Coll", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.MapLike@FilteredKeysextendsAbstractMap[A,B]withDefaultMap[A,B]"
	expected = shared.Namespace{
		Path:    []string{"scala", "collection", "MapLike"},
		Members: []shared.Member{{Name: "FilteredKeys", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}
}
