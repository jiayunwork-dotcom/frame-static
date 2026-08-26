package model

type Node struct {
	ID        string  `json:"id"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Restraint [3]bool `json:"restraint"`
	Support   string  `json:"support"`
}

func (n Node) EffectiveRestraint() [3]bool {
	if n.Restraint != [3]bool{false, false, false} {
		return n.Restraint
	}
	switch n.Support {
	case "pin", "pinned", "hinge", "铰支":
		return [3]bool{true, true, false}
	case "fixed", "固支":
		return [3]bool{true, true, true}
	case "roller-x", "滑动-x", "roller":
		return [3]bool{false, true, false}
	case "roller-y", "滑动-y":
		return [3]bool{true, false, false}
	default:
		return n.Restraint
	}
}

type Element struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	E    float64   `json:"E"`
	A    float64   `json:"A"`
	I    float64   `json:"I"`
	Dist *DistLoad `json:"dist"`
}

type DistLoad struct {
	Q float64 `json:"q"`
}

type NodeLoad struct {
	Node string  `json:"node"`
	FX   float64 `json:"fx"`
	FY   float64 `json:"fy"`
	MZ   float64 `json:"mz"`
}

type Model struct {
	Nodes    []Node     `json:"nodes"`
	Elements []Element  `json:"elements"`
	Loads    []NodeLoad `json:"loads"`
}
