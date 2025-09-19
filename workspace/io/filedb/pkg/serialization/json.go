package serialization

import "encoding/json"

func JSONSerialize[T any](object T) ([]byte, error) {
	body, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func JSONSerializePretty[T any](object T) ([]byte, error) {
	body, err := json.MarshalIndent(object, "", "	")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func JSONDeserialize[T any](body []byte, object *T) error {
	err := json.Unmarshal(body, object)
	if err != nil {
		return err
	}
	return nil
}
