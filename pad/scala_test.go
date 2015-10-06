package main

import (
	"../shared"
	"testing"
)

// TODO

// {
//   "Path": [
//     "scala",
//     "AnyRef"
//   ],
//   "Members": [
//     {
//       "Name": "==(x$1:Any):Boolean",
//       "Target": "scala/scala-docs-2.11.7/api/scala-reflect/index.html#scala.AnyRef@==(x$1:Any):Boolean"
//     }
//   ]
// },

// {
//   "Path": [
//     "package"
//   ],
//   "Members": [
//     {
//       "Name": "",
//       "Target": "scala/scala-docs-2.11.7/api/scala-reflect/index.html#package"
//     }
//   ]
// },

// merging:
// {
//   "Path": "package",
//   "Members": [
//     {
//       "Name": "",
//       "Target": "scala/scala-docs-2.11.7/api/scala-library/index.html#package"
//     }
//   ]
// },
// {
//   "Path": "package",
//   "Members": [
//     {
//       "Name": "scala",
//       "Target": "scala/scala-docs-2.11.7/api/scala-library/index.html#package@scala"
//     }
//   ]
// },

func TestParseScalaEntry(t *testing.T) {

	source := "/some source file"
	target := "some target"

	v := "scala.reflect.macros.contexts.Parsers@notify():Unit"
	expected := shared.Namespace{
		Path:    "scala.reflect.macros.contexts.Parsers",
		Members: []shared.Member{{Name: "notify", Target: "/" + target + v}},
	}

	actual, _ := parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@notifyAll():Unit"
	expected = shared.Namespace{
		Path:    "scala.reflect.reify.utils.Extractors.SymDef",
		Members: []shared.Member{{Name: "notifyAll", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.reflect.reify.utils.Extractors$SymDef$@unapply(tree:Extractors.this.global.Tree):Option[(Extractors.this.global.Tree,Extractors.this.global.TermName,Long,Boolean)]"
	expected = shared.Namespace{
		Path:    "scala.reflect.reify.utils.Extractors.SymDef",
		Members: []shared.Member{{Name: "unapply", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.AnyRef@notify():Unit"
	expected = shared.Namespace{
		Path:    "scala.AnyRef",
		Members: []shared.Member{{Name: "notify", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.tools.cmd.Spec$@InfoextendsAnyRef"
	expected = shared.Namespace{
		Path:    "scala.tools.cmd.Spec",
		Members: []shared.Member{{Name: "Info", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.tools.ant.FastScalac"
	expected = shared.Namespace{
		Path:    "scala.tools.ant.FastScalac",
		Members: []shared.Member{{Name: "", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v expected \n%v.", actual, expected)
	}

	v = "scala.collection.MapLike$FilteredKeys@andThen[C](k:B=>C):PartialFunction[A,C]"
	expected = shared.Namespace{
		Path:    "scala.collection.MapLike.FilteredKeys",
		Members: []shared.Member{{Name: "andThen", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.util.Success@isFailure:Boolean"
	expected = shared.Namespace{
		Path:    "scala.util.Success",
		Members: []shared.Member{{Name: "isFailure", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "package"
	expected = shared.Namespace{
		Path:    "package",
		Members: []shared.Member{{Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.concurrent.TrieMap$@Coll=CC[_,_]"
	expected = shared.Namespace{
		Path:    "scala.collection.concurrent.TrieMap",
		Members: []shared.Member{{Name: "Coll", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.Enumeration$ValueSet@:\\[B](z:B)(op:(A,B)=\u003eB):B"
	expected = shared.Namespace{
		Path:    "scala.Enumeration.ValueSet",
		Members: []shared.Member{{Name: ":\\", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.IterableViewLike$Transformed@Self=Repr"
	expected = shared.Namespace{
		Path:    "scala.collection.IterableViewLike.Transformed",
		Members: []shared.Member{{Name: "Self", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.MapLike@FilteredKeysextendsAbstractMap[A,B]withDefaultMap[A,B]"
	expected = shared.Namespace{
		Path:    "scala.collection.MapLike",
		Members: []shared.Member{{Name: "FilteredKeys", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}

	v = "scala.collection.parallel.FutureThreadPoolTasks$@==(x$1:Any):Boolean"
	expected = shared.Namespace{
		Path:    "scala.collection.parallel.FutureThreadPoolTasks",
		Members: []shared.Member{{Name: "==", Target: "/" + target + v}},
	}

	actual, _ = parseNamespace(source, target, v)
	if !expected.Eq(actual) {
		t.Errorf("parsing scala entry failed. got\n%v\nexpected\n%v.", actual, expected)
	}
}
