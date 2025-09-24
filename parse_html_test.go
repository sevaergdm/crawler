package main

import (
	"strings"
	"testing"
)

func testGetH1FromHTMLBasic(t *testing.T) {
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	expected := "Test Title"
	actual, err := getH1FromHTML(inputBody)
	if err != nil {
		t.Errorf("expected no error, but got %v", err)
	}

	if expected != actual {
		t.Errorf("expected '%q', but got '%q'", expected, actual)
	}
}

func TestGetFirstParagraphFromHTMLMainPriority(t *testing.T) {
	testCases := []struct{
		name string
		inputBody string
		expected string
		errorContains string
	}{
		{
			name: "One paragraph in main",
			inputBody: `<html><body>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "Two paragraphs, one in main, one outside",
			inputBody: `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`,
			expected: "Main paragraph.",
		},
		{
			name: "One paragraph without main",
			inputBody: `<html><body>
		<p>Outside paragraph.</p>
	</body></html>`,
			expected: "Outside paragraph.",
		},
		{
			name: "No paragraphs",
			inputBody: `<html><body>
	</body></html>`,
			expected: "",
		},
	}
	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getFirstParagraphFromHTML(tc.inputBody)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if actual != tc.expected {
				t.Errorf("Test %v - '%s' FAIL: expected '%s' but got '%s'", i, tc.name, tc.expected, actual)
			}
		})
	}
}
