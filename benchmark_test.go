package b_plus_tree_test

import (
	"github.com/SarthakMakhija/b-plus-tree/index"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	bytes := make([]byte, n)
	for iterator := range bytes {
		bytes[iterator] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(bytes)
}

var keyPrefix = randStringBytes(100)
var valuePrefix = randStringBytes(600)
var options = index.Options{
	FileName:                       "./index.db",
	PageSize:                       os.Getpagesize(),
	PreAllocatedPagePoolSize:       8192,
	AllowedPageOccupancyPercentage: 80,
}
var bPlusTree, _ = index.CreateBPlusTree(options)

func Benchmark_Put(b *testing.B) {
	benchmarkPutKeys(b, bPlusTree)
}

func Benchmark_GetFixedKey(b *testing.B) {
	if b.N == 1 {
		initializeBPlusTree := func(count int) {
			for iterator := 1; iterator <= count; iterator++ {
				key := "KEY" + "_" + keyPrefix + "_" + strconv.Itoa(iterator)
				err := bPlusTree.Put([]byte(key), []byte("VALUE"+"_"+valuePrefix))
				if err != nil {
					b.Fatalf("Failed while putting key/values during preppring for get %v", err)
				}
			}
		}
		initializeBPlusTree(100000)
	}
	searchKey := "KEY" + "_" + keyPrefix + "_" + strconv.Itoa(60000)
	benchmarkGetFixedKey(b, bPlusTree, searchKey)
}

func Benchmark_GetKeys(b *testing.B) {
	if b.N == 1 {
		initializeBPlusTree := func(count int) {
			for iterator := 1; iterator <= count; iterator++ {
				key := "KEY" + "_" + keyPrefix + "_" + strconv.Itoa(iterator)
				err := bPlusTree.Put([]byte(key), []byte("VALUE"+"_"+valuePrefix))
				if err != nil {
					b.Fatalf("Failed while putting key/values during preppring for get %v", err)
				}
			}
		}
		initializeBPlusTree(100000)
	}
	benchmarkGetKeys(b, bPlusTree)
}

func benchmarkPutKeys(b *testing.B, bPlusTree *index.BPlusTree) {
	b.ResetTimer()
	for iterator := 0; iterator < b.N; iterator++ {
		key := "KEY" + "_" + keyPrefix + "_" + strconv.Itoa(iterator)
		err := bPlusTree.Put([]byte(key), []byte("VALUE"+"_"+valuePrefix))
		if err != nil {
			b.Fatalf("Failed while putting %v", err)
		}
	}
}

func benchmarkGetFixedKey(b *testing.B, bPlusTree *index.BPlusTree, key string) {
	b.ResetTimer()
	for iterator := 0; iterator < b.N; iterator++ {
		getResult := bPlusTree.Get([]byte(key))
		if getResult.Err != nil {
			b.Fatalf("Failed while getting %v", getResult.Err)
		}
	}
}

func benchmarkGetKeys(b *testing.B, bPlusTree *index.BPlusTree) {
	b.ResetTimer()
	for iterator := 0; iterator < b.N; iterator++ {
		key := "KEY" + "_" + keyPrefix + "_" + strconv.Itoa(iterator)
		getResult := bPlusTree.Get([]byte(key))
		if getResult.Err != nil {
			b.Fatalf("Failed while getting %v", getResult.Err)
		}
	}
}
