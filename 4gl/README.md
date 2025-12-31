# 4GL Test Cases

This directory contains OpenEdge 4GL test cases for the annotator tool.

## Test Files

- **HelloWorld.cls** - Simple class with basic method annotations
- **UserService.cls** - Service class with HTTP annotations, free annotations, and multi-line annotations
- **misc/string/StringHelper.cls** - Demonstrates nested directory structure and static methods
- **ComplexExample.cls** - Complex scenarios with comments, blank lines, and multiple annotations
- **NoAnnotations.cls** - Class without any annotations (should produce empty output)
- **EdgeCases.cls** - Edge cases including escaped characters, empty values, annotations in comments
- **api/controller/ProductController.cls** - REST controller pattern with multiple HTTP method annotations

## Expected Behavior

Run the annotator on this directory:
```
annotator.exe parse 4gl
```

The tool should:
1. Find all .cls files recursively
2. Parse annotations from each file
3. Generate annotations.json with all found annotations
4. Ignore annotations in comments
5. Handle multi-line annotations
6. Correctly identify class vs method vs free annotations
