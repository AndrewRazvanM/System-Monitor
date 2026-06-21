package layoutmanager

func (s *Split) UpdatePos(g Geometry, isVisible bool) {
	n := len(s.Children)
	if n == 0 {
		return
	}

	total := g.W
	if s.Direction == Column {
		total = g.H
	}

	weights := make([]float32, n)
	var sum float32
	for i, e := range s.Children {
		w := e.Weight
		if w == 0 {
			w = 1.0
		}
		weights[i] = w
		sum += w
	}

	offset := 0
	for i, e := range s.Children {
		share := int(float32(total) * weights[i] / sum)
		if i == n-1 {
			share = total - offset // last child absorbs the rounding remainder
		}

		child := g
		if s.Direction == Row {
			child.X = g.X + offset
			child.W = share
		} else {
			child.Y = g.Y + offset
			child.H = share
		}
		offset += share

		e.Node.UpdatePos(child, isVisible)
	}
}