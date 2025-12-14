package wordwrap

import (
	"errors"
	"reflect"
	"testing"
)

func TestSplitString(t *testing.T) {
	tests := []struct {
		input   string
		output  []string
		bytelim uint
	}{
		{"asdasd asd asdasd",
			[]string{"asda", "sd ", "asd ", "asda", "sd"}, 4},

		{"𠜎𠜱00𠝹𠱓𠱸𠲖𠳏𠳕",
			[]string{"𠜎𠜱0", "0𠝹𠱓", "𠱸𠲖", "𠳏𠳕"}, 9},

		{`If any earl, baron, or other person that holds lands directly of the Crown, for military service, shall die, and at his death his heir shall be of full age and owe a 'relief', the heir shall have his inheritance on payment of the ancient scale of 'relief'.`,
			[]string{
				"If any earl, baron, or other person that holds lands ",
				"directly of the Crown, for military service, shall die, and ",
				"at his death his heir shall be of full age and owe a ",
				"'relief', the heir shall have his inheritance on payment of ",
				"the ancient scale of 'relief'."}, 60},

		{`クラウンの直接土地を保持している任意の伯爵、男爵、または他の人は、兵役のために、死ぬ、と彼の死で彼の後継者は成年であることと「救済」を借りなければならない場合は、相続人は、支払いの彼の継承をもたなければなりません「救済」の古代規模の。`,
			[]string{
				"クラウンの直接土地を保持している任意の伯",
				"爵、男爵、または他の人は、兵役のために、",
				"死ぬ、と彼の死で彼の後継者は成年であるこ",
				"とと「救済」を借りなければならない場合は",
				"、相続人は、支払いの彼の継承をもたなけれ",
				"ばなりません「救済」の古代規模の。"}, 60},

		{`크라운 의 직접 토지 를 보유하고 있는 백작 , 남작 , 또는 다른 사람이 군 복무 를 위해 죽을 것이요, 그의 죽음 에 그의 후계자 가 전체 연령 하고' 구호 '을 빚을 해야 하는 경우, 상속인 이 지불 에 대한 자신의 상속을 가져야한다 ' 구호 ' 의 고대 규모의 `,
			[]string{
				"크라운 의 직접 토지 를 보유하고 있는 백작 ",
				", 남작 , 또는 다른 사람이 군 복무 를 위해 ",
				"죽을 것이요, 그의 죽음 에 그의 후계자 가 ",
				"전체 연령 하고' 구호 '을 빚을 해야 하는 ",
				"경우, 상속인 이 지불 에 대한 자신의 상속을 ",
				"가져야한다 ' 구호 ' 의 고대 규모의 "}, 60},

		// ZWJ sequences - family emoji
		{"Hello 👩‍👩‍👧‍👧 world",
			[]string{"Hello 👩‍👩‍👧‍👧 ", "world"}, 32},

		// ZWJ sequences - person with Christmas tree
		{"Test 🧑‍🎄 emoji here",
			[]string{"Test 🧑‍🎄 ", "emoji here"}, 20},

		// Long word with ZWJ emoji (no spaces to break on)
		{"abcdefgh👩‍👩‍👧‍👧ijklmn",
			[]string{"abcdefgh", "👩‍👩‍👧‍👧ijklm", "n"}, 30},

		// Multiple ZWJ emojis
		{"🧑‍🎄 and 👩‍👩‍👧‍👧 test",
			[]string{"🧑‍🎄 and ", "👩‍👩‍👧‍👧 ", "test"}, 30},

		// ZWJ emoji at the start
		{"👩‍👩‍👧‍👧 family",
			[]string{"👩‍👩‍👧‍👧 ", "family"}, 30},

		// ZWJ emoji at the end
		{"family 👩‍👩‍👧‍👧",
			[]string{"family ", "👩‍👩‍👧‍👧"}, 30},

		// Devanagari complex clusters
		{"नमस्ते क्षि test",
			[]string{"नमस्ते ", "क्षि test"}, 20},

		// Devanagari multiple clusters
		{"श्री त्र द्ध test",
			[]string{"श्री ", "त्र द्ध ", "test"}, 20},

		// Arabic with diacritics
		{"السلام عليكم مُحَمَّد test",
			[]string{"السلام عليكم ", "مُحَمَّد test"}, 25},

		// Hebrew with points
		{"שָׁלוֹם test word",
			[]string{"שָׁלוֹם test ", "word"}, 20},

		// Thai with tone marks
		{"สวัสดี ก้า test",
			[]string{"สวัสดี ", "ก้า test"}, 20},

		// Emoji with skin tone modifiers
		{"Hello 👋🏽 👍🏿 world",
			[]string{"Hello 👋🏽 ", "👍🏿 world"}, 20},

		// Emoji woman technologist (ZWJ with profession)
		{"Test 👩‍💻 code",
			[]string{"Test 👩‍💻 ", "code"}, 20},

		// Keycap sequences
		{"Numbers 1️⃣ 2️⃣ 3️⃣ here",
			[]string{"Numbers 1️⃣ ", "2️⃣ 3️⃣ ", "here"}, 20},

		// Regional indicator (flag emoji) - fits within limit
		{"Hello 🇺🇸 test",
			[]string{"Hello 🇺🇸 test"}, 20},

		// Bengali complex cluster
		{"বাংলা ক্ষ test",
			[]string{"বাংলা ", "ক্ষ test"}, 20},

		// Tamil with vowel signs
		{"தமிழ் நீ கூ test",
			[]string{"தமிழ் ", "நீ கூ test"}, 20},

		// Vietnamese with multiple combining marks
		{"Tiếng Việt ệ test",
			[]string{"Tiếng Việt ệ ", "test"}, 20},
	}

	for _, test := range tests {
		actual, err := SplitString(test.input, test.bytelim)
		if err != nil {
			t.Errorf(`SplitString(%#v) returned unexpected error: %v`, test.input, err)
			continue
		}

		if !reflect.DeepEqual(actual, test.output) {
			t.Errorf(`SplitString(%#v) = %#v; want %#v`, test.input, actual, test.output)
		}
	}
}

