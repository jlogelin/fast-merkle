package main

import (
	//"crypto/sha256"

	"sort"

	"github.com/minio/sha256-simd"
)

func leafHash(n []byte) []byte {
	h := sha256.Sum256(n)
	return h[:]
}

func parentHash(l, r []byte) []byte {
	h := sha256.New()
	h.Write(l)
	h.Write(r)
	return h.Sum(nil)
}

func foldr(f func([]byte, []byte) []byte, coll [][]byte) []byte {
	if len(coll) == 0 {
		return nil
	}
	res := coll[len(coll)-1]
	for i := len(coll) - 2; i >= 0; i-- {
		res = f(coll[i], res)
	}
	return res
}

func insert(s map[int][]byte, v []byte, n int) map[int][]byte {
	if _, ok := s[n]; ok {
		p := parentHash(s[n], v)
		return insert(del(s, n), p, n+1)
	}
	s[n] = v
	return s
}

func del(s map[int][]byte, n int) map[int][]byte {
	delete(s, n)
	return s
}

func finalize(s map[int][]byte) []byte {
	var keys []int
	for k := range s {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))
	var vals [][]byte
	for _, k := range keys {
		vals = append(vals, s[k])
	}
	return foldr(parentHash, vals)
}

func MerkleRoot(stream [][]byte) []byte {
	if len(stream) == 0 {
		return nil
	}
	m := make(map[int][]byte)
	for _, v := range stream {
		m = insert(m, leafHash(v), 0)
	}
	return finalize(m)
}