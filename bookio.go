package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	ASC = "ASC"
	DSC = "DSC"
)

type Output struct {
	Author string
	Books  []struct {
		Name     string
		Revision int
	}
}

type ISBNstruc struct {
	URL     string `json:"url"`
	Key     string `json:"key"`
	Title   string `json:"title"`
	Authors []struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	} `json:"authors"`
	NumberOfPages int    `json:"number_of_pages"`
	Pagination    string `json:"pagination"`
	ByStatement   string `json:"by_statement"`
	Identifiers   struct {
		Goodreads    []string `json:"goodreads"`
		Librarything []string `json:"librarything"`
		Isbn10       []string `json:"isbn_10"`
		Lccn         []string `json:"lccn"`
		Oclc         []string `json:"oclc"`
		Openlibrary  []string `json:"openlibrary"`
	} `json:"identifiers"`
	Classifications struct {
		LcClassifications []string `json:"lc_classifications"`
		DeweyDecimalClass []string `json:"dewey_decimal_class"`
	} `json:"classifications"`
	Publishers []struct {
		Name string `json:"name"`
	} `json:"publishers"`
	PublishPlaces []struct {
		Name string `json:"name"`
	} `json:"publish_places"`
	PublishDate string `json:"publish_date"`
	Subjects    []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"subjects"`
	SubjectPlaces []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"subject_places"`
	SubjectPeople []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"subject_people"`
	SubjectTimes []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"subject_times"`
	Excerpts []struct {
		Text    string `json:"text"`
		Comment string `json:"comment"`
	} `json:"excerpts"`
	Notes string `json:"notes"`
	Links []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"links"`
	Ebooks []struct {
		PreviewURL   string `json:"preview_url"`
		Availability string `json:"availability"`
		Formats      struct {
		} `json:"formats"`
		BorrowURL  string `json:"borrow_url"`
		Checkedout bool   `json:"checkedout"`
	} `json:"ebooks"`
	Cover struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"cover"`
}

type Books struct {
	Links struct {
		Self   string `json:"self"`
		Author string `json:"author"`
		Next   string `json:"next"`
	} `json:"links"`
	Size    int `json:"size"`
	Entries []struct {
		Title   string `json:"title"`
		Authors []struct {
			Author struct {
				Key string `json:"key"`
			} `json:"author"`
			Type struct {
				Key string `json:"key"`
			} `json:"type"`
		} `json:"authors"`
		Key  string `json:"key"`
		Type struct {
			Key string `json:"key"`
		} `json:"type"`
		LatestRevision int `json:"latest_revision"`
		Revision       int `json:"revision"`
		Created        struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"created"`
		LastModified struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"last_modified"`
		Covers           []int    `json:"covers,omitempty"`
		Description      string   `json:"description,omitempty"`
		FirstPublishDate string   `json:"first_publish_date,omitempty"`
		DeweyNumber      []string `json:"dewey_number,omitempty"`
		Subjects         []string `json:"subjects,omitempty"`
	} `json:"entries"`
}

func listOfBooks(Author string, Name string, SORT string) {
	//https://openlibrary.org/authors/OL26320A/works.json

	//Author := "OL26320A"
	var url = "https://openlibrary.org/authors/" + Author + "/works.json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var book Books
	json.Unmarshal([]byte(body), &book)

	output := Output{
		Author: Name,
		Books: []struct {
			Name     string
			Revision int
		}{},
	}
	lenght := len(book.Entries)
	output.Books = make([]struct {
		Name     string
		Revision int
	}, lenght)

	for i := 0; i < len(book.Entries); i++ {
		fmt.Printf("\n")

		output.Books[i].Name = book.Entries[i].Title
		output.Books[i].Revision = book.Entries[i].Revision

		//fmt.Println("Sorted by revision:", output)

	}

	if SORT == ASC {
		sort.SliceStable(output.Books, func(i, j int) bool {
			return output.Books[i].Revision < output.Books[j].Revision
		})
	} else if SORT == DSC {
		sort.SliceStable(output.Books, func(i, j int) bool {
			return output.Books[i].Revision > output.Books[j].Revision
		})
	} else {
		fmt.Println("Invalid SORT - no sorting")
		return
	}

	y_output := output
	yamlData, err := yaml.Marshal(&y_output)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	fmt.Printf("\n")
	fmt.Println(" --- YAML ---")
	fmt.Println(string(yamlData))
	//var output Output

}

func main() {
	//search by ISBN
	var ISBN string
	var SORT string

	fmt.Println("ISBN: ") //0395193958
	fmt.Scanf("%s\n", &ISBN)

	fmt.Println("Sort? ")
	fmt.Scanf("%s\n", &SORT)

	var url = "http://openlibrary.org/api/books?bibkeys=ISBN:" + ISBN + "&jscmd=data&format=json"
	//fmt.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)

	//nepekny triming
	countISBN := len(ISBN)
	sbJson := sb
	helper := countISBN + 10
	sbJson2 := sbJson[helper:]
	sbJson2 = sbJson2[:len(sbJson2)-1]

	var id ISBNstruc
	json.Unmarshal([]byte(sbJson2), &id)
	fmt.Printf("Book name: %s", id.Title)
	fmt.Printf("\nAutor: %s", id.Authors[0].URL)

	fmt.Printf("\n\n")

	var authorList []string
	for _, element := range id.Authors {

		array := strings.Split(element.URL, "/")

		//fmt.Printf("%d\n", index)
		fmt.Printf("%s\n", array[len(array)-2])
		//fmt.Printf("%s\n", element.Name)
		authorList = append(authorList, array[len(array)-2])
	}

	for i := 0; i < len(authorList); i++ {
		fmt.Printf("Author name: %s", id.Authors[i].Name)
		listOfBooks(authorList[i], id.Authors[i].Name, SORT)
		fmt.Printf("\n")

	}

}
