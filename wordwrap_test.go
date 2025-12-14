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
