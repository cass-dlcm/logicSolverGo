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

// Package logic implements the functionality of the library.
package logic

import (
	"fmt"
	"log"
	"reflect"
)

// Component is an interface listing the common functionality of what all logical components provide.
// Calculate gives an answer.
// Analyze generates the truth table.
// Name identifies the component.
type Component interface {
	Calculate() bool
	Analyze() ([]string, [][]bool, int)
	Name() string
	getChildren() ComponentList
	solveChildren(map[string]int, *[]bool) bool
}

// ComponentList is a slice of Component.
type ComponentList []Component

// Deduplicate keeps only one of each Component with the same name.
func (c *ComponentList) Deduplicate() {
	var tempC ComponentList
	for i := range *c {
		tempC = append(tempC, (*c)[i])
		log.Println((*c)[i].Name())
	}
	for i := 0; i < len(tempC); {
		log.Println(i)
		if tempC[i] == nil {
			tempC = append(tempC[:i], tempC[i+1:]...)
			i--
			continue
		}
		for j := i + 1; j < len(tempC); {
			if tempC[i] == tempC[j] {
				tempC = append(tempC[:j], tempC[j+1:]...)
			} else {
				j++
			}
		}
		i++
	}
	*c = tempC
}

// RemoveBoolContainers removes the BoolContainer from the ComponentList, leaving only operations.
func (c *ComponentList) RemoveBoolContainers() {
	var tempC ComponentList
	for i := range *c {
		if reflect.TypeOf((*c)[i]) != reflect.TypeOf(BoolContainer{
			Letter: "",
			Value:  nil,
		}) {
			tempC = append(tempC, (*c)[i])
		}
	}
	*c = tempC
}

// BoolContainer provides a Component to store bool primitives.
type BoolContainer struct {
	Letter string
	Value  *bool
}

// Name outputs the label of the BoolContainer.
func (b BoolContainer) Name() string {
	return b.Letter
}

// Calculate simply returns the stored value.
func (b BoolContainer) Calculate() bool {
	return *(b.Value)
}

// Analyze creates a two-row one-column truth table, of false and true.
func (b BoolContainer) Analyze() ([]string, [][]bool, int) {
	return []string{b.Letter}, [][]bool{{false}, {true}}, 1
}

func (b BoolContainer) getChildren() ComponentList {
	return ComponentList{b}
}

func (b BoolContainer) solveChildren(names map[string]int, values *[]bool) bool {
	log.Println(names, len(*values))
	if names[b.Name()] == 0 {
		names[b.Name()] = len(names) + 1
		*values = append(*values, b.Calculate())
	} else if names[b.Name()] > len(*values) {
		*values = append(*values, b.Calculate())
	}
	log.Println(names, *values)
	return (*values)[names[b.Name()]-1]
}

// UnaryOperator is a number to represent unary operations.
type UnaryOperator int

// NOT is the UnaryOperator generally represented in code as "!"
const (
	NOT UnaryOperator = iota
)

// UnaryOperation is the Component that uses an UnaryOperator and a single Component.
type UnaryOperation struct {
	Operator UnaryOperator
	Value    Component
}

// Calculate checks the operator and for NOT, it takes the negation of the [Component.Calculate] function.
func (u UnaryOperation) Calculate() bool {
	switch u.Operator {
	case NOT:
		return !u.Value.Calculate()
	}
	return false
}

// Analyze generates a truth table containing 2^n rows where n is the number of descendent BoolContainer, and the number of columns is the number of unique Component.
func (u UnaryOperation) Analyze() ([]string, [][]bool, int) {
	componentTableHeadingA, values, varCount := u.Value.Analyze()
	componentTableHeadingA = append(componentTableHeadingA, u.Name())
	componentTableHeadingFound := map[string]bool{}
	var componentTableHeadings []string
	for i := range componentTableHeadingA {
		if !componentTableHeadingFound[componentTableHeadingA[i]] {
			componentTableHeadings = append(componentTableHeadings, componentTableHeadingA[i])
			componentTableHeadingFound[componentTableHeadingA[i]] = true
		}
	}
	for i := range values {
		values[i] = append(values[i], UnaryOperation{u.Operator,
			BoolContainer{componentTableHeadingA[len(componentTableHeadingA)-1],
				&values[i][len(values[i])-1]}}.Calculate())
	}
	return componentTableHeadings, values, varCount
}

// Name recursively constructs the name of the UnaryOperation from the UnaryOperator and child [Component.Name].
func (u UnaryOperation) Name() string {
	switch u.Operator {
	case NOT:
		return fmt.Sprintf("!(%s)", u.Value.Name())

	}
	return ""
}

func (u UnaryOperation) getChildren() ComponentList {
	children := u.Value.getChildren()
	return append(children, u)
}

func (u UnaryOperation) solveChildren(names map[string]int, values *[]bool) bool {
	if names[u.Name()] == 0 {
		val := u.Value.solveChildren(names, values)
		u.Value = BoolContainer{u.Value.Name(), &val}
		names[u.Name()] = len(names) + 1
		*values = append(*values, u.Calculate())
	} else if names[u.Name()] >= len(*values) {
		val := u.Value.solveChildren(names, values)
		u.Value = BoolContainer{u.Value.Name(), &val}
		*values = append(*values, u.Calculate())
	}
	log.Println(names, *values)
	return (*values)[names[u.Name()]-1]
}

