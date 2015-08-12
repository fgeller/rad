package main

import "testing"

func TestParseScalaEntry(t *testing.T) {
	// scala.
	// reflect.
	// reify.
	// utils.
	// Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):
	// Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]

	// scala.
	// reflect.
	// reify.
	// utils.
	// Extractors$SymDef$@notifyAll():
	//Unit

	// scala.
	// AnyRef@notify():
	//Unit

	// scala.
	// reflect.
	// macros.
	// contexts.Parsers@notify():
	// Unit

	// scala.reflect.macros.contexts.Parsers@notify(): Unit

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
}
