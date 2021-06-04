package col

// OrderedMapItem key/value item
type OrderedMapItem struct {
	MapItem
	item *ListItem
}

// Next returns a pointer to the next item.
func (mi *OrderedMapItem) Next() *OrderedMapItem {
	return toOrderedMapItem(mi.item.Next())
}

// Prev returns a pointer to the previous item.
func (mi *OrderedMapItem) Prev() *OrderedMapItem {
	return toOrderedMapItem(mi.item.Prev())
}

func toOrderedMapItem(li *ListItem) *OrderedMapItem {
	if li == nil {
		return nil
	}
	return li.Value.(*OrderedMapItem)
}
