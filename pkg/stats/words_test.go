package stats

import "testing"

func TestWordCountInline(t *testing.T) {
	sut := `<p>some text</p>
		<p>More text <a href="/testing" class="some">link text</a></p>`
	want := 6
	wc := WordCount(sut)
	if wc != want {
		t.Errorf("TestWordCountInline: count is %v, should be %v", wc, want)
	}
}
