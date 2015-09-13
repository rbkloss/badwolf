// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package semantic

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/badwolf/bql/lexer"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/google/badwolf/triple/predicate"
)

func TestDataAccumulatorHook(t *testing.T) {
	st := &Statement{}
	ces := []ConsumedElement{
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemNode,
			Text: "/_<s>",
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemPredicate,
			Text: `"p"@[]`,
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemNode,
			Text: "/_<o>",
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemNode,
			Text: "/_<s>",
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemPredicate,
			Text: `"p"@[]`,
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemNode,
			Text: "/_<o>",
		}),
	}
	var (
		hook ElementHook
		err  error
	)
	hook = dataAccumulator(literal.DefaultBuilder())
	for _, ce := range ces {
		hook, err = hook(st, ce)
		if err != nil {
			t.Errorf("semantic.DataAccumulator hook should have never failed for %v with error %v", ce, err)
		}
	}
	data := st.Data()
	if len(data) != 2 {
		t.Errorf("semantic.DataAccumulator hook should have produced 2 triples; instead produced %v", st.Data())
	}
	for _, trpl := range data {
		if got, want := trpl.S().String(), "/_<s>"; got != want {
			t.Errorf("semantic.DataAccumulator hook failed to parse subject correctly; got %v, want %v", got, want)
		}
		if got, want := trpl.P().String(), `"p"@[]`; got != want {
			t.Errorf("semantic.DataAccumulator hook failed to parse prdicate correctly; got %v, want %v", got, want)
		}
		if got, want := trpl.O().String(), "/_<o>"; got != want {
			t.Errorf("semantic.DataAccumulator hook failed to parse object correctly; got %v, want %v", got, want)
		}
	}
}

func TestSemanticAcceptInsertDelete(t *testing.T) {
	st := &Statement{}
	ces := []ConsumedElement{
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemBinding,
			Text: "?foo",
		}),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemComma,
			Text: ",",
		}),
		NewConsumedSymbol("FOO"),
		NewConsumedToken(&lexer.Token{
			Type: lexer.ItemBinding,
			Text: "?bar",
		}),
	}
	var (
		hook ElementHook
		err  error
	)
	hook = graphAccumulator()
	for _, ce := range ces {
		hook, err = hook(st, ce)
		if err != nil {
			t.Errorf("semantic.GraphAccumulator hook should have never failed for %v with error %v", ce, err)
		}
	}
	data := st.Graphs()
	if len(data) != 2 {
		t.Errorf("semantic.GraphAccumulator hook should have produced 2 graph bindings; instead produced %v", st.Graphs())
	}
	for _, g := range data {
		if g != "?foo" && g != "?bar" {
			t.Errorf("semantic.GraphAccumulator hook failed to provied either ?foo or ?bar; got %v instead", g)
		}
	}
}

func TestTypeBindingClauseHook(t *testing.T) {
	f := TypeBindingClauseHook(Insert)
	st := &Statement{}
	f(st, Symbol("FOO"))
	if got, want := st.Type(), Insert; got != want {
		t.Errorf("semantic.TypeBidingHook failed to set the right type; got %s, want %s", got, want)
	}
}

func TestWhereInitClauseHook(t *testing.T) {
	f := whereInitWorkingClause()
	st := &Statement{}
	f(st, Symbol("FOO"))
	if st.WorkingClause() == nil {
		t.Errorf("semantic.WhereInitWorkingClause should have returned a valid working clause for statement %v", st)
	}
}

func TestWhereWorkingClauseHook(t *testing.T) {
	f := whereNextWorkingClause()
	st := &Statement{}
	st.ResetWorkingGraphClause()
	f(st, Symbol("FOO"))
	f(st, Symbol("FOO"))
	if len(st.GraphPatternClauses()) != 2 {
		t.Errorf("semantic.whereNextWorkingClause should have returned two clauses for statement %v", st)
	}
}

func TestWhereSubjectClauseHook(t *testing.T) {
	st := &Statement{}
	f := whereSubjectClause("S", "SBinding", "SAlias", "STypeAlias", "SIDAlias")
	st.ResetWorkingGraphClause()
	n, err := node.Parse("/_<foo>")
	if err != nil {
		t.Fatalf("node.Parse failed with error %v", err)
	}
	table := []struct {
		ces  []ConsumedElement
		want *GraphClause
	}{
		{
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemNode,
					Text: "/_<foo>",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemType,
					Text: "type",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				S:          n,
				SAlias:     "?bar",
				STypeAlias: "?bar2",
				SIDAlias:   "?bar3",
			},
		},
		{
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?foo",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemType,
					Text: "type",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				SBinding:   "?foo",
				SAlias:     "?bar",
				STypeAlias: "?bar2",
				SIDAlias:   "?bar3",
			},
		},
	}
	for _, entry := range table {
		for _, ce := range entry.ces {
			if _, err := f(st, ce); err != nil {
				t.Errorf("semantic.whereSubjectClause should have never failed with error: %v", err)
			}
		}
		if got, want := st.WorkingClause(), entry.want; !reflect.DeepEqual(got, want) {
			t.Errorf("smeantic.whereSubjectClause should have populated all subject fields; got %+v, want %+v", got, want)
		}
		st.ResetWorkingGraphClause()
	}
}

