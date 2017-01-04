# WordWrap

[![Go Report Card](https://goreportcard.com/badge/github.com/donatj/wordwrap)](https://goreportcard.com/report/github.com/donatj/wordwrap)
[![GoDoc](https://godoc.org/github.com/donatj/wordwrap?status.svg)](https://godoc.org/github.com/donatj/wordwrap)
[![Build Status](https://travis-ci.org/donatj/wordwrap.svg?branch=master)](https://travis-ci.org/donatj/wordwrap)

UTF-8 Safe Word Wrapping for Go based on number of bytes.

This is a word wrap library that doesn’t *break UTF-8 runes* **and** *operates on number of bytes* rather than runes. It’s preference is to break on a unicode space character, but will break long words if necessary. This is particularly useful for sending text over protocols where message size is limited by bytes.

### Samples

English:

```go
fmt.Println(wordwrap.WrapString(`If any earl, baron, or other person that holds lands directly of the Crown, for military service, shall die, and at his death his heir shall be of full age and owe a 'relief', the heir shall have his inheritance on payment of the ancient scale of 'relief'.`, 60))
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
fmt.Println(wordwrap.WrapString(`クラウンの直接土地を保持している任意の伯爵、男爵、または他の人は、兵役のために、死ぬ、と彼の死で彼の後継者は成年であることと「救済」を借りなければならない場合は、相続人は、支払いの彼の継承をもたなければなりません「救済」の古代規模の。`, 60))
```

Becomes:

```
クラウンの直接土地を保持している任意の伯  // 60 bytes
爵、男爵、または他の人は、兵役のために、  // 60 bytes
死ぬ、と彼の死で彼の後継者は成年であるこ  // 60 bytes
とと「救済」を借りなければならない場合は  // 60 bytes
、相続人は、支払いの彼の継承をもたなけれ  // 60 bytes
ばなりません「救済」の古代規模の。        // 51 bytes
```

---

Korean:

```go
fmt.Println(wordwrap.WrapString(`크라운 의 직접 토지 를 보유하고 있는 백작 , 남작 , 또는 다른 사람이 군 복무 를 위해 죽을 것이요, 그의 죽음 에 그의 후계자 가 전체 연령 하고' 구호 '을 빚을 해야 하는 경우, 상속인 이 지불 에 대한 자신의 상속을 가져야한다 ' 구호 ' 의 고대 규모의 `, 60))
```

Becomes:

```go
크라운 의 직접 토지 를 보유하고 있는 백작  // 59 bytes
, 남작 , 또는 다른 사람이 군 복무 를 위해  // 57 bytes
죽을 것이요, 그의 죽음 에 그의 후계자 가   // 57 bytes
전체 연령 하고' 구호 '을 빚을 해야 하는    // 55 bytes
경우, 상속인 이 지불 에 대한 자신의 상속을 // 60 bytes
가져야한다 ' 구호 ' 의 고대 규모의         // 47 bytes
```
