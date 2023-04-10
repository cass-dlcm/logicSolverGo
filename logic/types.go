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

package logic

import (
	"fmt"
	"log"
	"reflect"
)

type Component interface {
	Calculate() bool
	Analyze() ([]string, [][]bool, int)
	Name() string
	GetChildren() ComponentList
	SolveChildren(map[string]int, *[]bool) bool
}

type ComponentList []Component

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

type BoolContainer struct {
	Letter string
	Value  *bool
}

func (b BoolContainer) Name() string {
	return b.Letter
}

func (b BoolContainer) Calculate() bool {
	return *(b.Value)
}

func (b BoolContainer) Analyze() ([]string, [][]bool, int) {
	return []string{b.Letter}, [][]bool{{false}, {true}}, 1
}

func (b BoolContainer) GetChildren() ComponentList {
	return ComponentList{b}
}

func (b BoolContainer) SolveChildren(names map[string]int, values *[]bool) bool {
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

type UnaryOperator int

const (
	NOT UnaryOperator = iota
)

type UnaryOperation struct {
	Operator UnaryOperator
	Value    Component
}

func (u UnaryOperation) Calculate() bool {
	switch u.Operator {
	case NOT:
		return !u.Value.Calculate()
	}
	return false
}

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

func (u UnaryOperation) Name() string {
	switch u.Operator {
	case NOT:
		return fmt.Sprintf("!(%s)", u.Value.Name())

	}
	return ""
}

func (u UnaryOperation) GetChildren() ComponentList {
	children := u.Value.GetChildren()
	return append(children, u)
}

func (u UnaryOperation) SolveChildren(names map[string]int, values *[]bool) bool {
	if names[u.Name()] == 0 {
		val := u.Value.SolveChildren(names, values)
		u.Value = BoolContainer{u.Value.Name(), &val}
		names[u.Name()] = len(names) + 1
		*values = append(*values, u.Calculate())
	} else if names[u.Name()] >= len(*values) {
		val := u.Value.SolveChildren(names, values)
		u.Value = BoolContainer{u.Value.Name(), &val}
		*values = append(*values, u.Calculate())
	}
	log.Println(names, *values)
	return (*values)[names[u.Name()]-1]
}

type BinaryOperator int

const (
	AND BinaryOperator = iota
	OR
)

type BinaryOperation struct {
	ValueA   Component
	ValueB   Component
	Operator BinaryOperator
}

func (b BinaryOperation) Calculate() bool {
	switch b.Operator {
	case AND:
		return b.ValueA.Calculate() && b.ValueB.Calculate()
	case OR:
		return b.ValueA.Calculate() || b.ValueB.Calculate()
	}
	return false
}

func (b BinaryOperation) Analyze() ([]string, [][]bool, int) {
	var values [][]bool
	var tempBool bool
	componentList := b.GetChildren()
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
		b.SolveChildren(names, &values[i])
	}
	log.Println(names)
	for i, j := range names {
		log.Println(j)
		componentTableHeadings[j-1] = i
	}
	return componentTableHeadings, values, count
}

func (b BinaryOperation) SolveChildren(names map[string]int, values *[]bool) bool {
	if names[b.Name()] == 0 {
		valA := b.ValueA.SolveChildren(names, values)
		b.ValueA = BoolContainer{b.ValueA.Name(), &valA}
		valB := b.ValueB.SolveChildren(names, values)
		b.ValueB = BoolContainer{b.ValueB.Name(), &valB}
		names[b.Name()] = len(names) + 1
		*values = append(*values, b.Calculate())
	} else if names[b.Name()] >= len(*values) {
		valA := b.ValueA.SolveChildren(names, values)
		b.ValueA = BoolContainer{b.ValueA.Name(), &valA}
		valB := b.ValueB.SolveChildren(names, values)
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

func (b BinaryOperation) GetChildren() ComponentList {
	childrenA := b.ValueA.GetChildren()
	childrenB := b.ValueB.GetChildren()
	children := append(childrenA, childrenB...)
	return append(children, b)
}

func (b BinaryOperation) Name() string {
	switch b.Operator {
	case AND:
		return fmt.Sprintf("(%s && %s)", b.ValueA.Name(), b.ValueB.Name())
	case OR:
		return fmt.Sprintf("(%s ||  %s)", b.ValueA.Name(), b.ValueB.Name())
	}
	return ""
}

func RemoveIrrelevantVariables(header []string, values [][]bool, variableCount int) ([]string, [][]bool) {
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
