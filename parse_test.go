package main

import "testing"

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	v := "scala.reflect.macros.contexts.Parsers@notify():Unit"
	expected := entry{
		Namespace: []string{"scala", "reflect", "macros", "contexts"},
		Entity:    "Parsers",
		Member:    "notify",
		Signature: "():Unit",
		Target:    target + v,
		Source:    source,
	}

	actual, _ := parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got %v expected %v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit"
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Entity:    "SymDef",
		Member:    "notifyAll",
		Signature: "():Unit",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]"
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils", "Extractors"},
		Entity:    "SymDef",
		Member:    "unapply",
		Signature: "(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	v = "scala.AnyRef@notify():Unit"
	expected = entry{
		Namespace: []string{"scala"},
		Entity:    "AnyRef",
		Member:    "notify",
		Signature: "():Unit",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	v = "scala.tools.cmd.Spec$@InfoextendsAnyRef"
	expected = entry{
		Namespace: []string{"scala", "tools", "cmd"},
		Entity:    "Spec",
		Member:    "InfoextendsAnyRef",
		Signature: "",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.tools.ant.FastScalac"
	expected = entry{
		Namespace: []string{"scala", "tools", "ant"},
		Entity:    "FastScalac",
		Member:    "",
		Signature: "",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	v = "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]"
	expected = entry{
		Namespace: []string{"scala", "collection", "MapLike"},
		Entity:    "FilteredKeys",
		Member:    "andThen",
		Signature: "[C](k:B=>C):PartialFunction[A,C]",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.util.Success@isFailure:Boolean"
	expected = entry{
		Namespace: []string{"scala", "util"},
		Entity:    "Success",
		Member:    "isFailure",
		Signature: ":Boolean",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	v = "package"
	expected = entry{
		Namespace: []string{},
		Entity:    "package",
		Target:    target + v,
		Source:    source,
	}

	actual, _ = parseEntry(source, target, v)
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}
}
