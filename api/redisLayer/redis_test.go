package redisLayer

import (
	"fmt"
	"os"
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

	err := FlushDb()
	if err != nil {
		fmt.Printf("Error flushing redis db: %s", err.Error())
	}
}

func TestKeyReadWrite(t *testing.T) {
	for i, key := range testKeys {
		SetKeyString(key, testValues[i])

		exists, err := Exists(key)
		if err != nil {
			fmt.Printf("Error checking if key exists: %s", err.Error())
		}
		if !exists {
			t.Errorf("Key %s does not exist in Redis db", key)
		}

		value, err := GetKeyString(key)
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

		SetKeyBytes(key, []byte(testValues[i]))

		exists, err := Exists(key)
		if err != nil {
			fmt.Printf("Error checking if key exists: %s", err.Error())
		}
		if !exists {
			t.Errorf("Key %s does not exist in Redis db", key)
		}

		value, err := GetKeyBytes(key)
		if err != nil {
			fmt.Printf("Error looking up value for key: %s", err.Error())
		}

		if string(value) != testValues[i] {
			t.Errorf("Unexpected value for key %s in Redis db", key)
		}
	}
}

func teardown() {
	err := FlushDb()
	if err != nil {
		fmt.Printf("Error flushing redis db: %s", err.Error())
	}
}
