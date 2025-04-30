package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

func main() {
	// Assuming you have a toml file named 'data.toml'
	tomlFile, err := os.Open("data3.toml")
	if err != nil {
		fmt.Println("Error: data.toml not found.")
		return
	}
	defer tomlFile.Close()

	data, err := toml.LoadReader(tomlFile)
	if err != nil {
		fmt.Println("Error loading TOML data:", err)
		return
	}

	plantumlCode := createPlantUmlDiagram(data)

	// Save the PlantUML code to a file
	plantumlFile, err := os.Create("diagram3.puml")
	if err != nil {
		fmt.Println("Error creating diagram.puml:", err)
		return
	}
	defer plantumlFile.Close()

	_, err = plantumlFile.WriteString(plantumlCode)
	if err != nil {
		fmt.Println("Error writing to diagram.puml:", err)
		return
	}

	fmt.Println("Diagram generated successfully!")
}

// SourceObject represents a class in the object diagram
type SourceObject struct {
	Name string
	Kind string
	Type string
}

// TransformObject represents a class in the object diagram
type TransformObject struct {
	Name   string
	Kind   string
	Type   string
	Inputs []string
}

// SinkObject represents a class in the object diagram
type SinkObject struct {
	Name   string
	Kind   string
	Type   string
	Inputs []string
}

func createPlantUmlDiagram(data *toml.Tree) string {

	var transformObjs []TransformObject
	var sinkObjs []SinkObject
	var srcObjs []SourceObject
	// Create a simple PlantUML class diagram
	umlCode := "@startuml\n"
	umlCode += "left to right direction\n"
	umlCode += "title Data Processing Pipeline\n"
	umlCode += "\n"

	for _, section := range data.Keys() {
		sectionData, ok := data.Get(section).(*toml.Tree)
		if !ok {
			continue
		}
		name := fmt.Sprintf("%s", section)
		switch name { // asuming this hits only once per kind
		case "sources":
			fmt.Println("found sources", name)
			srcObjs = createSourceObject(sectionData)
			srcStr := convertSourceObjectsToString(srcObjs)
			fmt.Println(srcStr)
			umlCode += srcStr
		case "transforms":
			fmt.Println("found transforms", name)
			transformObjs = createTransformObject(sectionData)
			transformStr := convertTransformObjectsToString(transformObjs)
			fmt.Println(transformStr)
			umlCode += transformStr
		case "sinks":
			fmt.Println("found sinks", name)
			sinkObjs = createSinksObject(sectionData)
			sinkStr := convertSinkObjectsToString(sinkObjs)
			fmt.Println(sinkStr)
			umlCode += sinkStr
		default:
			fmt.Println("ignoring", name)
		}
	}

	mappingStr := createObjectsMapping(srcObjs, transformObjs, sinkObjs)
	fmt.Println(mappingStr)
	umlCode += mappingStr

	umlCode += "@enduml\n"
	return umlCode
}

func createSourceObject(obj *toml.Tree) []SourceObject {
	var sourceObjects []SourceObject
	for _, name := range obj.Keys() {
		subSection, ok := obj.Get(name).(*toml.Tree)
		if !ok {
			continue
		}
		typeName := subSection.Get("type").(string)
		sourceObjects = append(sourceObjects, SourceObject{
			Name: name,
			Type: typeName,
			Kind: "sources",
		})
	}
	return sourceObjects
}

func convertSourceObjectsToString(srcObjs []SourceObject) string {
	var result string
	result += "together {\n"
	for _, srcObj := range srcObjs {
		result += fmt.Sprintf("object \"%s\"\n", srcObj.Name)
		result += fmt.Sprintf("\"%s\" : Name = %s\n", srcObj.Name, srcObj.Name)
		result += fmt.Sprintf("\"%s\" : Kind = %s\n", srcObj.Name, srcObj.Kind)
		result += fmt.Sprintf("\"%s\" : Type = %s\n", srcObj.Name, srcObj.Type)
		result += "\n"
	}
	result += "}\n"
	return result
}

func createTransformObject(obj *toml.Tree) []TransformObject {
	var trnsObjects []TransformObject
	for _, name := range obj.Keys() {
		subSection, ok := obj.Get(name).(*toml.Tree)
		if !ok {
			continue
		}
		typeName := subSection.Get("type").(string)
		inputs, _ := subSection.GetArray("inputs").([]string)
		trnsObjects = append(trnsObjects, TransformObject{
			Name:   name,
			Type:   typeName,
			Kind:   "transforms",
			Inputs: inputs,
		})
	}
	return trnsObjects
}

