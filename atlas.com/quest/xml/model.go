package xml

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
)

type Noder interface {
	Name() string
}

func IntFromIntegerNode(root Noder) (int, error) {
	val, ok := root.(*IntegerNode)
	if !ok {
		return 0, errors.New("invalid xml structure")
	}
	return strconv.Atoi(val.Value())
}

func StringFromStringNode(root Noder) (string, error) {
	val, ok := root.(*StringNode)
	if !ok {
		return "", errors.New("invalid xml structure")
	}
	return val.Value(), nil
}

func IntFromStringNode(root Noder) (int, error) {
	val, err := StringFromStringNode(root)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

type Parent interface {
	Noder
	Children() []Noder
	ChildByName(name string) (Noder, error)
}

type Node struct {
	name     string
	children []Noder
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Children() []Noder {
	return n.children
}

func (n *Node) initFromXML(i compositeNode) {
	n.name = i.Name
	for _, c := range i.ChildNodes {
		child := Node{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
	for _, c := range i.CanvasNodes {
		child := CanvasNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
	for _, c := range i.IntegerNodes {
		child := IntegerNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
	for _, c := range i.StringNodes {
		child := StringNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
	for _, c := range i.PointNodes {
		child := PointNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
}

func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var i compositeNode
	err := d.DecodeElement(&i, &start)
	if err != nil {
		return err
	}
	n.initFromXML(i)
	return nil
}

func (n *Node) ChildByName(name string) (Noder, error) {
	segments := strings.Split(name, "/")
	if len(segments) == 1 {
		for _, c := range n.Children() {
			if c.Name() == name {
				return c, nil
			}
		}
		return nil, errors.New("child not found")
	}

	is, err := n.ChildByName(segments[0])
	if err != nil {
		return nil, err
	}
	intermediary, ok := is.(Parent)
	if !ok {
		return nil, errors.New("child not found")
	}
	return intermediary.ChildByName(strings.Join(segments[1:], "/"))
}

type compositeNode struct {
	XMLName      xml.Name        `xml:"imgdir"`
	Name         string          `xml:"name,attr"`
	ChildNodes   []compositeNode `xml:"imgdir"`
	CanvasNodes  []canvasNode    `xml:"canvas"`
	IntegerNodes []integerNode   `xml:"int"`
	StringNodes  []stringNode    `xml:"string"`
	PointNodes   []pointNode     `xml:"vector"`
}

type CanvasNode struct {
	name     string
	width    string
	height   string
	children []Noder
}

func (n *CanvasNode) Name() string {
	return n.name
}

func (n *CanvasNode) Children() []Noder {
	return n.children
}

func (n *CanvasNode) initFromXML(i canvasNode) {
	n.name = i.Name
	for _, c := range i.IntegerNodes {
		child := IntegerNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
	for _, c := range i.PointNodes {
		child := PointNode{}
		child.initFromXML(c)
		n.children = append(n.children, &child)
	}
}

func (n *CanvasNode) ChildByName(name string) (Noder, error) {
	segments := strings.Split(name, "/")
	if len(segments) == 1 {
		for _, c := range n.Children() {
			if c.Name() == name {
				return c, nil
			}
		}
		return nil, errors.New("child not found")
	}

	is, err := n.ChildByName(segments[0])
	if err != nil {
		return nil, err
	}
	intermediary, ok := is.(Parent)
	if !ok {
		return nil, errors.New("child not found")
	}
	return intermediary.ChildByName(strings.Join(segments[1:], "/"))
}

//func (i *compositeNode) GetShort(name string, def uint16) uint16 {
//	for _, c := range i.IntegerNodes {
//		if c.Name == name {
//			res, err := strconv.ParseUint(c.Value, 10, 16)
//			if err != nil {
//				return def
//			}
//			return uint16(res)
//		}
//	}
//	return def
//}
//
//func (i *compositeNode) GetString(name string, def string) string {
//	for _, c := range i.StringNodes {
//		if c.Name == name {
//			return c.Value
//		}
//	}
//	return def
//}
//

func GetInteger(n Parent, name string) (int32, error) {
	for _, c := range n.Children() {
		val, ok := c.(*IntegerNode)
		if !ok {
			continue
		}

		if val.Name() == name {
			res, err := strconv.ParseInt(val.Value(), 10, 32)
			if err != nil {
				return 0, err
			}
			return int32(res), nil
		}
	}
	return 0, errors.New("node not found")
}

func GetBoolean(n Parent, name string) (bool, error) {
	for _, c := range n.Children() {
		val, ok := c.(*IntegerNode)
		if !ok {
			continue
		}

		if val.Name() == name {
			res, err := strconv.ParseInt(val.Value(), 10, 32)
			if err != nil {
				return false, err
			}
			return res == 1, nil
		}
	}
	return false, errors.New("node not found")
}

func GetIntegerWithDefault(n Parent, name string, def int32) int32 {
	for _, c := range n.Children() {
		val, ok := c.(*IntegerNode)
		if !ok {
			continue
		}

		if val.Name() == name {
			res, err := strconv.ParseInt(val.Value(), 10, 32)
			if err != nil {
				return def
			}
			return int32(res)
		}
	}
	return def
}

func GetString(n Parent, name string) (string, error) {
	for _, c := range n.Children() {
		val, ok := c.(*StringNode)
		if !ok {
			continue
		}

		if val.Name() == name {
			return val.Value(), nil
		}
	}
	return "", errors.New("node not found")
}

//
//func (i *compositeNode) GetFloatWithDefault(name string, def float64) float64 {
//	for _, c := range i.IntegerNodes {
//		if c.Name == name {
//			res, err := strconv.ParseFloat(c.Value, 64)
//			if err != nil {
//				return def
//			}
//			return res
//		}
//	}
//	return def
//}
//
//func (i *compositeNode) GetPoint(name string, defX int32, defY int32) (int32, int32) {
//	for _, c := range i.PointNodes {
//		if c.Name == name {
//			x, err := strconv.ParseInt(c.X, 10, 32)
//			if err != nil {
//				return defX, defY
//			}
//			y, err := strconv.ParseInt(c.Y, 10, 32)
//			if err != nil {
//				return defX, defY
//			}
//			return int32(x), int32(y)
//		}
//	}
//	return defX, defY
//}

type canvasNode struct {
	Name         string        `xml:"name,attr"`
	Width        string        `xml:"width,attr"`
	Height       string        `xml:"height,attr"`
	IntegerNodes []integerNode `xml:"int"`
	PointNodes   []pointNode   `xml:"vector"`
}

type IntegerNode struct {
	name  string
	value string
}

func (n *IntegerNode) Name() string {
	return n.name
}

func (n *IntegerNode) Value() string {
	return n.value
}

func (n *IntegerNode) initFromXML(i integerNode) {
	n.name = i.Name
	n.value = i.Value
}

type integerNode struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type StringNode struct {
	name  string
	value string
}

func (n *StringNode) Name() string {
	return n.name
}

func (n *StringNode) initFromXML(i stringNode) {
	n.name = i.Name
	n.value = i.Value
}

func (n *StringNode) Value() string {
	return n.value
}

type stringNode struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type PointNode struct {
	name string
	x    string
	y    string
}

func (n *PointNode) Name() string {
	return n.name
}

func (n *PointNode) initFromXML(i pointNode) {
	n.name = i.Name
	n.x = i.X
	n.y = i.Y
}

type pointNode struct {
	Name string `xml:"name,attr"`
	X    string `xml:"x,attr"`
	Y    string `xml:"y,attr"`
}

//
//func (i *CanvasNode) GetIntegerWithDefault(name string, def int32) int32 {
//	for _, c := range i.IntegerNodes {
//		if c.Name == name {
//			res, err := strconv.ParseUint(c.Value, 10, 32)
//			if err != nil {
//				return def
//			}
//			return int32(res)
//		}
//	}
//	return def
//}
//
//func (i *CanvasNode) GetPoint(name string, defX int32, defY int32) (int32, int32) {
//	for _, c := range i.PointNodes {
//		if c.Name == name {
//			x, err := strconv.ParseInt(c.X, 10, 32)
//			if err != nil {
//				return defX, defY
//			}
//			y, err := strconv.ParseInt(c.Y, 10, 32)
//			if err != nil {
//				return defX, defY
//			}
//			return int32(x), int32(y)
//		}
//	}
//	return defX, defY
//}
