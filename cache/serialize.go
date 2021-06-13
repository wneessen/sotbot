package cache

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func SerializeObj(o interface{}) (string, error) {
	var serializeString bytes.Buffer
	gobEnc := gob.NewEncoder(&serializeString)
	if err := gobEnc.Encode(o); err != nil {
		return "", err
	}
	serializeBase64 := base64.StdEncoding.EncodeToString(serializeString.Bytes())

	return serializeBase64, nil
}

func DeserializeObj(s string, o interface{}) error {
	var objBuffer bytes.Buffer
	objString, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	objBuffer.Write(objString)
	gobDec := gob.NewDecoder(&objBuffer)
	if err := gobDec.Decode(o); err != nil {
		return err
	}

	return nil
}
