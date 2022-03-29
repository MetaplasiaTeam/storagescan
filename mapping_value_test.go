package storagescan

import (
	"fmt"
	"testing"
)

func Test_encodeIntString(t *testing.T) {
	v := encodeIntString("256")
	fmt.Println(v)
	v1 := encodeIntString("-256")
	fmt.Println(v1)
}
