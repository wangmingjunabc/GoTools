package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const SYSTEM_BIT int64 = 32 << (^uint(0) >> 63)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int64) bool {
	word, bit := x/SYSTEM_BIT, uint(x%SYSTEM_BIT)
	return int(word) < len(s.words) && s.words[word]&(1<<bit) != 0
}

func (s *IntSet) Add(x int64) {
	word, bit := x/SYSTEM_BIT, uint(x%SYSTEM_BIT)
	for int(word) >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

func (s *IntSet) Remove(x int64) {
	word, bit := x/SYSTEM_BIT, uint(x%SYSTEM_BIT)
	s.words[word] &^= 1 << bit
}

func (s *IntSet) Clear() {
	for i := range s.words {
		s.words[i] = 0
	}
}

func (s *IntSet) Copy() *IntSet {
	ss := &IntSet{}
	ss.words = make([]uint64, len(s.words))
	copy(ss.words, s.words)
	return ss
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) Len() int {
	count := 0
	for _, word := range s.words {
		for word != 0 {
			count++
			word &= word - 1
		}
	}
	return count
}

//字符串表现形式

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

//以上是实现位向量bitvector

type BloomFilter struct {
	BitVector *IntSet
	Seeds     []int64
}

func (bf *BloomFilter) Hash(value string, seed int64) int64 {
	var ret int64
	for _, ch := range value {
		ret += seed*ret + int64(ch)
	}
	return ret & ((1 << 31) - 1)
}

func (bf *BloomFilter) Insert(value string) {
	for _, seed := range bf.Seeds {
		ret := bf.Hash(value, seed)
		bf.BitVector.Add(ret)
	}
}

func (bf *BloomFilter) IsContain(value string) bool {
	if value == "" {
		return false
	}
	result := true
	for _, seed := range bf.Seeds {
		ret := bf.Hash(value, seed)
		result = bf.BitVector.Has(ret)
	}
	return result
}

func main() {
	bf := BloomFilter{}
	bf.BitVector = new(IntSet)
	bf.Seeds = []int64{5, 7, 11, 13, 19, 31, 37, 61}
	bf.BitVector.Add(1 << 32)

	f, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	buf := bufio.NewReader(f)
	for {
		b, errR := buf.ReadBytes('\n')
		if errR != nil {
			if errR == io.EOF {
				break
			}
			fmt.Println(errR.Error())
		}

		s := strings.TrimRight(string(b), "\r\n")
		if strings.Compare(s, "exit") == 0 {
			fmt.Println("complete and exit now")
			break
		} else if bf.IsContain(s) == false {
			bf.Insert(s)
		} else {
			fmt.Printf("url: %s is exist\n", s)
		}
	}

}