func TestSplitStringError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		bytelim uint
	}{
		{
			name:    "Family emoji too large",
			input:   "👩‍👩‍👧‍👧",
			bytelim: 20, // Family emoji is 25 bytes
		},
		{
			name:    "Person with tree emoji too large",
			input:   "🧑‍🎄",
			bytelim: 8, // Person with tree is 11 bytes
		},
		{
			name:    "Grapheme cluster in word too large",
			input:   "test👩‍👩‍👧‍👧end",
			bytelim: 20, // Cannot break within the emoji
		},
		{
			name:    "Devanagari single cluster too large",
			input:   "क्",
			bytelim: 5, // क् is 6 bytes
		},
		{
			name:    "Devanagari cluster at end too large",
			input:   "test नी",
			bytelim: 5, // "test " is 5 bytes, नी is 6 bytes, needs > 11 total, but नी alone exceeds 5
		},
		{
			name:    "Thai cluster single too large",
			input:   "ก้",
			bytelim: 5, // ก้ is 6 bytes
		},
		{
			name:    "Tag sequence flag too large",
			input:   "🏴󠁧󠁢󠁥󠁮󠁧󠁿",
			bytelim: 25, // England flag is 28 bytes
		},
		{
			name:    "Emoji with skin tone at end",
			input:   "test 👋🏽",
			bytelim: 7, // 👋🏽 is 8 bytes, "test " is 5 bytes, total 13, cannot fit at limit 7
		},
		{
			name:    "Keycap sequence too large",
			input:   "1️⃣",
			bytelim: 6, // 1️⃣ is 7 bytes
		},
		{
			name:    "Vietnamese combining marks too large",
			input:   "ệ",
			bytelim: 2, // ệ is 3 bytes
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := SplitString(test.input, test.bytelim)
			if err == nil {
				t.Errorf("SplitString(%#v, %d) should have returned an error", test.input, test.bytelim)
			}
			if !errors.Is(err, ErrGraphemeClusterTooLarge) {
				t.Errorf("SplitString(%#v, %d) returned wrong error: got %v, want %v", test.input, test.bytelim, err, ErrGraphemeClusterTooLarge)
			}
		})
	}
}

