package perimeter

import "testing"

func TestPerimeter(t *testing.T) {
	rectangle := Rectangle{10.0, 10.0}
	got := Perimeter(rectangle)
	want := 40.0

	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

func TestArea(t *testing.T) {

	areaTests := []struct{
		shape	Shape
		want	float64
	} {
		{shape: Rectangle{12.0, 6.0}, want: 72.0},
		{shape: Circle{10}, want: 314.1592653589793},
		{shape: Triangle{12,6}, want: 36},
	}

	for _, v := range areaTests {
		got := v.shape.Area()
		if v.want != got {
			t.Errorf("%#v got %g want %g", v, got, v.want)
		}
	}
}
