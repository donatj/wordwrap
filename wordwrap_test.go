package wordwrap

import (
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
	}

	for _, test := range tests {
		actual := SplitString(test.input, test.bytelim)

		if !reflect.DeepEqual(actual, test.output) {
			t.Errorf(`SplitString(%#v) = %#v; want %#v`, test.input, actual, test.output)
		}
	}
}
