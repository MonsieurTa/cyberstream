package bencode

type Decoder map[string]interface{}

func (d Decoder) GetString(key string) string {
	v, ok := d[key]
	if !ok {
		return ""
	}
	rv, ok := v.(string)
	if !ok {
		return ""
	}
	return rv
}

func (d Decoder) GetList(key string) []interface{} {
	v, ok := d[key]
	if !ok {
		return nil
	}
	rv, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return rv
}

func (d Decoder) GetInt(key string) int {
	v, ok := d[key]
	if !ok {
		return 0
	}
	rv, ok := v.(int64)
	if !ok {
		return 0
	}
	return int(rv)
}

func (d Decoder) GetDict(key string) map[string]interface{} {
	v, ok := d[key]
	if !ok {
		return nil
	}
	rv, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return rv
}
