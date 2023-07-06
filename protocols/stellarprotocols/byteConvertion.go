package stellarprotocols

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dileepaj/tracified-gateway/constants"
	"github.com/dileepaj/tracified-gateway/utilities"
	"github.com/sirupsen/logrus"
)

type bitString string

var logger = utilities.NewCustomLogger()

func stringToBin(s string) (binString string) {
	for _, c := range s {
		binString = fmt.Sprintf("%s%b", binString, c)
	}
	return
}

func (b bitString) AsByteSlice() ([]byte, error) {
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
			logger.LogWriter("Error converting :"+err.Error(), constants.ERROR)
			return nil, err
		}
		out = append([]byte{byte(v)}, out...)
	}
	return out, nil
}

func (b bitString) AsHexSlice() ([]string, error) {
	var out []string
	byteSlice, sliceErr := b.AsByteSlice()
	if sliceErr != nil {
		logger.LogWriter("Error AsByteSlice() in AsHexSlice():"+sliceErr.Error(), constants.ERROR)
		return nil, sliceErr
	}
	for _, b := range byteSlice {
		out = append(out, "0x"+hex.EncodeToString([]byte{b}))
	}
	return out, nil
}

func ConvertingBinaryToByteString(str string) (string, error) {
	bitValue := bitString(str)
	byteValue, sliceErr := bitValue.AsByteSlice()
	if sliceErr != nil {
		logger.LogWriter("Error AsByteSlice() in ConvertingBinaryToByteString() :"+sliceErr.Error(), constants.ERROR)
		return "", sliceErr
	}
	return string(byteValue), nil
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

func Int8ToByteString(value uint8) (string, error) {
	binary := strconv.FormatInt(int64(value), 2)

	if len(binary) < 8 {
		// add 0s to the rest of the name
		remain := 8 - len(binary)
		setReaminder := fmt.Sprintf("%s", strings.Repeat("0", remain))
		return ConvertingBinaryToByteString(setReaminder + binary)
	} else if len(binary) == 8 {
		return ConvertingBinaryToByteString(binary)
	} else {
		return binary, errors.New("Unit length shouldbe equal to 8")
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
func UInt64ToByteString(i uint64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return string(b)
}

// return convert usign int16 to byte string
func UInt16ToByteString(i uint16) string {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(i))
	return string(b)
}

// return convert usign int32 to byte string
func UInt32ToByteString(i uint32) string {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return string(b)
}

func Float64ToByteString(f float64) string {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b[:], math.Float64bits(f))
	return string(b)
}
