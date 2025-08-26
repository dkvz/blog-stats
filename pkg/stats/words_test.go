package stats

import "testing"

func TestInlineTagsRemoval(t *testing.T) {
	sut := `<p>some text</p>
		<p>More text <a href="/testing" class="some">link text</a></p>`

	want := 6
	wc := WordCount(sut)
	if wc != want {
		t.Errorf("TestWordCountInline: count is %v, should be %v", wc, want)
	}
}

func TestRemainingTagsRemoval(t *testing.T) {
	sut := `<h1>Title of the post</h1>
		<p>Code extract:</p>
		<pre class="screen"><code class="language-javascript">
		console.log("hello")
		</pre></code>
		<p>More text</p>
		<img src="img.png" alt="some image">`

	want := 8
	wc := WordCount(sut)
	if wc != want {
		t.Errorf("TestWordCountInline: count is %v, should be %v", wc, want)
	}
}
