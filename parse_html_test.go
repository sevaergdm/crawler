package main

import (
	"net/url"
	"reflect"
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

func TestGetURLsFromHTML(t *testing.T) {
	testCases := []struct{
		name string
		inputBody string
		inputURL string
		expected []string
	}{
		{
			name: "Absolute reference",
			inputBody: `<html><body><a href="https://blog.boot.dev"><span>Boot.dev</span></a></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name: "nested reference",
			inputBody: `<html><body><a href="/posts/myfancynewpost"><span>Boot.dev</span></a></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://blog.boot.dev/posts/myfancynewpost"},
		},
		{
			name: "external reference",
			inputBody: `<html><body><a href="https://www.google.com"><span>Boot.dev</span></a></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://www.google.com"},
		},
		{
			name: "multi reference",
			inputBody: `<html><body><a href="/posts/thisisonepost"><a href="/posts/thisisanotherpost"><span>Boot.dev</span></a></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://blog.boot.dev/posts/thisisonepost", "https://blog.boot.dev/posts/thisisanotherpost"},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}
			
			actual, err := getURLsFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected '%s' but got '%s'", i, tc.name, tc.expected, actual)
			}
		})

	}
}

func TestGetImagessFromHTML(t *testing.T) {
	testCases := []struct{
		name string
		inputBody string
		inputURL string
		expected []string
	}{
		{
			name: "Relative reference",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://blog.boot.dev/logo.png"},
		},
		{
			name: "Missing reference",
			inputBody: `<html><body></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{},
		},
		{
			name: "Multiple reference",
			inputBody: `<html><body><img src="/logo.png" alt="Logo"><img src="/logo2.png" alt="Logo2"></body></html>`,
			inputURL: "https://blog.boot.dev",
			expected: []string{"https://blog.boot.dev/logo.png", "https://blog.boot.dev/logo2.png"},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("couldn't parse input URL: %v", err)
				return
			}
			
			actual, err := getImagesFromHTML(tc.inputBody, baseURL)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected '%s' but got '%s'", i, tc.name, tc.expected, actual)
			}
		})

	}
}

func TestExtractPageData(t *testing.T) {
	testCases := []struct{
		name string
		inputBody string
		inputURL string
		expected PageData
	}{
		{
			name: "Basic",
			inputBody: `<html><body>
        <h1>Test Title</h1>
        <p>This is the first paragraph.</p>
        <a href="/link1">Link 1</a>
        <img src="/image1.jpg" alt="Image 1">
    </body></html>`, 
			inputURL: "https://blog.boot.dev",
			expected: PageData{
				URL: "https://blog.boot.dev",
				H1: "Test Title",
				FirstParagraph: "This is the first paragraph.",
				OutgoingLinks: []string{"https://blog.boot.dev/link1"},
				ImageURLs: []string{"https://blog.boot.dev/image1.jpg"},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual:= extractPageData(tc.inputBody, tc.inputURL)

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - '%s' FAIL: expected '%+v' but got '%+v'", i, tc.name, tc.expected, actual)
			}
		})

	}
}
