package main

// Attribute represents a single attribute-value pair in an annotation
type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Annotation represents a parsed annotation from a 4GL class file
type Annotation struct {
	Name           string      `json:"name"`
	Attributes     []Attribute `json:"attributes"`
	File           string      `json:"file"`
	ClassName      string      `json:"classname"`
	Type           string      `json:"type"` // "class", "method", or "free"
	ConstructName  string      `json:"constructName,omitempty"`
	AnnotationLine int         `json:"annotationLine"`
	ConstructLine  int         `json:"constructLine,omitempty"`
}

// Output represents the JSON output structure
type Output struct {
	Annotations map[string][]Annotation `json:"annotations"`
}

// NewOutput creates a new Output instance with initialized map
func NewOutput() *Output {
	return &Output{
		Annotations: make(map[string][]Annotation),
	}
}

// AddAnnotation adds an annotation to the output
func (o *Output) AddAnnotation(ann Annotation) {
	o.Annotations[ann.Name] = append(o.Annotations[ann.Name], ann)
}
