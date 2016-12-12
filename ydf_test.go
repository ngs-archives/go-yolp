package yolp

import (
	"encoding/xml"
	"testing"
)

func TestUnmarshalCoordinates(t *testing.T) {
	var c Coordinates
	err := xml.Unmarshal([]byte("<Coordinates>35.62172852580437,139.6999476850032</Coordinates>"), &c)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	Test{35.62172852580437, c.Latitude}.Compare(t)
	Test{139.6999476850032, c.Longitude}.Compare(t)
	err = xml.Unmarshal([]byte("<Coordinates>35.62172852580437,---</Coordinates>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`strconv.ParseFloat: parsing "---": invalid syntax`, err.Error()}.Compare(t)
	err = xml.Unmarshal([]byte("<Coordinates>---,139.6999476850032</Coordinates>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`strconv.ParseFloat: parsing "---": invalid syntax`, err.Error()}.Compare(t)
	err = xml.Unmarshal([]byte("<Coordinates>---</Coordinates>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`Invalid format ---`, err.Error()}.Compare(t)
}

func TestUnmarshalRadius(t *testing.T) {
	var c Radius
	err := xml.Unmarshal([]byte("<Radius>0.01,0.02</Radius>"), &c)
	if err != nil {
		t.Errorf("Got error %v", err)
	}
	Test{0.01, c.Horizontal}.Compare(t)
	Test{0.02, c.Vertical}.Compare(t)
	err = xml.Unmarshal([]byte("<Radius>0.01,---</Radius>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`strconv.ParseFloat: parsing "---": invalid syntax`, err.Error()}.Compare(t)
	err = xml.Unmarshal([]byte("<Radius>---,0.02</Radius>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`strconv.ParseFloat: parsing "---": invalid syntax`, err.Error()}.Compare(t)
	err = xml.Unmarshal([]byte("<Radius>---</Radius>"), &c)
	if err == nil {
		t.Errorf("Expected not nil but got nil: %v", c)
	}
	Test{`Invalid format ---`, err.Error()}.Compare(t)
}
