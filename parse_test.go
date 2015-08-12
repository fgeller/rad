package main

import "testing"

func TestParseScalaEntry(t *testing.T) {

	expected := entry{
		namespace: []string{"scala", "reflect", "macros", "contexts"},
		entity:    "Parsers",
		function:  "notify",
		signature: "():Unit",
	}

	actual, _ := parseEntry("scala.reflect.macros.contexts.Parsers@notify():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got %v expected %v.", actual, expected)
	}

	// scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit
	expected = entry{
		namespace: []string{"scala", "reflect", "reify", "utils"},
		entity:    "Extractors$SymDef$",
		function:  "notifyAll",
		signature: "():Unit",
	}

	actual, _ = parseEntry("scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]
	expected = entry{
		namespace: []string{"scala", "reflect", "reify", "utils"},
		entity:    "Extractors$SymDef$",
		function:  "unapply",
		signature: "(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]",
	}

	actual, _ = parseEntry("scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.AnyRef@notify():Unit
	expected = entry{
		namespace: []string{"scala"},
		entity:    "AnyRef",
		function:  "notify",
		signature: "():Unit",
	}

	actual, _ = parseEntry("scala.AnyRef@notify():Unit")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.tools.cmd.Spec$@InfoextendsAnyRef
	expected = entry{
		namespace: []string{"scala", "tools", "cmd"},
		entity:    "Spec$",
		function:  "InfoextendsAnyRef",
		signature: "",
	}

	actual, _ = parseEntry("scala.tools.cmd.Spec$@InfoextendsAnyRef")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}

	// scala.tools.ant.FastScalac
	expected = entry{
		namespace: []string{"scala", "tools", "ant"},
		entity:    "FastScalac",
		function:  "",
		signature: "",
	}

	actual, _ = parseEntry("scala.tools.ant.FastScalac")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}
}
