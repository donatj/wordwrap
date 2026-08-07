# WordWrap

[![Go Report Card](https://goreportcard.com/badge/github.com/donatj/wordwrap)](https://goreportcard.com/report/github.com/donatj/wordwrap)
[![CI](https://github.com/donatj/wordwrap/actions/workflows/ci.yml/badge.svg)](https://github.com/donatj/wordwrap/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/donatj/wordwrap.svg)](https://pkg.go.dev/github.com/donatj/wordwrap)


UTF-8 Grapheme Cluster Safe Word Wrapping / Line Splitting for Go based on number of bytes.

This library wraps text without breaking UTF-8 grapheme clusters. It operates on byte count, not runes. It prefers breakable whitespace, then other Unicode whitespace, before breaking between grapheme clusters. It never splits emojis like 👩‍👩‍👧‍👧 or characters with combining marks.

This is useful for protocols where message size is limited by bytes.

## Requirements

Go 1.18 or later.

## Samples

English:

```go
package main

import (
	"fmt"
	"log"

	"github.com/donatj/wordwrap"
)

func main() {
	wrapped, err := wordwrap.WrapString(
		`If any earl, baron, or other person that holds lands directly of the Crown, for military service, shall die, and at his death his heir shall be of full age and owe a 'relief', the heir shall have his inheritance on payment of the ancient scale of 'relief'.`,
		60)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wrapped)
}
```

Becomes:

```
If any earl, baron, or other person that holds lands         // 53 bytes
directly of the Crown, for military service, shall die, and  // 60 bytes
at his death his heir shall be of full age and owe a         // 53 bytes
'relief', the heir shall have his inheritance on payment of  // 60 bytes
the ancient scale of 'relief'.                               // 30 bytes
```


---

Japanese:

```go
wrapped, err := wordwrap.WrapString(
	`クラウンの直接土地を保持している任意の伯爵、男爵、または他の人は、兵役のために、死ぬ、と彼の死で彼の後継者は成年であることと「救済」を借りなければならない場合は、相続人は、支払いの彼の継承をもたなければなりません「救済」の古代規模の。`,
	60)
if err != nil {
	log.Fatal(err)
}
fmt.Println(wrapped)
```

Becomes:

```
クラウンの直接土地を保持している任意の伯  // 60 bytes
爵、男爵、または他の人は、兵役のために、  // 60 bytes
死ぬ、と彼の死で彼の後継者は成年であるこ  // 60 bytes
とと「救済」を借りなければならない場合は  // 60 bytes
、相続人は、支払いの彼の継承をもたなけれ  // 60 bytes
ばなりません「救済」の古代規模の。       // 51 bytes
```

---

Korean:

```go
wrapped, err := wordwrap.WrapString(
	`크라운 의 직접 토지 를 보유하고 있는 백작 , 남작 , 또는 다른 사람이 군 복무 를 위해 죽을 것이요, 그의 죽음 에 그의 후계자 가 전체 연령 하고' 구호 '을 빚을 해야 하는 경우, 상속인 이 지불 에 대한 자신의 상속을 가져야한다 ' 구호 ' 의 고대 규모의 `,
	60)
if err != nil {
	log.Fatal(err)
}
fmt.Println(wrapped)
```

Becomes:

```
크라운 의 직접 토지 를 보유하고 있는 백작   // 59 bytes
, 남작 , 또는 다른 사람이 군 복무 를 위해  // 57 bytes
죽을 것이요, 그의 죽음 에 그의 후계자 가   // 57 bytes
전체 연령 하고' 구호 '을 빚을 해야 하는    // 55 bytes
경우, 상속인 이 지불 에 대한 자신의 상속을 // 60 bytes
가져야한다 ' 구호 ' 의 고대 규모의       // 47 bytes
```

---

Grapheme Clusters:

```go
wrapped, err := wordwrap.WrapString(
	`Hello 👩‍👩‍👧‍👧 family 🧑‍🎄 celebrating café with naïve résumé क्षि`,
	30)
if err != nil {
	log.Fatal(err)
}
fmt.Println(wrapped)
```

Becomes:

```
Hello                                 // 6 bytes
👩‍👩‍👧‍👧                                    // 26 bytes
family 🧑‍🎄                             // 19 bytes
celebrating café with naïve           // 30 bytes
résumé क्षि                             // 21 bytes
```

### Error Handling

The library returns an `ErrGraphemeClusterTooLarge` error when a grapheme cluster exceeds the byte limit. This allows you to handle these situations gracefully.

Single Japanese character on 2-byte limit:

```go
// import "errors"
_, err := wordwrap.SplitString("し", 2)  // error: し is 3 bytes
if errors.Is(err, wordwrap.ErrGraphemeClusterTooLarge) {
	// Handle the error
}
```

Family emoji on 20-byte limit:

```go
// import "errors"
_, err := wordwrap.SplitString("👩‍👩‍👧‍👧", 20)  // error: emoji is 25 bytes
if errors.Is(err, wordwrap.ErrGraphemeClusterTooLarge) {
	// Handle the error
}
```

Person with tree on 8-byte limit:

```go
// import "errors"
_, err := wordwrap.SplitString("🧑‍🎄", 8)  // error: emoji is 11 bytes
if errors.Is(err, wordwrap.ErrGraphemeClusterTooLarge) {
	// Handle the error
}
```

String ending with oversized cluster:

```go
// import "errors"
_, err := wordwrap.SplitString("test 👩‍👩‍👧‍👧", 20)  // error: emoji is 25 bytes
if errors.Is(err, wordwrap.ErrGraphemeClusterTooLarge) {
	// Handle the error
}
```
