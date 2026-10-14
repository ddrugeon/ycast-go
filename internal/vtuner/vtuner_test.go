package vtuner

import (
	"fmt"
	"testing"
)

func TestXMLMatchesPython(t *testing.T) {
	// Golden strings produced by the original vtuner.py.
	const goldenDirs = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?><ListOfItems><ItemCount>2</ItemCount><Item><ItemType>Dir</ItemType><Title>Genres</Title><UrlDir>http://denon.vtuner.com/ycast/radiobrowser/genre/?vtuner=true</UrlDir><UrlDirBackUp>http://denon.vtuner.com/ycast/radiobrowser/genre/?vtuner=true</UrlDirBackUp><DirCount>12</DirCount></Item><Item><ItemType>Dir</ItemType><Title>R&amp;B &lt;Soul&gt;</Title><UrlDir>http://denon.vtuner.com/ycast/my_stations/R%26B%20%3CSoul%3E?vtuner=true</UrlDir><UrlDirBackUp>http://denon.vtuner.com/ycast/my_stations/R%26B%20%3CSoul%3E?vtuner=true</UrlDirBackUp><DirCount>3</DirCount></Item></ListOfItems>`
	const goldenDisplay = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?><ListOfItems><ItemCount>1</ItemCount><Item><ItemType>Display</ItemType><Display>No stations found</Display></Item></ListOfItems>`

	p := &Page{Count: 2}
	p.Add(DirItem{Title: "Genres", URL: "http://denon.vtuner.com/ycast/radiobrowser/genre/", Count: 12})
	p.Add(DirItem{Title: "R&B <Soul>", URL: "http://denon.vtuner.com/ycast/my_stations/R%26B%20%3CSoul%3E", Count: 3})
	if got := p.String(); got != goldenDirs {
		t.Errorf("dirs XML differs:\n got: %s\nwant: %s", got, goldenDirs)
	}
	if got := MessagePage("No stations found").String(); got != goldenDisplay {
		t.Errorf("display XML differs:\n got: %s\nwant: %s", got, goldenDisplay)
	}
}

func TestXMLEscapingAndControlChars(t *testing.T) {
	got := xmlText("a & b <c> \x00\x07 d\u00e9 \xff")
	if got != "a &amp; b &lt;c&gt;  d\u00e9 " {
		t.Errorf("xmlText = %q", fmt.Sprintf("%q", got))
	}
}
