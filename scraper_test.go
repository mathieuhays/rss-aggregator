package rss_aggregator

import (
	"testing"
)

func TestExtractRSSItems(t *testing.T) {
	source := `
<?xml version="1.0" encoding="utf-8" standalone="yes"?>
<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom">
  <channel>
    <title>Boot.dev Blog</title>
    <link>https://blog.boot.dev/</link>
    <description>Recent content on Boot.dev Blog</description>
    <generator>Hugo</generator>
    <language>en-us</language>
    <lastBuildDate>Wed, 10 Jul 2024 00:00:00 +0000</lastBuildDate>
    <atom:link href="https://blog.boot.dev/index.xml" rel="self" type="application/rss+xml" />
    <item>
      <title>The Boot.dev Beat. July 2024</title>
      <link>https://blog.boot.dev/news/bootdev-beat-2024-07/</link>
      <pubDate>Wed, 10 Jul 2024 00:00:00 +0000</pubDate>
      <guid>https://blog.boot.dev/news/bootdev-beat-2024-07/</guid>
      <description>One million lessons. Well, to be precise, you have all completed 1,122,050 lessons just in June.</description>
    </item>
    <item>
      <title>The Boot.dev Beat. June 2024</title>
      <link>https://blog.boot.dev/news/bootdev-beat-2024-06/</link>
      <pubDate>Wed, 05 Jun 2024 00:00:00 +0000</pubDate>
      <guid>https://blog.boot.dev/news/bootdev-beat-2024-06/</guid>
      <description>ThePrimeagen&amp;rsquo;s new Git course is live. A new boss battle is on the horizon, and we&amp;rsquo;ve made massive speed improvements to the site.</description>
    </item>
</channel>
</rss>`

	items, err := extractRSSItems(source)

	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 items got %d instead", len(items))
	}

	expectedFirstTitle := "The Boot.dev Beat. July 2024"
	if items[0].Title != expectedFirstTitle {
		t.Errorf("invalid decoding. Expected title %q, got %q", expectedFirstTitle, items[0].Title)
	}
}