func TestWrapString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		bytelim  uint
		expected string
	}{
		{
			name:     "Simple wrapping",
			input:    "Hello world this is a test",
			bytelim:  10,
			expected: "Hello \nworld \nthis is a \ntest",
		},
		{
			name:     "English text with spaces",
			input:    "If any earl, baron, or other person that holds lands directly of the Crown",
			bytelim:  30,
			expected: "If any earl, baron, or other \nperson that holds lands \ndirectly of the Crown",
		},
		{
			name:     "Unicode Japanese text",
			input:    "クラウンの直接土地を保持している任意の伯爵、男爵",
			bytelim:  30,
			expected: "クラウンの直接土地を\n保持している任意の伯\n爵、男爵",
		},
		{
			name:     "Text with emoji",
			input:    "Hello 👋🏽 world",
			bytelim:  15,
			expected: "Hello 👋🏽 \nworld",
		},
		{
			name:     "Single line that fits",
			input:    "Short",
			bytelim:  20,
			expected: "Short",
		},
		{
			name:     "Multiple ZWJ emojis",
			input:    "🧑‍🎄 and 👩‍👩‍👧‍👧 test",
			bytelim:  30,
			expected: "🧑‍🎄 and \n👩‍👩‍👧‍👧 \ntest",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := WrapString(test.input, test.bytelim)
			if err != nil {
				t.Errorf("WrapString(%#v, %d) returned unexpected error: %v", test.input, test.bytelim, err)
				return
			}

			if actual != test.expected {
				t.Errorf("WrapString(%#v, %d) = %#v; want %#v", test.input, test.bytelim, actual, test.expected)
			}
		})
	}
}

func TestWrapStringError(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		bytelim uint
	}{
		{
			name:    "Family emoji too large",
			input:   "👩‍👩‍👧‍👧",
			bytelim: 20, // Family emoji is 25 bytes
		},
		{
			name:    "Person with tree emoji too large",
			input:   "🧑‍🎄",
			bytelim: 8, // Person with tree is 11 bytes
		},
		{
			name:    "Grapheme cluster in text too large",
			input:   "test👩‍👩‍👧‍👧end",
			bytelim: 20, // Cannot break within the emoji
		},
		{
			name:    "Single character too large",
			input:   "し",
			bytelim: 2, // し is 3 bytes
		},
		{
			name:    "Thai cluster too large",
			input:   "ก้",
			bytelim: 5, // ก้ is 6 bytes
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := WrapString(test.input, test.bytelim)
			if err == nil {
				t.Errorf("WrapString(%#v, %d) should have returned an error", test.input, test.bytelim)
			}
			if !errors.Is(err, ErrGraphemeClusterTooLarge) {
				t.Errorf("WrapString(%#v, %d) returned wrong error: got %v, want %v", test.input, test.bytelim, err, ErrGraphemeClusterTooLarge)
			}
		})
	}
}

func TestSplitBuilder_DefaultBehavior(t *testing.T) {
	// Test that default SplitBuilder matches SplitString behavior
	input := "asdasd asd asdasd"
	bytelim := uint(4)
	expected := []string{"asda", "sd ", "asd ", "asda", "sd"}

	sb := NewSplitBuilder()
	
	var actual []string
	for _, line := range sb.Split(input, bytelim) {
		actual = append(actual, line)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("SplitBuilder.Split(%#v) = %#v; want %#v", input, actual, expected)
	}
}

func TestSplitBuilder_WithIndex(t *testing.T) {
	input := "Hello world this is a test"
	bytelim := uint(10)
	
	sb := NewSplitBuilder()
	
	expectedLines := []string{"Hello ", "world ", "this is a ", "test"}
	actualLines := []string{}
	actualIndices := []int{}
	
	for idx, line := range sb.Split(input, bytelim) {
		actualIndices = append(actualIndices, idx)
		actualLines = append(actualLines, line)
	}
	
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Errorf("Lines mismatch: got %#v; want %#v", actualLines, expectedLines)
	}
	
	expectedIndices := []int{0, 1, 2, 3}
	if !reflect.DeepEqual(actualIndices, expectedIndices) {
		t.Errorf("Indices mismatch: got %#v; want %#v", actualIndices, expectedIndices)
	}
}

