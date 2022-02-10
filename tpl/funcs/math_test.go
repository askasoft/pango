package funcs

import "testing"

func TestAdd_int_int(t *testing.T) {
	result, err := Add(1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(3) {
		t.Errorf("expected %d to be %d", result, int64(3))
	}
}

func TestAdd_int_uint(t *testing.T) {
	result, err := Add(1, uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(3) {
		t.Errorf("expected %d to be %d", result, int64(3))
	}
}

func TestAdd_int_float(t *testing.T) {
	result, err := Add(1, 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(3) {
		t.Errorf("expected %f to be %f", result, float64(3))
	}
}

func TestAdd_uint_int(t *testing.T) {
	result, err := Add(uint(1), 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(3) {
		t.Errorf("expected %d to be %d", result, int64(3))
	}
}

func TestAdd_uint_uint(t *testing.T) {
	result, err := Add(uint(1), uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != uint64(3) {
		t.Errorf("expected %d to be %d", result, uint64(3))
	}
}

func TestAdd_uint_float(t *testing.T) {
	result, err := Add(uint(1), 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(3) {
		t.Errorf("expected %f to be %f", result, float64(3))
	}
}

func TestAdd_float_int(t *testing.T) {
	result, err := Add(1.0, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(3) {
		t.Errorf("expected %f to be %f", result, float64(3))
	}
}

func TestAdd_float_uint(t *testing.T) {
	result, err := Add(1.0, uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(3) {
		t.Errorf("expected %f to be %f", result, float64(3))
	}
}

func TestAdd_float_float(t *testing.T) {
	result, err := Add(1.0, 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(3) {
		t.Errorf("expected %f to be %f", result, float64(3))
	}
}

func TestAdd_string_int(t *testing.T) {
	_, err := Add("foo", 2)
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "add: unknown type for \"foo\" (string)"
	if err.Error() != expected {
		t.Errorf("expected %q to be %q", err.Error(), expected)
	}
}

func TestSubtract_int_int(t *testing.T) {
	result, err := Subtract(1, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(1) {
		t.Errorf("expected %d to be %d", result, int64(1))
	}
}

func TestSubtract_int_uint(t *testing.T) {
	result, err := Subtract(1, uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(1) {
		t.Errorf("expected %d to be %d", result, int64(1))
	}
}

func TestSubtract_int_float(t *testing.T) {
	result, err := Subtract(1, 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(1) {
		t.Errorf("expected %f to be %f", result, float64(1))
	}
}

func TestSubtract_uint_int(t *testing.T) {
	result, err := Subtract(uint(1), 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(1) {
		t.Errorf("expected %d to be %d", result, int64(1))
	}
}

func TestSubtract_uint_uint(t *testing.T) {
	result, err := Subtract(uint(1), uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != uint64(1) {
		t.Errorf("expected %d to be %d", result, uint64(1))
	}
}

func TestSubtract_uint_float(t *testing.T) {
	result, err := Subtract(uint(1), 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(1) {
		t.Errorf("expected %f to be %f", result, float64(1))
	}
}

func TestSubtract_float_int(t *testing.T) {
	result, err := Subtract(1.0, 2)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(1) {
		t.Errorf("expected %f to be %f", result, float64(1))
	}
}

func TestSubtract_float_uint(t *testing.T) {
	result, err := Subtract(1.0, uint(2))
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(1) {
		t.Errorf("expected %f to be %f", result, float64(1))
	}
}

func TestSubtract_float_float(t *testing.T) {
	result, err := Subtract(1.0, 2.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(1) {
		t.Errorf("expected %f to be %f", result, float64(1))
	}
}

func TestSubtract_string_int(t *testing.T) {
	_, err := Subtract("foo", 2)
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "subtract: unknown type for \"foo\" (string)"
	if err.Error() != expected {
		t.Errorf("expected %q to be %q", err.Error(), expected)
	}
}

func TestMultiply_int_int(t *testing.T) {
	result, err := Multiply(2, 3)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(6) {
		t.Errorf("expected %d to be %d", result, int64(6))
	}
}

func TestMultiply_int_uint(t *testing.T) {
	result, err := Multiply(2, uint(3))
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(6) {
		t.Errorf("expected %d to be %d", result, int64(6))
	}
}

func TestMultiply_int_float(t *testing.T) {
	result, err := Multiply(2, 3.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(6) {
		t.Errorf("expected %f to be %f", result, float64(6))
	}
}

func TestMultiply_uint_int(t *testing.T) {
	result, err := Multiply(uint(2), 3)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(6) {
		t.Errorf("expected %d to be %d", result, int64(6))
	}
}

func TestMultiply_uint_uint(t *testing.T) {
	result, err := Multiply(uint(2), uint(3))
	if err != nil {
		t.Fatal(err)
	}

	if result != uint64(6) {
		t.Errorf("expected %d to be %d", result, uint64(6))
	}
}

func TestMultiply_uint_float(t *testing.T) {
	result, err := Multiply(uint(2), 3.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(6) {
		t.Errorf("expected %f to be %f", result, float64(6))
	}
}

func TestMultiply_float_int(t *testing.T) {
	result, err := Multiply(2.0, 3)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(6) {
		t.Errorf("expected %f to be %f", result, float64(6))
	}
}

func TestMultiply_float_uint(t *testing.T) {
	result, err := Multiply(2.0, uint(3))
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(6) {
		t.Errorf("expected %f to be %f", result, float64(6))
	}
}

func TestMultiply_float_float(t *testing.T) {
	result, err := Multiply(2.0, 3.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(6) {
		t.Errorf("expected %f to be %f", result, float64(6))
	}
}

func TestMultiply_string_int(t *testing.T) {
	_, err := Multiply("foo", 2)
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "multiply: unknown type for \"foo\" (string)"
	if err.Error() != expected {
		t.Errorf("expected %q to be %q", err.Error(), expected)
	}
}

func TestDivide_int_int(t *testing.T) {
	result, err := Divide(2, 10)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(5) {
		t.Errorf("expected %d to be %d", result, int64(5))
	}
}

func TestDivide_int_uint(t *testing.T) {
	result, err := Divide(2, uint(10))
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(5) {
		t.Errorf("expected %d to be %d", result, int64(5))
	}
}

func TestDivide_int_float(t *testing.T) {
	result, err := Divide(2, 10.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(5) {
		t.Errorf("expected %f to be %f", result, float64(5))
	}
}

func TestDivide_uint_int(t *testing.T) {
	result, err := Divide(uint(2), 10)
	if err != nil {
		t.Fatal(err)
	}

	if result != int64(5) {
		t.Errorf("expected %d to be %d", result, int64(5))
	}
}

func TestDivide_uint_uint(t *testing.T) {
	result, err := Divide(uint(2), uint(10))
	if err != nil {
		t.Fatal(err)
	}

	if result != uint64(5) {
		t.Errorf("expected %d to be %d", result, uint64(5))
	}
}

func TestDivide_uint_float(t *testing.T) {
	result, err := Divide(uint(2), 10.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(5) {
		t.Errorf("expected %f to be %f", result, float64(5))
	}
}

func TestDivide_float_int(t *testing.T) {
	result, err := Divide(2.0, 10)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(5) {
		t.Errorf("expected %f to be %f", result, float64(5))
	}
}

func TestDivide_float_uint(t *testing.T) {
	result, err := Divide(2.0, uint(10))
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(5) {
		t.Errorf("expected %f to be %f", result, float64(5))
	}
}

func TestDivide_float_float(t *testing.T) {
	result, err := Divide(2.0, 10.0)
	if err != nil {
		t.Fatal(err)
	}

	if result != float64(5) {
		t.Errorf("expected %f to be %f", result, float64(5))
	}
}

func TestDivide_string_int(t *testing.T) {
	_, err := Divide("foo", 2)
	if err == nil {
		t.Fatal("expected error, but nothing was returned")
	}

	expected := "divide: unknown type for \"foo\" (string)"
	if err.Error() != expected {
		t.Errorf("expected %q to be %q", err.Error(), expected)
	}
}
