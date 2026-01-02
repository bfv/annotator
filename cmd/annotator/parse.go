package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	outputFile   string
	stdout       bool
	compact      bool
	logLevel     string
	logToConsole bool
)

var parseCmd = &cobra.Command{
	Use:   "parse <directory>",
	Short: "Parse .cls files for annotations",
	Long:  `Recursively scan a directory for .cls files and extract annotations.`,
	Args:  cobra.ExactArgs(1),
	Run:   runParse,
}

func init() {
	parseCmd.Flags().StringVarP(&outputFile, "output", "o", "annotations.json", "Output file path")
	parseCmd.Flags().BoolVar(&stdout, "stdout", false, "Output to stdout")
	parseCmd.Flags().BoolVar(&compact, "compact", false, "Compact JSON output")
	parseCmd.Flags().StringVarP(&logLevel, "loglevel", "l", "info", "Log level (none, error, info, debug, trace)")
	parseCmd.Flags().BoolVar(&logToConsole, "logtoconsole", false, "Log to console instead of file")
}

func runParse(cmd *cobra.Command, args []string) {
	directory := args[0]
	startTime := time.Now()

	// Initialize logger
	logFile := "annotations.log"
	if err := InitLogger(LogLevel(logLevel), logToConsole, logFile); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.Info().Str("directory", directory).Msg("Starting parse")

	// Validate directory
	info, err := os.Stat(directory)
	if err != nil {
		logger.Error().Err(err).Str("directory", directory).Msg("Cannot access directory")
		fmt.Fprintf(os.Stderr, "Error: Cannot access directory %s: %v\n", directory, err)
		os.Exit(2)
	}
	if !info.IsDir() {
		logger.Error().Str("directory", directory).Msg("Path is not a directory")
		fmt.Fprintf(os.Stderr, "Error: %s is not a directory\n", directory)
		os.Exit(2)
	}

	// Scan for .cls files
	clsFiles, err := findClsFiles(directory)
	if err != nil {
		logger.Error().Err(err).Msg("Error scanning for .cls files")
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	logger.Info().Int("count", len(clsFiles)).Msg("Found .cls files")

	// Parse all files
	output := NewOutput()
	for _, file := range clsFiles {
		relPath, _ := filepath.Rel(directory, file)
		logger.Debug().Str("file", relPath).Msg("Parsing file")

		annotations, err := parseClsFile(file, directory)
		if err != nil {
			logger.Error().Err(err).Str("file", relPath).Msg("Error parsing file")
			continue
		}

		for _, ann := range annotations {
			output.AddAnnotation(ann)
		}
	}

	logger.Info().Int("total", countAnnotations(output)).Msg("Parsing complete")

	// Write output
	if err := writeOutput(output); err != nil {
		logger.Error().Err(err).Msg("Error writing output")
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Calculate and log elapsed time
	elapsed := time.Since(startTime)
	var elapsedStr string
	if elapsed.Seconds() < 5 {
		elapsedStr = fmt.Sprintf("%dms", elapsed.Milliseconds())
	} else {
		elapsedStr = fmt.Sprintf("%.1fs", elapsed.Seconds())
	}
	logger.Info().Str("elapsed", elapsedStr).Msg("Done")
}

func findClsFiles(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warn().Err(err).Str("path", path).Msg("Error accessing path")
			return nil // Continue despite errors
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".cls") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

func parseClsFile(filePath, baseDir string) ([]Annotation, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(baseDir, filePath)
	relPath = filepath.ToSlash(relPath)

	// Extract class name from file content
	className := extractClassName(string(content), relPath)

	// Parse the file for annotations
	annotations := extractAnnotations(string(content), relPath, className)

	return annotations, nil
}

func extractClassName(content, filePath string) string {
	// Remove comments first
	content = removeComments(content)

	// Look for CLASS statement
	re := regexp.MustCompile(`(?i)\bCLASS\s+([\w.]+)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}

	// Fallback: derive from file path
	className := strings.TrimSuffix(filePath, ".cls")
	className = strings.ReplaceAll(className, "\\", ".")
	className = strings.ReplaceAll(className, "/", ".")
	return className
}

func extractAnnotations(content, filePath, className string) []Annotation {
	var annotations []Annotation

	lines := strings.Split(content, "\n")

	// Track line types (annotation, comment, blank, class, method, code)
	lineTypes := classifyLines(lines)

	// Find all annotations
	processedLines := make(map[int]bool)
	for i := 0; i < len(lines); i++ {
		if lineTypes[i] == "annotation" && !processedLines[i] {
			ann, endLine := parseAnnotation(lines, i, filePath, className, lineTypes)
			if ann != nil {
				annotations = append(annotations, *ann)
				// Mark all lines of this annotation as processed
				for j := i; j <= endLine; j++ {
					processedLines[j] = true
				}
			}
		}
	}

	return annotations
}

func classifyLines(lines []string) []string {
	types := make([]string, len(lines))
	inBlockComment := false
	inAnnotation := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle block comments - check entire line, not just prefix
		if inBlockComment {
			types[i] = "comment"
			if strings.Contains(line, "*/") {
				inBlockComment = false
			}
			continue
		}

		// Check if block comment starts anywhere in the line
		if strings.Contains(line, "/*") {
			types[i] = "comment"
			// Check if it also closes on the same line
			commentStart := strings.Index(line, "/*")
			commentEnd := strings.Index(line[commentStart:], "*/")
			if commentEnd == -1 {
				// Comment doesn't close on this line
				inBlockComment = true
			}
			continue
		}

		// Handle line comments
		if strings.HasPrefix(trimmed, "//") {
			types[i] = "comment"
			continue
		}

		// Handle annotations (only if not in a comment)
		if strings.HasPrefix(trimmed, "@") {
			types[i] = "annotation"
			if !strings.Contains(trimmed, ".") {
				inAnnotation = true
			}
			continue
		}

		if inAnnotation {
			types[i] = "annotation"
			if strings.Contains(trimmed, ".") {
				inAnnotation = false
			}
			continue
		}

		// Handle blank lines
		if trimmed == "" {
			types[i] = "blank"
			continue
		}

		// Handle CLASS statement
		if regexp.MustCompile(`(?i)\bCLASS\s+`).MatchString(trimmed) {
			types[i] = "class"
			continue
		}

		// Handle METHOD statement
		if regexp.MustCompile(`(?i)\bMETHOD\s+`).MatchString(trimmed) {
			types[i] = "method"
			continue
		}

		// Handle DEFINE PROPERTY statement
		if regexp.MustCompile(`(?i)\bDEFINE\s+.*?\bPROPERTY\s+`).MatchString(trimmed) {
			types[i] = "property"
			continue
		}

		// Default
		types[i] = "code"
	}

	return types
}

func parseAnnotation(lines []string, startLine int, filePath, className string, lineTypes []string) (*Annotation, int) {
	// Collect full annotation text
	annotationText := ""
	endLine := startLine

	for i := startLine; i < len(lines); i++ {
		annotationText += lines[i]
		if strings.Contains(lines[i], ".") {
			endLine = i
			break
		}
	}

	// Parse annotation name and attributes
	name, attributes := parseAnnotationText(annotationText)
	if name == "" {
		return nil, endLine
	}

	// Find what this annotation is attached to
	constructType := "free"
	constructName := ""
	constructLine := 0

	// Look for class or method after the annotation
	for i := endLine + 1; i < len(lines); i++ {
		lt := lineTypes[i]

		if lt == "blank" || lt == "comment" || lt == "annotation" {
			continue
		}

		if lt == "class" {
			constructType = "class"
			constructLine = i + 1 // 1-based line number
			break
		}

		if lt == "method" {
			constructType = "method"
			constructName = extractMethodName(lines[i])
			constructLine = i + 1 // 1-based line number
			break
		}

		if lt == "property" {
			constructType = "property"
			constructName = extractPropertyName(lines[i])
			constructLine = i + 1 // 1-based line number
			break
		}

		// Hit other code, it's a free annotation
		break
	}

	ann := &Annotation{
		Name:           name,
		Attributes:     attributes,
		File:           filePath,
		ClassName:      className,
		Type:           constructType,
		ConstructName:  constructName,
		AnnotationLine: startLine + 1, // 1-based line number
	}

	if constructType != "free" {
		ann.ConstructLine = constructLine
	}

	return ann, endLine
}

func parseAnnotationText(text string) (string, []Attribute) {
	// Remove @ and final period
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "@")
	text = strings.TrimSuffix(text, ".")

	// Find annotation name and attributes
	parenIdx := strings.Index(text, "(")
	if parenIdx == -1 {
		// No attributes
		return strings.TrimSpace(text), []Attribute{}
	}

	name := strings.TrimSpace(text[:parenIdx])
	attrText := text[parenIdx+1:]

	// Remove closing paren
	if idx := strings.LastIndex(attrText, ")"); idx != -1 {
		attrText = attrText[:idx]
	}

	attributes := parseAttributes(attrText)

	return name, attributes
}

func parseAttributes(text string) []Attribute {
	var attributes []Attribute

	// Split by comma, but respect quoted strings
	parts := smartSplit(text, ',')

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split by = to get name and value
		eqIdx := strings.Index(part, "=")
		if eqIdx == -1 {
			continue
		}

		name := strings.TrimSpace(part[:eqIdx])
		value := strings.TrimSpace(part[eqIdx+1:])

		// Remove quotes from value
		value = strings.Trim(value, `"`)

		attributes = append(attributes, Attribute{
			Name:  name,
			Value: value,
		})
	}

	return attributes
}

func smartSplit(text string, sep rune) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, ch := range text {
		if ch == '"' {
			inQuotes = !inQuotes
			current.WriteRune(ch)
		} else if ch == sep && !inQuotes {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func extractMethodName(line string) string {
	// Look for method name after VOID or return type
	// Pattern: METHOD [modifiers] {VOID|return-type} method-name (
	re := regexp.MustCompile(`(?i)\bMETHOD\b.*?\b(?:VOID|[\w.]+)\s+([\w]+)\s*\(`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractPropertyName(line string) string {
	// Look for property name after PROPERTY keyword
	// Pattern: DEFINE [modifiers] PROPERTY property-name
	re := regexp.MustCompile(`(?i)\bDEFINE\b.*?\bPROPERTY\s+([\w]+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func removeComments(content string) string {
	// Remove block comments
	blockRe := regexp.MustCompile(`(?s)/\*.*?\*/`)
	content = blockRe.ReplaceAllString(content, "")

	// Remove line comments
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if idx := strings.Index(line, "//"); idx != -1 {
			lines[i] = line[:idx]
		}
	}

	return strings.Join(lines, "\n")
}

func countAnnotations(output *Output) int {
	count := 0
	for _, anns := range output.Annotations {
		count += len(anns)
	}
	return count
}

func writeOutput(output *Output) error {
	var data []byte
	var err error

	if compact {
		data, err = json.Marshal(output)
	} else {
		data, err = json.MarshalIndent(output, "", "  ")
	}

	if err != nil {
		return err
	}

	// Write to stdout if requested (-o takes precedence)
	if stdout && outputFile == "annotations.json" {
		fmt.Println(string(data))
		return nil
	}

	// Write to file
	return os.WriteFile(outputFile, data, 0644)
}
