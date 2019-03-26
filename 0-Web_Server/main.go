package main
import (
"fmt"
"net/http"
"sort"
"log"
"github.com/RealJK/rss-parser-go"
)

type rssItems []rss.Item

func (items rssItems) Less(i, j int) bool {
 return items[i].PubDate < items[j].PubDate 
}

func (items rssItems) Len() int {
 return len(items)
}

func (items rssItems) Swap(i, j int) { items[i], items[j] = items[j], items[i] }

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)

	rssObject, err_1 := rss.ParseRSS("http://blagnews.ru/rss_vk.xml")
	rssObject_2, err_2 := rss.ParseRSS("https://lenta.ru/rss")
	rssObject_3, err_3 := rss.ParseRSS("https://news.mail.ru/rss/90/")

	body := ""

	if err_1 != nil && err_2 != nil && err_3 != nil  {
		body += fmt.Sprintf("<p> Parsed RSS channels:" +
			"<ul><li><a href=\"/first_rss\">%s</a></li>"+
			"<li><a href=\"/sec_rss\">%s</a></li>"+
			"<li><a href=\"/third_rss\">%s</a></li></ul></p>", 
			rssObject.Channel.Title, rssObject_2.Channel.Title, rssObject_3.Channel.Title)
		rss_Items := append(rssObject.Channel.Items, rssObject_2.Channel.Items...)
		rss_Items = append(rss_Items, rssObject_3.Channel.Items...)

		items := rssItems(rss_Items)
		sort.Sort(items)

		for v := range items {
			item := items[v]
			body += fmt.Sprintf("<p> <b>%s</b>" +
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div style=\"background-color: #dffcff\">%s</div></p>", 
			item.Title, item.Description, item.Guid.Value, item.Link, item.PubDate)
		}
	}

	fmt.Fprintf(w,  "<h1>Main page</h1>" +
		"<div><i>GET: </i>%s</div><div>%s</div>", r.URL.RawQuery, body)
}

func FirstRSSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)

	rssObject, err_1 := rss.ParseRSS("http://blagnews.ru/rss_vk.xml")

	body := ""

	if err_1 != nil  {
		body += fmt.Sprintf("<a href=\"/\">Main</a>" +
			"<p>Current RSS:</p><h1>%s</h1>", 
			rssObject.Channel.Title)
		items := rssObject.Channel.Items

		for v := range items {
			item := items[v]
			body += fmt.Sprintf("<p> <b>%s</b>" +
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div style=\"background-color: #dffcff\">%s</div></p>", 
			item.Title, item.Description, item.Guid.Value, item.Link, item.PubDate)
		}
	}

	fmt.Fprintf(w, "<div>%s</div>", body)
}

func SecondRSSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)

	rssObject, err_1 := rss.ParseRSS("https://lenta.ru/rss")

	body := ""

	if err_1 != nil  {
		body += fmt.Sprintf("<a href=\"/\">Main</a>" +
			"<p>Current RSS:</p><h1>%s</h1>", 
			rssObject.Channel.Title)
		items := rssObject.Channel.Items

		for v := range items {
			item := items[v]
			body += fmt.Sprintf("<p> <b>%s</b>" +
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div style=\"background-color: #dffcff\">%s</div></p>", 
			item.Title, item.Description, item.Guid.Value, item.Link, item.PubDate)
		}
	}

	fmt.Fprintf(w, "<div>%s</div>", body)
}

func ThirdRSSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path", r.URL.Path)

	rssObject, err_1 := rss.ParseRSS("https://news.mail.ru/rss/90/")

	body := ""

	if err_1 != nil  {
		body += fmt.Sprintf("<a href=\"/\">Main</a>" +
			"<p>Current RSS:</p><h1>%s</h1>", 
			rssObject.Channel.Title)
		items := rssObject.Channel.Items

		for v := range items {
			item := items[v]
			body += fmt.Sprintf("<p> <b>%s</b>" +
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div>%s</div>"+
			"<div style=\"background-color: #dffcff\">%s</div></p>", 
			item.Title, item.Description, item.Guid.Value, item.Link, item.PubDate)
		}
	}

	fmt.Fprintf(w, "<div>%s</div>", body)
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Menu</h1>"+
    	"<a href=\"/\"=/>Main</a><br>"+
        "<form>"+
        "Name:<br>"+
        "<input type=\"text\" value=\"Save\">"+
        "</form>")
}

func MenuHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Menu</h1>"+
    	"<a href=\"/form\"=/>Form</a><br>"+
        "<iframe  width = \"320\" height = \"240\" src = \"https://www.youtube.com/embed/d9TpRfDdyU0?autoplay=1?loop=1&start=28&color=white\"></iframe>")
}

func main() {
	http.HandleFunc("/", HomeRouterHandler)
	http.HandleFunc("/form/", FormHandler)
	http.HandleFunc("/menu/", MenuHandler)
	http.HandleFunc("/first_rss/", FirstRSSHandler)
	http.HandleFunc("/sec_rss/", SecondRSSHandler)
	http.HandleFunc("/third_rss/", ThirdRSSHandler)



	err := http.ListenAndServe(":9016", nil) // задаем слушать порт

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