// BinaryOperator is a number to represent binary operations.
type BinaryOperator int

// AND is the && BinaryOperator.
// OR is the || BinaryOperator.
const (
	AND BinaryOperator = iota
	OR
)

// BinaryOperation is the Component that uses a BinaryOperator and a pair of Component.
type BinaryOperation struct {
	ValueA   Component
	ValueB   Component
	Operator BinaryOperator
}

// Calculate returns the calculation of the BinaryOperator with the values of each Component operand.
func (b BinaryOperation) Calculate() bool {
	switch b.Operator {
	case AND:
		return b.ValueA.Calculate() && b.ValueB.Calculate()
	case OR:
		return b.ValueA.Calculate() || b.ValueB.Calculate()
	}
	return false
}

// Analyze generates a truth table consisting of 2^n rows, where n is the number of unique BoolContainer, and m columns where m is the number of unique Component in the list of components.
func (b BinaryOperation) Analyze() ([]string, [][]bool, int) {
	var values [][]bool
	var tempBool bool
	componentList := b.getChildren()
	log.Println(componentList)
	componentList.Deduplicate()
	log.Println(componentList)
	componentTableHeadings := make([]string, len(componentList))
	var variableIndices []int
	names := map[string]int{}
	count := 0
	for i := range componentList {
		if reflect.TypeOf(componentList[i]) == reflect.TypeOf(BoolContainer{"", &tempBool}) {
			variableIndices = append(variableIndices, i)
			componentTableHeadings[count] = componentList[i].Name()
			names[componentList[i].Name()] = count + 1
			count++
		}
	}
	values = make([][]bool, 1<<len(variableIndices))
	for i := range values {
		var row []bool
		for j := 0; j < len(variableIndices); j++ {
			row = append(row, i>>j%2 == 1)
		}
		values[i] = row
	}
	log.Println(len(componentTableHeadings))
	for i := range values {
		b.solveChildren(names, &values[i])
	}
	log.Println(names)
	for i, j := range names {
		log.Println(j)
		componentTableHeadings[j-1] = i
	}
	return componentTableHeadings, values, count
}

func (b BinaryOperation) solveChildren(names map[string]int, values *[]bool) bool {
	if names[b.Name()] == 0 {
		valA := b.ValueA.solveChildren(names, values)
		b.ValueA = BoolContainer{b.ValueA.Name(), &valA}
		valB := b.ValueB.solveChildren(names, values)
		b.ValueB = BoolContainer{b.ValueB.Name(), &valB}
		names[b.Name()] = len(names) + 1
		*values = append(*values, b.Calculate())
	} else if names[b.Name()] >= len(*values) {
		valA := b.ValueA.solveChildren(names, values)
		b.ValueA = BoolContainer{b.ValueA.Name(), &valA}
		valB := b.ValueB.solveChildren(names, values)
		b.ValueB = BoolContainer{b.ValueB.Name(), &valB}
		if names[b.Name()] >= len(*values) {
			*values = append(*values, b.Calculate())
		} else {
			(*values)[names[b.Name()]-1] = b.Calculate()
		}
	}
	log.Println(names, *values)
	return (*values)[names[b.Name()]-1]
}

func (b BinaryOperation) getChildren() ComponentList {
	childrenA := b.ValueA.getChildren()
	childrenB := b.ValueB.getChildren()
	children := append(childrenA, childrenB...)
	return append(children, b)
}

// Name recursively constructs the name of the BinaryOperation from the BinaryOperator and the two child [Component.Name].
func (b BinaryOperation) Name() string {
	switch b.Operator {
	case AND:
		return fmt.Sprintf("(%s && %s)", b.ValueA.Name(), b.ValueB.Name())
	case OR:
		return fmt.Sprintf("(%s ||  %s)", b.ValueA.Name(), b.ValueB.Name())
	}
	return ""
}

// RemoveIrrelevantTerms simplifies a truth table to only all terms that affect the output.
// Tautologies and contradictions are reduced to a 1x1 [][]bool, containing true for a tautology and false for a contradiction.
func RemoveIrrelevantTerms(header []string, values [][]bool, variableCount int) ([]string, [][]bool) {
	log.Println(header, values)
	i := 0
	for i < variableCount {
		j := 0
		count := 0
		symmetric := true
		for j+1<<i < 1<<variableCount {
			index := j + 1<<i
			log.Println(i, j, index)
			log.Println(values[j][len(values[j])-1], values[index][len(values[index])-1])
			symmetric = values[j][len(values[j])-1] == values[index][len(values[index])-1] && symmetric
			if !symmetric {
				break
			}
			count++
			if count%1<<i == 0 {
				j++
			} else {
				j += 1 << (i + 1)
			}
		}
		if symmetric == true {
			if i < len(values[j])-1 {
				for j := range values {
					values[j] = append(values[j][:i], values[j][i+1:]...)
				}
				header = append(header[:i], header[i+1:]...)
			} else {
				return []string{"Value"}, [][]bool{{values[i][0]}}
			}
		} else {
			i++
		}
	}
	return header, values
}
