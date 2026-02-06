package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractClassNameWithInterface(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		filePath string
		expected string
	}{
		{
			name: "Simple interface",
			content: `interface UserRepository:
				method public void save().
			end interface.`,
			filePath: "test/UserRepository.cls",
			expected: "UserRepository",
		},
		{
			name: "Interface with inheritance",
			content: `interface PaymentService inherits BaseService:
				method public void process().
			end interface.`,
			filePath: "test/PaymentService.cls",
			expected: "PaymentService",
		},
		{
			name: "Class and interface mixed - interface first",
			content: `interface IService:
				method public void doWork().
			end interface.
			
			class ServiceImpl implements IService:
			end class.`,
			filePath: "test/ServiceImpl.cls",
			expected: "IService", // Should pick first one found
		},
		{
			name: "Class and interface mixed - class first",
			content: `class ServiceImpl implements IService:
			end class.
			
			interface IService:
				method public void doWork().
			end interface.`,
			filePath: "test/ServiceImpl.cls",
			expected: "ServiceImpl", // Should pick first one found
		},
		{
			name:     "Fallback to filename",
			content:  `// No class or interface declaration`,
			filePath: "test/FallbackExample.cls",
			expected: "test.FallbackExample",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractClassName(tt.content, tt.filePath)
			if result != tt.expected {
				t.Errorf("extractClassName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClassifyLinesWithInterface(t *testing.T) {
	lines := []string{
		"@api(version=\"1.0\").",
		"interface UserService:",
		"",
		"  @http(method=\"get\").",
		"  method public User getUser().",
		"",
		"end interface.",
	}

	types := classifyLines(lines)

	expected := []string{
		"annotation",
		"interface",
		"blank",
		"annotation",
		"method",
		"blank",
		"code",
	}

	for i, expectedType := range expected {
		if i >= len(types) {
			t.Errorf("Missing line type at index %d", i)
			continue
		}
		if types[i] != expectedType {
			t.Errorf("Line %d: expected type %s, got %s", i, expectedType, types[i])
		}
	}
}

func TestParseInterfaceAnnotations(t *testing.T) {
	// Create a temporary test file
	tempDir, err := os.MkdirTemp("", "interface_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testContent := `@service(name="UserService", version="2.0").
@transactional.
interface UserService:

  @query(sql="SELECT * FROM users").
  @cache(ttl="300").
  method public User[] getAllUsers().

  @command(operation="insert").
  method public void createUser(input user as User).

  // Free annotation
  @todo(what="add validation").

end interface.`

	testFile := filepath.Join(tempDir, "UserService.cls")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the file
	annotations, err := parseClsFile(testFile, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Verify we have the expected annotations
	expectedAnnotations := []struct {
		name           string
		annotationType string
		constructName  string
		lineNum        int
	}{
		{"service", "interface", "", 1},
		{"transactional", "interface", "", 2},
		{"query", "method", "getAllUsers", 5},
		{"cache", "method", "getAllUsers", 6},
		{"command", "method", "createUser", 9},
		{"todo", "free", "", 12},
	}

	if len(annotations) != len(expectedAnnotations) {
		t.Fatalf("Expected %d annotations, got %d", len(expectedAnnotations), len(annotations))
	}

	for i, expected := range expectedAnnotations {
		ann := annotations[i]
		if ann.Name != expected.name {
			t.Errorf("Annotation %d: expected name %s, got %s", i, expected.name, ann.Name)
		}
		if ann.Type != expected.annotationType {
			t.Errorf("Annotation %d: expected type %s, got %s", i, expected.annotationType, ann.Type)
		}
		if ann.ConstructName != expected.constructName {
			t.Errorf("Annotation %d: expected construct name %s, got %s", i, expected.constructName, ann.ConstructName)
		}
		if ann.AnnotationLine != expected.lineNum {
			t.Errorf("Annotation %d: expected line %d, got %d", i, expected.lineNum, ann.AnnotationLine)
		}
	}

	// Verify specific annotation attributes
	serviceAnn := annotations[0]
	if len(serviceAnn.Attributes) != 2 {
		t.Errorf("Service annotation: expected 2 attributes, got %d", len(serviceAnn.Attributes))
	}

	// Check service attributes
	nameFound := false
	versionFound := false
	for _, attr := range serviceAnn.Attributes {
		if attr.Name == "name" && attr.Value == "UserService" {
			nameFound = true
		}
		if attr.Name == "version" && attr.Value == "2.0" {
			versionFound = true
		}
	}
	if !nameFound {
		t.Error("Service annotation missing 'name' attribute")
	}
	if !versionFound {
		t.Error("Service annotation missing 'version' attribute")
	}
}

func TestInterfaceMethodExtraction(t *testing.T) {
	testCases := []struct {
		name         string
		line         string
		expectedName string
	}{
		{
			name:         "Simple method",
			line:         "method public void doSomething():",
			expectedName: "doSomething",
		},
		{
			name:         "Method with return type",
			line:         "method public User getUserById(input id as integer):",
			expectedName: "getUserById",
		},
		{
			name:         "Method with modifiers",
			line:         "method protected static User findUser():",
			expectedName: "findUser",
		},
		{
			name:         "Interface method declaration (no implementation)",
			line:         "method public User[] getAllUsers().",
			expectedName: "getAllUsers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractMethodName(tc.line)
			if result != tc.expectedName {
				t.Errorf("extractMethodName(%q) = %q, want %q", tc.line, result, tc.expectedName)
			}
		})
	}
}

func TestMixedClassAndInterface(t *testing.T) {
	// Create temporary test file with mixed content
	tempDir, err := os.MkdirTemp("", "mixed_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testContent := `@contract.
interface IPaymentService:
  @http(method="POST").
  method public void processPayment().
end interface.

@service.
class PaymentServiceImpl implements IPaymentService:
  @override.
  method public void processPayment():
  end method.
end class.`

	testFile := filepath.Join(tempDir, "PaymentServiceImpl.cls")
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Parse the file
	annotations, err := parseClsFile(testFile, tempDir)
	if err != nil {
		t.Fatal(err)
	}

	// Should have 4 annotations: @contract (interface), @http (method), @service (class), @override (method)
	if len(annotations) != 4 {
		t.Fatalf("Expected 4 annotations, got %d", len(annotations))
	}

	// Check that we get the correct types
	expectedTypes := []string{"interface", "method", "class", "method"}
	for i, expectedType := range expectedTypes {
		if annotations[i].Type != expectedType {
			t.Errorf("Annotation %d: expected type %s, got %s", i, expectedType, annotations[i].Type)
		}
	}
}
