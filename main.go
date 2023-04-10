/*
LogicSolverGo solves boolean logic expressions and generates truth tables.
Copyright (C) 2023  Cassandra de la Cruz-Munoz <me@cass-dlcm.dev>

        This program is free software: you can redistribute it and/or modify
        it under the terms of the GNU Affero General Public License as
        published by the Free Software Foundation, either version 3 of the
        License, or any later version.

        This program is distributed in the hope that it will be useful,
        but WITHOUT ANY WARRANTY; without even the implied warranty of
        MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
        GNU Affero General Public License for more details.

        You should have received a copy of the GNU Affero General Public License
        along with this program.  If not, see <https://www.gnu.org/licenses/>;.
*/

package main

import (
	"log"
	"logicSolverGo/logic"
)

func main() {
	a := true
	b := false
	c := false
	log.Println(logic.RemoveIrrelevantTerms(logic.BinaryOperation{
		ValueA: logic.UnaryOperation{Operator: logic.NOT, Value: logic.BoolContainer{
			Letter: "a",
			Value:  &a,
		}},
		ValueB: logic.BinaryOperation{
			ValueA: logic.BoolContainer{
				Letter: "c",
				Value:  &c,
			}, ValueB: logic.BoolContainer{
				Letter: "b",
				Value:  &b,
			}, Operator: logic.AND}, Operator: logic.OR}.Analyze()))
	log.Println(logic.RemoveIrrelevantTerms(logic.BinaryOperation{
		ValueA: logic.UnaryOperation{logic.NOT, logic.BoolContainer{
			"a", &a,
		}},
		ValueB:   logic.BoolContainer{"a", &a},
		Operator: logic.OR}.Analyze()))
	log.Println(logic.RemoveIrrelevantTerms(logic.BinaryOperation{
		ValueA: logic.UnaryOperation{logic.NOT, logic.BoolContainer{
			"a", &a,
		}},
		ValueB:   logic.BoolContainer{"a", &a},
		Operator: logic.AND}.Analyze()))
}
