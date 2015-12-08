package abstract

import (
	"bytes"
	"fmt"
)

/*
Adjust marshaling size for the Secret-structure - needs to be adjusted
with changes to how suites are stored
*/
func (s *Secret) MarshalSize() int {
	return s.SecretInterface.MarshalSize() + 8
}

/*
Marshal the suite, then the binary representation of the secret.
*/
func (s *Secret) MarshalBinary() (data []byte, err error) {
	var b bytes.Buffer
	bvalue, err := s.SecretInterface.MarshalBinary()
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(&b, s.GetSuite().String(), len(bvalue))
	b.Write(bvalue)
	return b.Bytes(), nil
}

/*
Unmarshal first the suite, create the secret, and unmarshal the
binary representation of the secret.
*/
func (s *Secret) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	b := bytes.NewBuffer(data)
	var length int
	var suiteStr string
	_, err := fmt.Fscanln(b, &suiteStr, &length)
	bvalue := make([]byte, length)
	b.Read(bvalue)
	suite, err := StringToSuite(suiteStr)
	if err != nil {
		return err
	}
	secret := suite.Secret()
	s.SecretInterface = secret.SecretInterface
	s.SecretInterface.SetSuite(suite)
	s.SecretInterface.UnmarshalBinary(bvalue)
	return err
}

/*
Adjust the marshal-size with regard to our storage of the suite
*/
func (p *Point) MarshalSize() int {
	return p.PointInterface.MarshalSize() + 8
}

/*
First write the suite, then the binary representation of the point.
*/
func (p *Point) MarshalBinary() (data []byte, err error) {
	var b bytes.Buffer
	bvalue, err := p.PointInterface.MarshalBinary()
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(&b, p.GetSuite().String(), len(bvalue))
	b.Write(bvalue)
	return b.Bytes(), nil
}

/*
Fetch the suite, create the point if it doesn't exist yet,
then unmarshal the binary representation into it.
*/
func (p *Point) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	b := bytes.NewBuffer(data)
	var length int
	var suiteStr string
	_, err := fmt.Fscanln(b, &suiteStr, &length)
	bvalue := make([]byte, length)
	b.Read(bvalue)
	suite, err := StringToSuite(suiteStr)
	if err != nil {
		return err
	}
	point := suite.Point()
	point.PointInterface.UnmarshalBinary(bvalue)
	if p.PointInterface != nil {
		p.Null()
		p.Add(p, point)
	} else {
		p.PointInterface = point.PointInterface
		p.SetSuite(suite)
	}
	return err
}