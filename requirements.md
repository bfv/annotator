# Annotator project requirements

## goals
The goal of the annotator project is to scan OpenEdge 4GL class file (file with a .cls extension) and scan these for so called annotations. All the annotation should be outputted in JSON format, so other tools can use these to do their jobs. The other tools are out of scope.

## technical stack
- We use the latest Go version, 1.25.5.
  - use cobra/viper for parameters
    - go get github.com/spf13/viper v1.21.0
    - go get -u github.com/spf13/cobra@latest
  - use zerolog for logging 
- We use `make` for building
- output is for windows (x64 only) and linux (x64 only)
- the sources are in `cmd/annotator`, the default structure for Go CLI tools
- main.go drives the command. For now there's just the `parse` command. Put the logic in `parse.go`. 
- data structures in `models.go`
- logging related in `logging.go` 
- split up code in logical files 
  
## function requirements
- the relevant 4gl syntax is in the `4gl-syntax.md` file
- annotator has a starting directory
- it recursively scans for .cls file
- each .cls file is parsed for:
  - the class name, fully qualified 
  - annotations
- there are 3 types of annotations:
  - class
  - method
  - free
- a annotation is tied to a method or a class, if there is between the annotation and the class or method there is:
  - zero lines
  - blank lines
  - comments
  - other annotations
- if not tied to a class or method, the type is free
- multiple annotations per construct can occur
- for every annotation the following is recorded:
  - the annotation name
  - an array of attribute/value pairs
  - the relative directory name of the containing .cls file
  - the fully qualified class name
  - the type
  - if applicable the method name, in `constructName`
  - the line number of the start of the annotation 
  - the line number of the class or method the annotation belongs to
- if an annotation spans multiple lines, the starting line is registered
- annotations in comments are not to be parsed
- there are nog particulair naming conventiions for attributes.
- if an annotation is not complete, the source will not compile. Since this tools is ran *after* compilation, correct annotations can be expected.
- escape character are the to be interpreted. Just register as-is.
- symlinks should be followed.
  
## output
The output should be an array of annotation name, which in turn has an array of all occurrences of that annotation.
The output should go to `annotations.json`, unless the `-o` parameter (with filename is specified). If is `--stdout` is specified, you know to do. `-o` takes preference over `--stdout`
The output should be pretty printed, unless the `--compact` is specified
The minimum output is in case of no annotations found is:
```
{ "annotations": {} }
```

## parameters
The command syntax should look like:
`annotator <command> <directory> [parameters]`
There's 1 command for now: `parse`. All described is is the parse command.
A `--help` or `-h` should display help.
A `--version` or `-v` should be supported. 
`<directory>` is obligatory, if not present fail and display help.

## logging
There's loglevels:
- none
- error
- info (default)
- debug
- trace
This is set by the `--loglevel` of `-l` parameter
We'll see later what goes where. You can make some obvious choices here straightaway.
Display status from info and up.
Logging should to `annotations.log`, unless `--logtoconsole` is specified.

## errors
If files or directories cannot be read, log the event and proceed.

## performance
project can range from a couple of dozen to thousand of line. Everything should be done in memory.

## exit code
Standard: 0 = success, 1 = general error, 2 = invalid usage

## executables
- windows: `annotator.exe` 
- linux: `annotator`

