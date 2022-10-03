package stellarprotocols

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type bitString string

func stringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%b", binString, c)
	}
	return
}

func (b bitString) AsByteSlice() []byte {
	var out []byte
	var str string

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			panic(err)
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out
}

func (b bitString) AsHexSlice() []string {
	var out []string
	byteSlice := b.AsByteSlice()
	for _, b := range byteSlice {
		out = append(out, "0x"+hex.EncodeToString([]byte{b}))
	}
	return out
}

func ConvertingBinaryToByteString(str string) string {
	bitValue := bitString(str)
	byteValue := bitValue.AsByteSlice()
	return string(byteValue)
}

func GetDataType(value string) string {
	if value == "int" {
		return "1"
	} else if value == "float64" {
		return "2"
	} else if value == "bool" {
		return "3"
	} else if value == "string" {
		return "4"
	}
	return "0"
}

func UnitToBinary(value int64) (string, error) {
	binary := strconv.FormatInt(value, 2)

	if len(binary) < 16 {
		// add 0s to the rest of the name
		remain := 16 - len(binary)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		return setReaminder + binary, nil
	} else if len(binary) == 16 {
		return binary, nil
	} else {
		return binary, errors.New("Unit length shouldbe equal to 16")
	}
}

func StringToBinary(value int64) (string, error) {
	binary := strconv.FormatInt(value, 2)

	if len(binary) < 8 {
		// add 0s to the rest of the name
		remain := 8 - len(binary)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		return setReaminder + binary, nil
	} else if len(binary) == 8 {
		return binary, nil
	} else {
		return binary, errors.New("Unit length shouldbe equal to 8")
	}
}

func IDToBinary(value int64) (string, error) {
	binary := strconv.FormatInt(value, 2)

	if len(binary) < 64 {
		// add 0s to the rest of the name
		remain := 64 - len(binary)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		return setReaminder + binary, nil
	} else if len(binary) == 64 {
		return binary, nil
	} else {
		return binary, errors.New("Unit length shouldbe equal to 64")
	}
}

func TenantIDToBinary(value int64) (string, error) {
	binary := strconv.FormatInt(value, 2)

	if len(binary) < 32 {
		// add 0s to the rest of the name
		remain := 32 - len(binary)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		return setReaminder + binary, nil
	} else if len(binary) == 32 {
		return binary, nil
	} else {
		return binary, errors.New("Unit length shouldbe equal to 32")
	}
}

func ByteStingToInteger(byteValue string) (int64, error) {
	strVal := []byte(byteValue)
	encodedString := hex.EncodeToString(strVal)
	intValue, err := strconv.ParseInt(encodedString, 16, 64)
	if err != nil {
		logrus.Printf("Conversion failed: %s\n", err)
		return 0, errors.New("Conversion failed: %s\n" + err.Error())
	} else {
		return intValue, nil
	}
}

// return convert usign int64 to byte string
func UInt64ToByteString(i int64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return string(b)
}