func TestSplitBuilder_TrimTrailingWhiteSpace(t *testing.T) {
	input := "Hello world this is a test"
	bytelim := uint(10)
	
	sb := NewSplitBuilder(TrimTrailingWhiteSpace(true))
	
	expectedLines := []string{"Hello", "world", "this is a", "test"}
	actualLines := []string{}
	
	for _, line := range sb.Split(input, bytelim) {
		actualLines = append(actualLines, line)
	}
	
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Errorf("Lines with trim: got %#v; want %#v", actualLines, expectedLines)
	}
}

func TestSplitBuilder_TrimTrailingWhiteSpace_MultipleSpaces(t *testing.T) {
	input := "test   more   data"
	bytelim := uint(10)
	
	sb := NewSplitBuilder(TrimTrailingWhiteSpace(true))
	
	expectedLines := []string{"test", "more", "data"}
	actualLines := []string{}
	
	for _, line := range sb.Split(input, bytelim) {
		actualLines = append(actualLines, line)
	}
	
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Errorf("Lines with multiple spaces trim: got %#v; want %#v", actualLines, expectedLines)
	}
}

func TestSplitBuilder_ContinueOnError(t *testing.T) {
	// Test with a grapheme cluster that's too large
	input := "test 👩‍👩‍👧‍👧 end"
	bytelim := uint(10) // Family emoji is 25 bytes, which exceeds limit
	
	sb := NewSplitBuilder(ContinueOnError(true))
	
	var lines []string
	for _, line := range sb.Split(input, bytelim) {
		lines = append(lines, line)
	}
	
	// With continueOnError, we should get some output
	if len(lines) == 0 {
		t.Errorf("Expected some lines with ContinueOnError, got none")
	}
}

func TestSplitBuilder_BreakGraphemeClusters(t *testing.T) {
	// Test breaking within a grapheme cluster
	input := "test 👩‍👩‍👧‍👧 end"
	bytelim := uint(10)
	
	sb := NewSplitBuilder(BreakGraphemeClusters(true))
	
	var lines []string
	for _, line := range sb.Split(input, bytelim) {
		lines = append(lines, line)
	}
	
	// With breakGraphemeClusters, we should get multiple lines
	if len(lines) < 2 {
		t.Errorf("Expected multiple lines with BreakGraphemeClusters, got %d", len(lines))
	}
}

func TestSplitBuilder_CombinedOptions(t *testing.T) {
	input := "Hello world  test"
	bytelim := uint(10)
	
	sb := NewSplitBuilder(
		TrimTrailingWhiteSpace(true),
		ContinueOnError(true),
	)
	
	expectedLines := []string{"Hello", "world", "test"}
	actualLines := []string{}
	
	for _, line := range sb.Split(input, bytelim) {
		actualLines = append(actualLines, line)
	}
	
	if !reflect.DeepEqual(actualLines, expectedLines) {
		t.Errorf("Combined options: got %#v; want %#v", actualLines, expectedLines)
	}
}

func TestSplitBuilder_EmptyString(t *testing.T) {
	input := ""
	bytelim := uint(10)
	
	sb := NewSplitBuilder()
	
	var lines []string
	for _, line := range sb.Split(input, bytelim) {
		lines = append(lines, line)
	}
	
	if len(lines) != 0 {
		t.Errorf("Expected no lines for empty string, got %d", len(lines))
	}
}

func TestSplitBuilder_Unicode(t *testing.T) {
	input := "クラウンの直接土地を保持している任意の伯爵、男爵"
	bytelim := uint(30)
	
	sb := NewSplitBuilder()
	
	var lines []string
	for _, line := range sb.Split(input, bytelim) {
		lines = append(lines, line)
	}
	
	// Verify we got multiple lines and each is within byte limit
	if len(lines) < 2 {
		t.Errorf("Expected multiple lines for long Unicode text, got %d", len(lines))
	}
	
	for i, line := range lines {
		if len(line) > int(bytelim) {
			t.Errorf("Line %d exceeds byte limit: %d > %d", i, len(line), bytelim)
		}
	}
}
