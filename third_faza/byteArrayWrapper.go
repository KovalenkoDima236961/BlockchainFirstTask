package third_faza

import (
	"hash/fnv"
)

type ByteArrayWrapper struct {
	contents []byte
}

func NewByteArrayWrapper(content []byte) *ByteArrayWrapper {
	copied := make([]byte, len(content))
	copy(copied, content)
	return &ByteArrayWrapper{contents: copied}
}

func (wrapper *ByteArrayWrapper) Equals(other *ByteArrayWrapper) bool {
	if other == nil {
		return false
	}

	b := other.contents

	if wrapper.contents == nil {
		if b == nil {
			return true
		} else {
			return false
		}
	} else {
		if b == nil {
			return false
		} else {
			if len(wrapper.contents) != len(b) {
				return false
			}
			for i := range len(b) {
				if wrapper.contents[i] != b[i] {
					return false
				}
			}
			return true
		}
	}
}

// HashCode function
func (wrapper *ByteArrayWrapper) HashCode() uint32 {
	hasher := fnv.New32a()                   // What is it
	_, err := hasher.Write(wrapper.contents) // what is it
	if err != nil {
		return 0
	}
	return hasher.Sum32()
}
