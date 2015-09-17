package main

import "testing"

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	v := "scala.reflect.macros.contexts.Parsers@notify():Unit"
	expected := entry{
		Namespace: []string{"scala", "reflect", "macros", "contexts"},
		Name:      "Parsers",
		Members:   []member{{Name: "notify", Signature: "():Unit", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ := parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit"
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Name:      "SymDef",
		Members:   []member{{Name: "notifyAll", Signature: "():Unit", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]"
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Name:      "SymDef",
		Members: []member{{Name: "unapply",
			Signature: "(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]",
			Target:    "/" + target + v,
			Source:    source}},
		Source: source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.AnyRef@notify():Unit"
	expected = entry{
		Namespace: []string{"scala"},
		Name:      "AnyRef",
		Members:   []member{{Name: "notify", Signature: "():Unit", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.tools.cmd.Spec$@InfoextendsAnyRef"
	expected = entry{
		Namespace: []string{"scala", "tools", "cmd"},
		Name:      "Spec",
		Members:   []member{{Name: "Info", Signature: "extendsAnyRef", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.tools.ant.FastScalac"
	expected = entry{
		Namespace: []string{"scala", "tools", "ant"},
		Name:      "FastScalac",
		Members:   []member{{Name: "", Signature: "", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]"
	expected = entry{
		Namespace: []string{"scala", "collection", "MapLike"},
		Name:      "FilteredKeys",
		Members:   []member{{Name: "andThen", Signature: "[C](k:B=>C):PartialFunction[A,C]", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.util.Success@isFailure:Boolean"
	expected = entry{
		Namespace: []string{"scala", "util"},
		Name:      "Success",
		Members:   []member{{Name: "isFailure", Signature: ":Boolean", Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "package"
	expected = entry{
		Namespace: []string{},
		Name:      "package",
		Members:   []member{{Target: "/" + target + v, Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.concurrent.TrieMap$@Coll=CC[_,_]"
	expected = entry{
		Namespace: []string{"scala", "collection", "concurrent"},
		Name:      "TrieMap",
		Members:   []member{{Name: "Coll", Target: "/" + target + v, Signature: "=CC[_,_]", Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.MapLike@FilteredKeysextendsAbstractMap[A,B]withDefaultMap[A,B]"
	expected = entry{
		Namespace: []string{"scala", "collection"},
		Name:      "MapLike",
		Members:   []member{{Name: "FilteredKeys", Target: "/" + target + v, Signature: "extendsAbstractMap[A,B]withDefaultMap[A,B]", Source: source}},
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}
}