func convertTransformObjectsToString(transformObjs []TransformObject) string {
	var result string
	for _, transformObj := range transformObjs {
		result += fmt.Sprintf("object \"%s\"\n", transformObj.Name)
		result += fmt.Sprintf("\"%s\" : Name = %s\n", transformObj.Name, transformObj.Name)
		result += fmt.Sprintf("\"%s\" : Kind = %s\n", transformObj.Name, transformObj.Kind)
		result += fmt.Sprintf("\"%s\" : Type = %s\n", transformObj.Name, transformObj.Type)
		result += "\n"
	}
	return result
}

func createSinksObject(obj *toml.Tree) []SinkObject {
	var trnsObjects []SinkObject
	for _, name := range obj.Keys() {
		subSection, ok := obj.Get(name).(*toml.Tree)
		if !ok {
			continue
		}
		typeName := subSection.Get("type").(string)
		inputs, _ := subSection.GetArray("inputs").([]string)
		trnsObjects = append(trnsObjects, SinkObject{
			Name:   name,
			Type:   typeName,
			Kind:   "sinks",
			Inputs: inputs,
		})
	}
	return trnsObjects
}

func convertSinkObjectsToString(sinkObjs []SinkObject) string {
	var result string
	//result += "together {\n"
	for _, sinkObj := range sinkObjs {
		result += fmt.Sprintf("object \"%s\"\n", sinkObj.Name)
		result += fmt.Sprintf("\"%s\" : Name = %s\n", sinkObj.Name, sinkObj.Name)
		result += fmt.Sprintf("\"%s\" : Kind = %s\n", sinkObj.Name, sinkObj.Kind)
		result += fmt.Sprintf("\"%s\" : Type = %s\n", sinkObj.Name, sinkObj.Type)
		result += "\n"
	}
	//result += "}\n"
	return result
}

func createObjectsMapping(srcObjs []SourceObject, transformObjs []TransformObject, sinkObjs []SinkObject) string {
	var mapping string
	srcNames := getArrayOfSourceOnjNames(srcObjs)
	trnsNames := getArrayOfTransformOnjNames(transformObjs)
	mapping = "\n"

	mapping += getTransformMappings(transformObjs, srcNames, trnsNames)
	mapping += "\n"

	mapping += getSinkMappings(sinkObjs, srcNames, trnsNames)
	mapping += "\n"

	return mapping
}

func getTransformMappings(transformObjs []TransformObject, srcNames []string, trnsNames []string) string {
	var mapping string
	for _, transformObj := range transformObjs {
		for _, source := range transformObj.Inputs {
			src := findMatchingElements(srcNames, source) // match sources
			for _, s := range src {
				tr := transformObj.Name
				mapping += fmt.Sprintf("\"%s\" --[#green]> \"%s\"\n", s, tr)
			}

			trnsNames := findMatchingElements(trnsNames, source) // match transforms
			for _, t := range trnsNames {
				tr := transformObj.Name
				mapping += fmt.Sprintf("\"%s\" ..[#blue]> \"%s\"\n", t, tr)
			}
		}
	}
	return mapping
}

func getSinkMappings(sinkObjs []SinkObject, srcNames []string, trnsNames []string) string {
	var mapping string
	for _, sinkObj := range sinkObjs {
		for _, sinkSource := range sinkObj.Inputs {
			src := findMatchingElements(srcNames, sinkSource) // match sources
			for _, s := range src {
				sn := sinkObj.Name
				mapping += fmt.Sprintf("\"%s\" --[#red]> \"%s\"\n", s, sn)
			}

			trnsNames := findMatchingElements(trnsNames, sinkSource) // match transforms
			for _, t := range trnsNames {
				sn := sinkObj.Name
				mapping += fmt.Sprintf("\"%s\" --[#red]> \"%s\"\n", t, sn)
			}
		}
	}
	return mapping
}

func findMatchingElements(arr []string, criteria string) []string {
	var result []string
	if strings.HasSuffix(criteria, "*") {
		prefix := strings.TrimSuffix(criteria, "*")
		for _, elem := range arr {
			if strings.HasPrefix(elem, prefix) {
				result = append(result, elem)
			}
		}
	} else {
		for _, elem := range arr {
			if elem == criteria {
				result = append(result, elem)
			}
		}
	}
	return result
}

func getArrayOfSourceOnjNames(srcObjs []SourceObject) []string {
	var srcNames []string
	for _, src := range srcObjs {
		srcNames = append(srcNames, src.Name)
	}
	return srcNames
}

func getArrayOfTransformOnjNames(trObjs []TransformObject) []string {
	var trNames []string
	for _, tr := range trObjs {
		trNames = append(trNames, tr.Name)
	}
	return trNames
}
