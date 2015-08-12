package main

import "testing"

func TestParseScalaEntry(t *testing.T) {

	// scala.AnyRef@notify():Unit

	expected := entry{
		namespace: []string{"scala", "reflect", "macros", "contexts"},
		entity:    "Parsers",
		function:  "notify",
		signature: "():Unit",
	}

	actual, _ := parseScalaEntry("scala.reflect.macros.contexts.Parsers@notify():Unit")
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

	actual, _ = parseScalaEntry("scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit")
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

	actual, _ = parseScalaEntry("scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]")
	if !expected.eq(actual) {
		t.Errorf("parsing scala entry failed. got \n%v expected \n%v.", actual, expected)
	}
}