func TestWherePredicatClauseHook(t *testing.T) {
	st := &Statement{}
	f := wherePredicateClause(
		"P", "PAlias", "PID", "PAnchorBinding", "PBinding", "PLowerBound", "PUpperBound",
		"PLowerBoundAlias", "PUpperBoundAlias", "PIDAlias", "PAnchorAlias")
	st.ResetWorkingGraphClause()
	p, err := predicate.Parse(`"foo"@[2015-07-19T13:12:04.669618843-07:00]`)
	if err != nil {
		t.Fatalf("predicate.Parse failed with error %v", err)
	}
	tlb, err := time.Parse(time.RFC3339Nano, `2015-07-19T13:12:04.669618843-07:00`)
	if err != nil {
		t.Fatalf("time.Parse failed to parse valid lower time bound with error %v", err)
	}
	tub, err := time.Parse(time.RFC3339Nano, `2016-07-19T13:12:04.669618843-07:00`)
	if err != nil {
		t.Fatalf("time.Parse failed to parse valid upper time bound with error %v", err)
	}
	table := []struct {
		valid bool
		id    string
		ces   []ConsumedElement
		want  *GraphClause
	}{
		{
			valid: true,
			id:    "valid predicate",
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemPredicate,
					Text: `"foo"@[2015-07-19T13:12:04.669618843-07:00]`,
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAt,
					Text: "at",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				P:            p,
				PAlias:       "?bar",
				PIDAlias:     "?bar2",
				PAnchorAlias: "?bar3",
			},
		},
		{
			valid: true,
			id:    "valid predicate with binding",
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemPredicate,
					Text: `"foo"@[?foo]`,
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAt,
					Text: "at",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				PID:            "foo",
				PAnchorBinding: "?foo",
				PAlias:         "?bar",
				PIDAlias:       "?bar2",
				PAnchorAlias:   "?bar3",
			},
		},
		{
			valid: true,
			id:    "valid bound with bindings",
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemPredicateBound,
					Text: `"foo"@[?fooLower,?fooUpper]`,
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAt,
					Text: "at",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				PID:              "foo",
				PLowerBoundAlias: "?fooLower",
				PUpperBoundAlias: "?fooUpper",
				PAlias:           "?bar",
				PIDAlias:         "?bar2",
				PAnchorAlias:     "?bar3",
			},
		},
		{
			valid: true,
			id:    "valid bound with dates",
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemPredicateBound,
					Text: `"foo"@[2015-07-19T13:12:04.669618843-07:00,2016-07-19T13:12:04.669618843-07:00]`,
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAt,
					Text: "at",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{
				PID:          "foo",
				PLowerBound:  &tlb,
				PUpperBound:  &tub,
				PAlias:       "?bar",
				PIDAlias:     "?bar2",
				PAnchorAlias: "?bar3",
			},
		},
		{
			valid: false,
			id:    "invalid bound with dates",
			ces: []ConsumedElement{
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemPredicateBound,
					Text: `"foo"@[2016-07-19T13:12:04.669618843-07:00,2015-07-19T13:12:04.669618843-07:00]`,
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAs,
					Text: "as",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar",
				}),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemID,
					Text: "id",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar2",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemAt,
					Text: "at",
				}),
				NewConsumedSymbol("FOO"),
				NewConsumedToken(&lexer.Token{
					Type: lexer.ItemBinding,
					Text: "?bar3",
				}),
				NewConsumedSymbol("FOO"),
			},
			want: &GraphClause{},
		},
	}
	failed := false
	for _, entry := range table {
		for _, ce := range entry.ces {
			if _, err := f(st, ce); err != nil {
				if entry.valid {
					t.Errorf("semantic.wherePredicateClause case %q should have never failed with error: %v", entry.id, err)
				} else {
					failed = true
				}
			}
		}
		if entry.valid {
			if got, want := st.WorkingClause(), entry.want; !reflect.DeepEqual(got, want) {
				t.Errorf("semantic.wherePredicateClause case %q should have populated all subject fields; got %+v, want %+v", entry.id, got, want)
			}
		} else {
			if !failed {
				t.Errorf("semantic.wherePredicateClause failed to reject invalid case %q", entry.id)
			}
		}
		st.ResetWorkingGraphClause()
	}
}

func TestSetClauseStringField(t *testing.T) {
	c := &GraphClause{}
	setClauseStringField(c, "SBinding", "string value")
	if got, want := c.SBinding, "string value"; got != want {
		t.Errorf("setClauseStringField failed to set SBinding field in %+v; got %s, want %s", c, got, want)
	}
}

func TestSetClausePointerField(t *testing.T) {
	c := &GraphClause{}
	n, err := node.Parse("/foo<bar>")
	if err != nil {
		t.Fatalf("node.Parse failed to pars /foo<bar>")
	}
	setClausePointerField(c, "S", n)
	if got, want := c.S, n; !reflect.DeepEqual(got, want) {
		t.Errorf("setClausePointerField failed to set SBinding field in %+v; got %v, want %v", c, got, want)
	}
}
