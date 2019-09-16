package keypath

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeyPathDecoder(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type SubStruct struct {
			Key string `json:"key"`
		}

		type Struct struct {
			SubStruct
			Name   string     `json:"name"`
			Slice  []string   `json:"slice"`
			Sub    SubStruct  `json:"sub"`
			PtrSub *SubStruct `json:"ptrSub,omitempty"`
		}

		s := &Struct{}
		s.Slice = []string{"!"}

		decoder := NewKeyPathDecoder(map[string]string{
			"name":       "name",
			"slice[0]":   "slice",
			"key":        "key",
			"ptrSub.key": "ptrKey",
			"sub.key":    "subKey",
		})

		err := decoder.Decode(s)
		require.NoError(t, err)

		expect := &Struct{
			Name:  "name",
			Slice: []string{"slice"},
			Sub: SubStruct{
				Key: "subKey",
			},
			PtrSub: &SubStruct{
				Key: "ptrKey",
			},
			SubStruct: SubStruct{
				Key: "key",
			},
		}

		require.Equal(t, expect, s)
	})

	t.Run("deep struct", func(t *testing.T) {
		type DepthStruct struct {
			A struct {
				B *struct {
					C string `json:"c"`
				} `json:"b"`
			} `json:"a"`
		}

		decoder := NewKeyPathDecoder(map[string]string{
			"a.b.c": "a.b.c",
		})

		s := &DepthStruct{}
		err := decoder.Decode(s)
		require.NoError(t, err)

		require.Equal(t, "a.b.c", s.A.B.C)
	})

	t.Run("slice && map", func(t *testing.T) {
		type SubStruct struct {
			Key string `json:"key"`
		}

		type SliceStruct struct {
			Slice           []string          `json:"slice"`
			Map             map[string]string `json:"map"`
			SliceWithStruct []SubStruct       `json:"sliceWithStruct"`
		}

		decoder := NewKeyPathDecoder(map[string]string{
			"slice[1]":               "slice",
			"sliceWithStruct[0].key": "k",
			"map.key":                "mapKey",
			"map.key1":               "mapKey1",
		})

		s := &SliceStruct{}
		s.SliceWithStruct = []SubStruct{{}}
		s.Map = map[string]string{
			"key": "1",
		}

		err := decoder.Decode(s)
		require.NoError(t, err)

		require.Equal(t, 0, len(s.Slice))
		require.Equal(t, "k", s.SliceWithStruct[0].Key)
		require.Equal(t, "mapKey", s.Map["key"])
		require.Equal(t, "mapKey1", s.Map["key1"])
	})
}
