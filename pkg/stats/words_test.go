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
		t.Errorf("TestRemainingTagsRemoval: count is %v, should be %v", wc, want)
	}
}

func TestBigTag(t *testing.T) {
	sut := `<h2>Title <b>with html in it</b></h2>

	<p>Some text&nbsp;here</p>

	<div class="card-panel z-depth-3 article-image center-image" style="max-width: 1000px">
<a href="/wp-content/blog.png" target="_blank"><img src="/wp-content/blog.png" alt="Some image" class="responsive-img"></a>
<div class="image-legend">Image legends are currently not counted in</div>
</div>`

	want := 8
	wc := WordCount(sut)
	if wc != want {
		t.Errorf("TestBigTag: count is %v, should be %v", wc, want)
	}
}
