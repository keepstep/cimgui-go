package main

import (
	"fmt"
	"os"
	"strings"
)

func generateGoStructs(prefix string, structs []StructDef) []string {
	var sb strings.Builder

	sb.WriteString(goPackageHeader)
	sb.WriteString(fmt.Sprintf(
		`// #include "%s_wrapper.h"
import "C"
import "unsafe"

`, prefix))

	// Save all struct name into a map
	var structNames []string

	for _, s := range structs {
		if shouldSkipStruct(s.Name) {
			continue
		}

		sb.WriteString(fmt.Sprintf("%s\n", s.CommentAbove))
		sb.WriteString(fmt.Sprintf(`type %[1]s uintptr

func (data %[1]s) handle() *C.%[2]s {
  return (*C.%[2]s)(unsafe.Pointer(data))
}

func (data %[1]s) c() C.%[2]s {
  return *(data.handle())
}

func new%[1]sFromC(cvalue C.%[2]s) %[1]s {
  return %[1]s(unsafe.Pointer(&cvalue))
}

`, renameGoIdentifier(s.Name), s.Name))

		structNames = append(structNames, s.Name)
	}

	structFile, err := os.Create(fmt.Sprintf("%s_structs.go", prefix))
	if err != nil {
		panic(err.Error())
	}

	defer structFile.Close()

	_, _ = structFile.WriteString(sb.String())

	return structNames
}
