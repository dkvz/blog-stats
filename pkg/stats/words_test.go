package stats

import "testing"

func TestInlineTagsRemoval(t *testing.T) {
	sut := `<p>some text</p>
		<p>More text <a href="/testing" class="some">link text</a></p>`

	want := 6
	wc := WordCount(&sut)
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
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestRemainingTagsRemoval: count is %v, should be %v", wc, want)
	}
}

func TestBigTag(t *testing.T) {
	sut := `<h2>Title <b>with html in it</b></h2>

	<p>Some text&nbsp;here</p>

	<div class="card-panel z-depth-3 article-image center-image" style="max-width: 1000px">
<a href="/wp-content/blog.png" target="_blank"><img src="/wp-content/blog.png" alt="Some image" class="responsive-img"></a>
<div class="image-legend">Image legends are currently well counted in</div>
</div>`

	want := 15
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestBigTag: count is %v, should be %v", wc, want)
	}
}

func TestUntaggedText(t *testing.T) {
	sut := `text out of any tag <i>italic</i> word
	<p>Some text</p>
	untagged again
	</div>
`

	want := 11
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestUntaggedText: count is %v, should be %v", wc, want)
	}
}

func TestBetweenBrs(t *testing.T) {
	sut := `text out of any tag word<br /><br />
	Some text
	<br /><br />
	untagged again `

	want := 10
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestBetweenBrs: count is %v, should be %v", wc, want)
	}
}

func TestHtmlComments(t *testing.T) {
	sut := `<!-- some comment here -->
	<!-- Another comment -->
	<p>Article text</p>
	<p>Article text again</p>
	<!-- Ends with a comment -->`

	want := 5
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestHtmlComments: count is %v, should be %v", wc, want)
	}
}

func TestBrAtTheEnd(t *testing.T) {
	sut := `<h1>Title</h1><!-- some comment here -->
	<p>Article text</p>
	<br>
	<p>Article text again</p>
	<br>`

	want := 6
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestBrAtTheEnd: count is %v, should be %v", wc, want)
	}
}

func TestImgAtBothEnds(t *testing.T) {
	sut := `<h1>Title</h1>
	<div class="card-panel z-depth-3 article-image center-image" style="max-width: 1000px">
	<a href="/wp-content/blog.png" target="_blank"><img src="/wp-content/blog.png" alt="Some image" class="responsive-img"></a>
	<div class="image-legend">Image legends are currently well counted in</div>
	</div>
	<p>Article text</p>
	<p>Article text again</p>
	<div class="card-panel z-depth-3 article-image center-image" style="max-width: 1000px">
	<a href="/wp-content/blog.png" target="_blank"><img src="/wp-content/blog.png" alt="Some image" class="responsive-img"></a>
	<div class="image-legend">Image legends are currently well counted in</div>
	</div>
	`

	want := 20
	wc := WordCount(&sut)
	if wc != want {
		t.Errorf("TestImgAtBothEnds: count is %v, should be %v", wc, want)
	}
}
