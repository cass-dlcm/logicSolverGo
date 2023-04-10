package main

import (
	"log"
	"logicSolverGo/logic"
)

func main() {
	a := true
	b := false
	c := false
	log.Println(logic.RemoveIrrelevantVariables(logic.BinaryOperation{
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
	log.Println(logic.RemoveIrrelevantVariables(logic.BinaryOperation{
		ValueA: logic.UnaryOperation{logic.NOT, logic.BoolContainer{
			"a", &a,
		}},
		ValueB:   logic.BoolContainer{"a", &a},
		Operator: logic.OR}.Analyze()))
	log.Println(logic.RemoveIrrelevantVariables(logic.BinaryOperation{
		ValueA: logic.UnaryOperation{logic.NOT, logic.BoolContainer{
			"a", &a,
		}},
		ValueB:   logic.BoolContainer{"a", &a},
		Operator: logic.AND}.Analyze()))
}
