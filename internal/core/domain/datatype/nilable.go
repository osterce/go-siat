package datatype

import (
	"encoding/json"
	"encoding/xml"
)

// Nilable es un tipo genérico para manejar valores que pueden ser nulos en el XML del SIAT.
// Es necesario porque el SIAT requiere que ciertos campos opcionales se envíen con el atributo
// xsi:nil="true" cuando no tienen valor, en lugar de omitir la etiqueta por completo.
type Nilable[T any] struct {
	Value *T
}

// MarshalXML implementa la interfaz xml.Marshaler para manejar la nulidad explícita con xsi:nil.
func (n Nilable[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Value == nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(*n.Value, start)
}

// UnmarshalXML implementa la interfaz xml.Unmarshaler para manejar la nulidad explícita con xsi:nil.
func (n *Nilable[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var isNil bool
	for _, attr := range start.Attr {
		if attr.Name.Local == "nil" && attr.Value == "true" {
			isNil = true
			break
		}
	}

	if isNil {
		n.Value = nil
		return d.Skip()
	}

	var val T
	if err := d.DecodeElement(&val, &start); err != nil {
		return err
	}
	n.Value = &val
	return nil
}

// MarshalJSON implementa la interfaz json.Marshaler.
func (n Nilable[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Value)
}

// UnmarshalJSON implementa la interfaz json.Unmarshaler.
func (n *Nilable[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &n.Value)
}
