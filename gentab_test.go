package gentab

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

type (
	testLabel struct {
		id          int64  `tab:"id"`
		Name        string `tab:"good"`
		Description string `tab:"-"`
		Test        string
	}

	testLabeln struct {
		testLabel
		id2 int64 `tab:"id2"`
	}

	testFooterS []testFooter

	testFooter struct {
		id     int64  `tab:"-"`
		name   string `tab:"name"`
		desc   string `tab:"desc"`
		amount float64
	}
)

func (t *testFooterS) Footer() ([]string, error) {
	sum := float64(0.0)
	for _, t := range *t {
		sum += t.amount
	}
	sums := fmt.Sprintf("%.2f", sum)
	return []string{"", "total", sums}, nil
}

var (
	exheader     = []string{"id", "good", "Test"}
	exheaderNest = []string{"id", "good", "Test", "id2"}
	exbodyStruct = [][]string{
		[]string{"1", "nam", "hello"},
	}
	l  = testLabel{1, "nam", "desc", "hello"}
	ln = testLabeln{l, 2}
	ls = []testLabel{
		testLabel{1, "line1", "desc1", "test1"},
		testLabel{2, "line2", "desc2", "test2"},
	}
	exbodySlice = [][]string{
		[]string{"1", "line1", "test1"},
		[]string{"2", "line2", "test2"},
	}
	exfooter = []string{}
	ts       = testFooterS{
		testFooter{1, "test1", "testing", 29.2},
		testFooter{1, "test1", "testing2", 30.8},
	}

	toPrint = []reflect.Type{
		reflect.TypeOf(testLabeln{}),
		reflect.TypeOf(testLabel{}),
	}
)

func TestGetHeader_struct(t *testing.T) {
	h, err := getHeader(&l)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exheader) != len(h) {
		t.Error("header not correct lenght, got:", len(h),
			"expected:", len(exheader))
		return
	}
	for i, eh := range exheader {
		if eh != h[i] {
			t.Error("header incorrect, got:", eh,
				"expected:", exheader)
			return
		}
	}
}

func TestGetHeader_slice(t *testing.T) {
	h, err := getHeader(&ls)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exheader) != len(h) {
		t.Error("header not correct lenght, got:", len(h),
			"expected:", len(exheader))
		return
	}
	for i, eh := range exheader {
		if eh != h[i] {
			t.Error("header incorrect, got:", eh,
				"expected:", exheader)
			return
		}
	}
}

func TestGetHeader_nested(t *testing.T) {
	typesToPrint = toPrint
	h, err := getHeader(&ln)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exheaderNest) != len(h) {
		t.Error("header not correct lenght, got:", h,
			"expected:", exheaderNest)
		return
	}
	for i, eh := range exheaderNest {
		if eh != h[i] {
			t.Error("header incorrect, got:", eh,
				"expected:", exheaderNest)
			return
		}
	}
}

func TestGetBody_struct(t *testing.T) {
	h, err := getBody(&l)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exbodyStruct) != len(h) {
		t.Error("body not correct lenght, got:", len(h),
			"expected:", len(exheader))
		return
	}
	for i, ehr := range exbodyStruct {
		for j, eh := range ehr {
			if eh != h[i][j] {
				t.Error("header incorrect, got:", eh,
					"expected:", exheader)
				return
			}
		}
	}
}

func TestGetBody_slice(t *testing.T) {
	h, err := getBody(&ls)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exbodySlice) != len(h) {
		t.Error("body row not correct lenght, got:", len(h),
			"expected:", len(exbodySlice))
		return
	}
	for i, ehr := range exbodySlice {
		if len(ehr) != len(h[i]) {
			t.Error("body column not correct lenght, got:", len(h[i]),
				"expected:", len(ehr))
			return
		}
		for j, eh := range ehr {
			if eh != h[i][j] {
				t.Error("body incorrect, got:", eh,
					"expected:", h[i][j])
				return
			}
		}
	}
}

func TestGetBody_nested(t *testing.T) {
	typesToPrint = toPrint
	h, err := getBody(&ln)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exbodyStruct) != len(h) {
		t.Error("body not correct lenght, got:", len(h),
			"expected:", len(exheader))
		return
	}
	for i, ehr := range exbodyStruct {
		for j, eh := range ehr {
			if eh != h[i][j] {
				t.Error("header incorrect, got:", eh,
					"expected:", exheader)
				return
			}
		}
	}
}

func TestGetFooter_struct(t *testing.T) {
	h, err := getFooter(l)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exfooter) != len(h) {
		t.Error("header not correct lenght, got:", len(h),
			"expected:", len(exfooter))
		return
	}
	for i, eh := range exfooter {
		if eh != h[i] {
			t.Error("header incorrect, got:", eh,
				"expected:", exfooter)
			return
		}
	}
}

func TestGetFooter_slice(t *testing.T) {
	h, err := getFooter(ls)
	if err != nil {
		t.Error("error declared", err)
		return
	}
	if len(exfooter) != len(h) {
		t.Error("header not correct lenght, got:", len(h),
			"expected:", len(exfooter))
		return
	}
	for i, eh := range exfooter {
		if eh != h[i] {
			t.Error("header incorrect, got:", eh,
				"expected:", exfooter)
			return
		}
	}
}

func TestPrintTable(t *testing.T) {
	PrintTable(os.Stdout, &ts, []reflect.Type{
		reflect.TypeOf(testLabeln{}),
	})
}
