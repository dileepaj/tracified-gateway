package commons

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"math"

	"github.com/sirupsen/logrus"
)

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func ByteArryToFloat64(b []byte, name string) (byteLength int, extractedValueFormByte float64) {
	c := []byte(b)
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", c)
	logrus.Info("\t\tByte array length:", len(c))
	base64EncodedStr := base64.StdEncoding.EncodeToString(c)
	decodedToBytes, _ := base64.StdEncoding.DecodeString(base64EncodedStr)
	logrus.Info("\t\tValue in float:", Float64frombytes(decodedToBytes))
	extractedValueFormByte = Float64frombytes(decodedToBytes)
	byteLength = len(c)
	return
}

func ByteArryToInt64(b []byte, name string) (byteLength int, extractedValueFormByte int64) {
	c := []byte(b)
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", c)
	logrus.Info("\t\tByte array length:", len(b))
	i := int64(binary.LittleEndian.Uint64(c))
	logrus.Info("\t\tActual value in int:", i)
	extractedValueFormByte = i
	byteLength = len(c)
	return
}

func ByteArryToInt32(b []byte, name string) (byteLength int, extractedValueFormByte int32) {
	c := []byte(b)
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", c)
	logrus.Info("\t\tByte array length:", len(b))
	i := int32(binary.LittleEndian.Uint32(c))
	logrus.Info("\t\tActual value in int:", i)
	extractedValueFormByte = i
	byteLength = len(c)
	return
}

func ByteArryToInt16(b []byte, name string) (byteLength int, extractedValueFormByte int16) {
	c := []byte(b)
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", c)
	logrus.Info("\t\tByte array length:", len(b))
	i := int16(binary.LittleEndian.Uint16(c))
	logrus.Info("\t\tActual value in int:", i)
	extractedValueFormByte = i
	byteLength = len(c)
	return
}

func ByteArryToHexString(b []byte, name string) (byteLength int, extractedValueFormByte string) {
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", b)
	logrus.Info("\t\tByte array length:", len(b))
	myString := hex.EncodeToString(b)
	logrus.Info("\t\tActual value as a hex string:", myString)
	extractedValueFormByte = myString
	byteLength = len(b)
	return
}

func ByteArryToString(b []byte, name string) (byteLength int, extractedValueFormByte string) {
	logrus.Info()
	logrus.Info("\t", name)
	logrus.Info("\t\tByte array:", b)
	logrus.Info("\t\tByte array length:", len(b))
	myString := string(b)
	logrus.Info("\t\tActual value as a string:", myString)
	extractedValueFormByte = myString
	byteLength = len(b)
	return
}
