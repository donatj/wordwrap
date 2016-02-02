package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

func WrapString(s string, lim uint) string {
	output := []string{}
	l := len(s)

	bc := []byte{}
	brk := [2]int{0, 0}
	lastValid := [2]int{0, 0}

	var lineLen, start int
	for i := 0; i <= l-1; i++ {

		bc = append(bc, s[i])

		_ = start
		if lineLen == int(lim) {
			if brk[1] == 0 {
				if lastValid[1] == 0 {
					panic(fmt.Sprint(lastValid, i, start, lineLen, bc))
				}

				fmt.Println("i:", i, "line:", lineLen, "lim:", lim, lastValid)
				ss := s[start:lastValid[0]]
				fmt.Printf("long-val: '%s' - %d\n", ss, len(ss))
				i = lastValid[0] + lastValid[1] - 1
				start = i

				lastValid = [2]int{}
			} else {
				//				fmt.Println(start, brk[0], i)

				i = brk[0] + brk[1] - 1
				ss := s[start:brk[0]]
				fmt.Printf("spac-val: '%s' - %d\n", ss, len(ss))
				start = i

				//				fmt.Println(start, i)

				brk = [2]int{0, 0}
				//				return ""
			}

			bc = []byte{}
			fmt.Println("bc-c", i)
			lineLen = 0
		}

		if len(bc) > 0 && utf8.Valid(bc) {

			r, size := utf8.DecodeRune(bc)
			//			fmt.Println(size, len(bc))
			fmt.Println("r:", string(r), i, size)
			lastValid = [2]int{i - size + 1, size}

			if unicode.IsSpace(r) {
				brk = lastValid
			}

			bc = []byte{}
		} else if len(bc) > 4 {
			panic(fmt.Sprint("too long ", i, " ", string(bc)))
		}

		_ = brk
		lineLen++
		//		fmt.Print(brk)

		time.Sleep(10 * time.Millisecond)
	}

	fmt.Print("\n")
	return strings.Join(output, "\n")
}

func main() {
	//	fmt.Println(WrapString("This istoo longandwontbreak Lets This is a test test test test boop thisisalsotoodamnlongandwnotbreak", 10))
	//	fmt.Println(WrapString("DecodeRq	une\xe2\x80\x83unpacks the\xe2\x80\x83first UTF-8 encoding\xe2\x80\x83in p and returns the rune and its width in bytes. If p is empty it returns (RuneError, 0). Otherwise, if the encoding is invalid, it returns (RuneError, 1). Both are impossible results for correct UTF-8.", 40))
	fmt.Println(WrapString(`1890 年代に アメbリカ 合衆国ニ ュ ージャージー州で、行政単位としての小さな「ボロ」が大量に作られたことを 表す言葉である。ニュージャージー州議会が地方政府と教育体系を改革しようとしたことが、当時あったタウンシップを小さなボロに分割することに繋がった ウンゼンド子爵位を前身とし、1787年に陸まる爵位である……`, 20))

}
