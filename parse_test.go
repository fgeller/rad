package main

import "testing"

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	expected := entry{
		Namespace: []string{"scala", "reflect", "macros", "contexts"},
		Entity:    "Parsers",
		Function:  "notify",
		Signature: "():Unit",
	}

	actual, _ := parseEntry(source, target, "scala.reflect.macros.contexts.Parsers@notify():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got %v expected %v.", actual, expected)
	}

	// scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils"},
		Entity:    "SymDef",
		Function:  "notifyAll",
		Signature: "():Unit",
	}

	actual, _ = parseEntry(source, target, "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]
	expected = entry{
		Namespace: []string{"scala", "reflect", "reify", "utils"},
		Entity:    "SymDef",
		Function:  "unapply",
		Signature: "(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]",
	}

	actual, _ = parseEntry(source, target, "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.AnyRef@notify():Unit
	expected = entry{
		Namespace: []string{"scala"},
		Entity:    "AnyRef",
		Function:  "notify",
		Signature: "():Unit",
	}

	actual, _ = parseEntry(source, target, "scala.AnyRef@notify():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.tools.cmd.Spec$@InfoextendsAnyRef
	expected = entry{
		Namespace: []string{"scala", "tools", "cmd"},
		Entity:    "Spec",
		Function:  "InfoextendsAnyRef",
		Signature: "",
	}

	actual, _ = parseEntry(source, target, "scala.tools.cmd.Spec$@InfoextendsAnyRef")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	// scala.tools.ant.FastScalac
	expected = entry{
		Namespace: []string{"scala", "tools", "ant"},
		Entity:    "FastScalac",
		Function:  "",
		Signature: "",
	}

	actual, _ = parseEntry(source, target, "scala.tools.ant.FastScalac")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]
	expected = entry{
		Namespace: []string{"scala", "collection"},
		Entity:    "FilteredKeys",
		Function:  "andThen",
		Signature: "[C](k:B=>C):PartialFunction[A,C]",
	}

	actual, _ = parseEntry(source, target, "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	// scala.util.Success@isFailure:Boolean
	expected = entry{
		Namespace: []string{"scala", "util"},
		Entity:    "Success",
		Function:  "isFailure",
		Signature: ":Boolean",
	}

	actual, _ = parseEntry(source, target, "scala.util.Success@isFailure:Boolean")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}

	// package
	expected = entry{
		Namespace: []string{},
		Entity:    "package",
	}

	actual, _ = parseEntry(source, target, "package")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v\nexpected\n%v.", actual, expected)
	}
}

// TODO: MapLike$DefaultValuesIterable should be DefaultValuesIterable
// TODO: Map$ should be Map
