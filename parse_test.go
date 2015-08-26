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
		Entity:    "Extractors$SymDef$",
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
		Entity:    "Extractors$SymDef$",
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
		Entity:    "Spec$",
		Function:  "InfoextendsAnyRef",
		Signature: "",
	}

	actual, _ = parseEntry(source, target, "scala.tools.cmd.Spec$@InfoextendsAnyRef")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
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
}
