package redis_test

import (
	"fmt"
	"os"
	"otherside/api/redisLayer"
	"testing"
)

var testKeys []string
var testValues []string

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {
	testKeys = append(testKeys, "Apple", "Orange", "Banana")
	testValues = append(testValues, "Red", "Orange", "Yellow")

	err := redisLayer.FlushDb()
	if err != nil {
		fmt.Printf("Error flushing redis db: %s", err.Error())
	}
}

func TestKeyReadWrite(t *testing.T) {
	for i, key := range testKeys {
		redisLayer.SetKeyString(key, testValues[i])

		exists, err := redisLayer.Exists(key)
		if err != nil {
			fmt.Printf("Error checking if key exists: %s", err.Error())
		}
		if !exists {
			t.Errorf("Key %s does not exist in Redis db", key)
		}

		value, err := redisLayer.GetKeyString(key)
		if err != nil {
			fmt.Printf("Error looking up value for key: %s", err.Error())
		}

		if value != testValues[i] {
			t.Errorf("Unexpected value for key %s in Redis db", key)
		}
	}
}

func TestBytesReadWrite(t *testing.T) {
	for i, key := range testKeys {

		redisLayer.SetKeyBytes(key, []byte(testValues[i]))

		exists, err := redisLayer.Exists(key)
		if err != nil {
			fmt.Printf("Error checking if key exists: %s", err.Error())
		}
		if !exists {
			t.Errorf("Key %s does not exist in Redis db", key)
		}

		value, err := redisLayer.GetKeyBytes(key)
		if err != nil {
			fmt.Printf("Error looking up value for key: %s", err.Error())
		}

		if string(value) != testValues[i] {
			t.Errorf("Unexpected value for key %s in Redis db", key)
		}
	}
}

func teardown() {
	err := redisLayer.FlushDb()
	if err != nil {
		fmt.Printf("Error flushing redis db: %s", err.Error())
	}
}
