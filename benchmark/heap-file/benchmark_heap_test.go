package heap_file_test

import (
	heap_file "github.com/SarthakMakhija/heap-file/heap-file"
	"github.com/SarthakMakhija/heap-file/heap-file/field"
	"github.com/SarthakMakhija/heap-file/heap-file/tuple"
	"github.com/SarthakMakhija/heap-file/index"
	"math/rand"
	"os"
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

var stringField1Value = randStringBytes(500)
var stringField2Value = randStringBytes(500)
var uint32FieldValue = uint32(1000)

var indexOptions = index.Options{
	FileName:                       "./index.db",
	PageSize:                       os.Getpagesize(),
	PreAllocatedPagePoolSize:       8192,
	AllowedPageOccupancyPercentage: 90,
}

var dbOptions = heap_file.DbOptions{
	HeapFileOptions: heap_file.HeapFileOptions{
		PageSize:                 os.Getpagesize(),
		FileName:                 "./heap.db",
		PreAllocatedPagePoolSize: 8192,
		TupleDescriptor: tuple.TupleDescriptor{
			FieldTypes: []field.FieldType{
				field.StringFieldType{}, field.StringFieldType{}, field.Uint32FieldType{},
			},
		},
	},
	IndexOptions: indexOptions,
}

var db, _ = heap_file.Open(dbOptions)

func Benchmark_Put(b *testing.B) {
	benchmarkPutKeys(b, db)
}

func Benchmark_GetFixedKey(b *testing.B) {
	if b.N == 1 {
		initializeDatabase := func(count int) {
			for iterator := 1; iterator <= count; iterator++ {
				aTuple := tuple.NewTuple()
				aTuple.AddField(field.NewStringField(stringField1Value))
				aTuple.AddField(field.NewStringField(stringField2Value))
				aTuple.AddField(field.NewUint32Field(uint32FieldValue))
				uint32FieldValue = uint32FieldValue + 1

				_, err := db.Put(aTuple)
				if err != nil {
					b.Fatalf("Failed while putting key/values during preparation for get %v", err)
				}
			}
		}
		initializeDatabase(100000)
	}
	searchKey := uint32(5000)
	benchmarkGetFixedKey(b, db, field.NewUint32Field(searchKey))
}

func Benchmark_GetKeys(b *testing.B) {
	if b.N == 1 {
		initializeDatabase := func(count int) {
			for iterator := 1; iterator <= count; iterator++ {
				aTuple := tuple.NewTuple()
				aTuple.AddField(field.NewStringField(stringField1Value))
				aTuple.AddField(field.NewStringField(stringField2Value))
				aTuple.AddField(field.NewUint32Field(uint32FieldValue))
				uint32FieldValue = uint32FieldValue + 1

				_, err := db.Put(aTuple)
				if err != nil {
					b.Fatalf("Failed while putting key/values during preparation for get %v", err)
				}
			}
		}
		initializeDatabase(100000)
	}
	benchmarkGetKeys(b, db, 100000+uint32FieldValue)
}

func benchmarkPutKeys(b *testing.B, db *heap_file.Db) {
	b.ResetTimer()
	for iterator := 0; iterator < b.N; iterator++ {
		aTuple := tuple.NewTuple()
		aTuple.AddField(field.NewStringField(stringField1Value))
		aTuple.AddField(field.NewStringField(stringField2Value))
		aTuple.AddField(field.NewUint32Field(uint32FieldValue))
		uint32FieldValue = uint32FieldValue + 1

		_, err := db.Put(aTuple)
		if err != nil {
			b.Fatalf("Failed while putting %v", err)
		}
	}
}

func benchmarkGetFixedKey(b *testing.B, db *heap_file.Db, key field.Uint32Field) {
	b.ResetTimer()
	for iterator := 0; iterator < b.N; iterator++ {
		_, err := db.GetByKey(key)
		if err != nil {
			b.Fatalf("Failed while getting %v", err)
		}
	}
}

func benchmarkGetKeys(b *testing.B, db *heap_file.Db, maxKeyValue uint32) {
	b.ResetTimer()
	startingKeyValue := uint32(1000)
	for iterator := 0; iterator < b.N; iterator++ {
		keyValue := startingKeyValue + uint32(iterator)
		if keyValue > maxKeyValue {
			keyValue = maxKeyValue
		}
		key := field.NewUint32Field(keyValue)
		_, err := db.GetByKey(key)
		if err != nil {
			b.Fatalf("Failed while getting %v", err)
		}
	}
}
