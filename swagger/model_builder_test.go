package swagger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestJsonTags(t *testing.T) {
	type X struct {
		A string
		B string `json:"-"`
		C int    `json:",string"`
		D int    `json:","`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A",
    "C",
    "D"
   ],
   "properties": {
    "A": {
     "type": "string",
     "description": ""
    },
    "C": {
     "type": "string",
     "description": "(int as string)"
    },
    "D": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestJsonTagOmitempty(t *testing.T) {
	type X struct {
		A int `json:",omitempty"`
		B int `json:"C,omitempty"`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "properties": {
    "A": {
     "type": "integer",
     "description": ""
    },
    "C": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestJsonTagName(t *testing.T) {
	type X struct {
		A string `json:"B"`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "string",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousStruct(t *testing.T) {
	type X struct {
		A struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "properties": {
    "A": {
     "type": "swagger.X.A",
     "description": ""
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousPtrStruct(t *testing.T) {
	type X struct {
		A *struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "swagger.X.A",
     "description": ""
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousArrayStruct(t *testing.T) {
	type X struct {
		A []struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.X.A"
     }
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousPtrArrayStruct(t *testing.T) {
	type X struct {
		A *[]struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.X.A"
     }
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

// go test -v -test.run TestEmbeddedStruct_Issue98 ...swagger
func TestEmbeddedStruct_Issue98(t *testing.T) {
	type Y struct {
		A int
	}
	type X struct {
		Y
	}
	testJsonFromStruct(t, X{}, `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`)
}

type File struct {
	History     []File
	HistoryPtrs []*File
}

// go test -v -test.run TestRecursiveStructure ...swagger
func TestRecursiveStructure(t *testing.T) {
	testJsonFromStruct(t, File{}, `{
  "swagger.File": {
   "id": "swagger.File",
   "required": [
    "History",
    "HistoryPtrs"
   ],
   "properties": {
    "History": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.File"
     }
    },
    "HistoryPtrs": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.File.HistoryPtrs"
     }
    }
   }
  },
  "swagger.File.HistoryPtrs": {
   "id": "swagger.File.HistoryPtrs",
   "properties": {}
  }
 }`)
}

//1
type A1 struct {
	B struct {
		Id int
	}
}

//2
type A2 struct {
	C `json:"B"`
}
type C struct{ Id int }

//3
type A3 struct {
	B D
}
type D struct {
	Id int
}

// go test -v -test.run TestEmbeddedStructA1 ...swagger
func TestEmbeddedStructA1(t *testing.T) {
	output, _ := json.MarshalIndent(A1{}, " ", " ")
	print(string(output))
	testJsonFromStruct(t, A1{}, `{
  "swagger.A1": {
   "id": "swagger.A1",
   "properties": {
    "B": {
     "type": "swagger.A1.B",
     "description": ""
    }
   }
  },
  "swagger.A1.B": {
   "id": "swagger.A1.B",
   "required": [
    "Id"
   ],
   "properties": {
    "Id": {
     "type": "integer",
     "description": ""
    }
   }
  }
 }`)
}

// Utils

func testJsonFromStruct(t *testing.T, sample interface{}, expectedJson string) {
	compareJson(t, false, modelsFromStruct(sample), expectedJson)
}

func modelsFromStruct(sample interface{}) map[string]Model {
	models := map[string]Model{}
	builder := modelBuilder{models}
	builder.addModel(reflect.TypeOf(sample), "")
	return models
}

func compareJson(t *testing.T, flatCompare bool, value interface{}, expectedJsonAsString string) {
	var output []byte
	var err error
	if flatCompare {
		output, err = json.Marshal(value)
	} else {
		output, err = json.MarshalIndent(value, " ", " ")
	}
	if err != nil {
		t.Error(err.Error())
		return
	}
	actual := string(output)
	if actual != expectedJsonAsString {
		t.Errorf("Mismatch JSON doc")
		// Use simple fmt to create a pastable output :-)
		fmt.Println("---- expected -----")
		fmt.Println(expectedJsonAsString)
		fmt.Println("---- actual -----")
		fmt.Println(actual)
	}
}