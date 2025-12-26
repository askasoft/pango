package ref

import "testing"

type M map[string]any

func TestMapGet(t *testing.T) {
	// test one level map
	m1 := map[string]int64{
		"a": 1,
		"1": 2,
	}

	if res, err := MapGet(m1, "a"); err == nil {
		if res.(int64) != 1 {
			t.Errorf("Should return 1, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, "1"); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if res, err := MapGet(m1, 1); err == nil {
		if res.(int64) != 2 {
			t.Errorf("Should return 2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 2 level map
	m2 := M{
		"1": map[string]float64{
			"2": 3.5,
		},
	}

	if res, err := MapGet(m2, 1, 2); err == nil {
		if res.(float64) != 3.5 {
			t.Errorf("Should return 3.5, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// test 5 level map
	m5 := M{
		"1": M{
			"2": M{
				"3": M{
					"4": M{
						"5": 1.2,
					},
				},
			},
		},
	}

	if res, err := MapGet(m5, 1, 2, 3, 4, 5); err == nil {
		if res.(float64) != 1.2 {
			t.Errorf("Should return 1.2, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	// check whether element not exists in map
	if res, err := MapGet(m5, 5, 4, 3, 2, 1); err == nil {
		if res != nil {
			t.Errorf("Should return nil, but return %v", res)
		}
	} else {
		t.Errorf("Error happens %v", err)
	}
}

func TestMapSet(t *testing.T) {
	// test one level map
	m1 := map[string]int64{
		"a": 0,
		"1": 0,
	}

	if err := MapSet(m1, "a", "1"); err == nil {
		if m1["a"] != 1 {
			t.Errorf(`m1["a"] = %v, want 1`, m1["a"])
		}
	} else {
		t.Errorf("Error happens %v", err)
	}

	if err := MapSet(m1, 1, "1"); err == nil {
		if m1["1"] != 1 {
			t.Errorf(`m1["1"] = %v, want 1`, m1["1"])
		}
	} else {
		t.Errorf("Error happens %v", err)
	}
}
