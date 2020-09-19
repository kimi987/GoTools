package must

import "github.com/lightpaw/logrus"

type marshaler interface {
	Marshal() (dAtA []byte, err error)
}

func Marshal(m marshaler) []byte {
	data, err := m.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("must.Marshal fail")
	}
	return data
}

func MarshalInterface(i interface{}) []byte {
	m, ok := i.(marshaler)
	if !ok {
		logrus.Errorf("must.MarshalInterface cast type fail")
	}

	data, err := m.Marshal()
	if err != nil {
		logrus.WithError(err).Errorf("must.Marshal fail")
	}
	return data
}

func MarshalArray(ma ...marshaler) [][]byte {
	datas := make([][]byte, len(ma))

	for i, m := range ma {
		datas[i] = Marshal(m)
	}

	return datas
}
