package main

import (
	"os"
	"fmt"
	"math"
	"io/ioutil"
	"bufio"
	"strings"
	"strconv"
	"text/template"
	"github.com/cheggaaa/pb"
)

/*
Inspiration
https://jan.newmarch.name/go/template/chapter-template.html
http://blog.zmxv.com/2011/09/go-template-examples.html
https://gohugo.io/templates/go-templates/
http://andlabs.lostsig.com/blog/2014/05/26/8/the-go-templates-post
http://blog.davidvassallo.me/2014/08/12/nugget-post-golang-html-templates/
http://shadynasty.biz/blog/2012/07/30/quick-and-clean-in-go/
*/


type CM_Test struct {
	Application string
	Migrated bool
	Financement bool
	Duree int
	Passion bool
	Cotitulaire bool
	Clients []Client
	Tickets []Ticket	
}

type Client struct {
	Id string
	Nom string
	Prenom string
}

type Article struct {
	Id int
	Pv float64
}

type Ticket struct {
	Articles []Article
	Total float64
	PointsFid float64
	CompteurCA float64
}

var fmap = template.FuncMap{
	"addPv": addPvFunc,
	"add" : add,
	"trunc": trunc,
	"for": forFunc,
	"discount": discount,
	"float" : floatFunc,
	"enumStr" : enumStrFunc,
	"enumInt" : enumIntFunc,
	"include": includeFunc,
	"cell": NewCell,
}


func main() {
	cwd, _ := os.Getwd()

    argsWithoutProg := os.Args[1:]
    //fmt.Println(argsWithoutProg)    

  	TplFileName := argsWithoutProg[0]

  	bar := pb.StartNew(1)
  
	var t = template.Must(template.New(TplFileName).Funcs(fmap).ParseFiles( cwd + "/" + TplFileName ))

	sParamFile := cwd + "/param."+TplFileName

	if _, err := os.Stat(cwd + "/param."+TplFileName); os.IsNotExist(err) {
	// path/to/whatever does not exist
		//fmt.Printf("Attention :  %v \n", err)
		sParamFile = cwd + "/param.default.txt"
		//fmt.Printf("Utilisation de la configuration par défaut :  %v \n", sParamFile)
		if _, err := os.Stat(cwd + "/param.default.txt"); os.IsNotExist(err) {
			fmt.Printf("Erreur irrecouvrable :  %v \n", err)
			os.Exit(1)
		}
	}

	inFile, e := os.Open(sParamFile)
	if e != nil {panic(e)}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines) 
	var lofTs []CM_Test
	for scanner.Scan() {	
		var t CM_Test
		//var sLofClient string
		//var sLofTickets string
		s := strings.Split(scanner.Text(), ";")	

		t.Application 	= s[0]
		t.Migrated 		= Str2Bool(s[1],"true")
		t.Financement 	= Str2Bool(s[2],"true")
		t.Duree,_		= strconv.Atoi(s[3])
		t.Passion 		= Str2Bool(s[4],"true")
		t.Cotitulaire 	= Str2Bool(s[5],"true")

		inFile, e := os.Open(s[6])
		if e != nil {panic(e)}
		defer inFile.Close()    		
		scanner := bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines) 
        //fmt.Println(f)
        var lofC []Client
        i:=0
        for scanner.Scan() {
        	if i > 0 && i < 6  {	// as i don't want more than 5 clients and not the first line
			var c Client
			s := strings.Split(scanner.Text(), ",")
			c.Id = s[1]
			c.Nom = s[2]
			c.Prenom = s[3]

			//fmt.Printf("client  : %v\n", c )
			lofC = append(lofC, c )//.String())
			//fmt.Printf("lofC %v\n", lofC )	
			}
			i++			
		}

		t.Clients = lofC

		inFile, e = os.Open(cwd + "\\" + s[7])

		// attention parfois il y a des espaces en fin de nom de fichier 

		if e != nil {panic(e)}
		defer inFile.Close()
		scanner = bufio.NewScanner(inFile)
		scanner.Split(bufio.ScanLines) 
		var lofT []Ticket
		for scanner.Scan() {	
			var t Ticket
			s := strings.Split(scanner.Text(), ";")	
				inFile2, e2 := os.Open(cwd + "/"+s[0] )
				if e2 != nil {panic(e2)}
				defer inFile2.Close()
				scanner2 := bufio.NewScanner(inFile2)	
				scanner2.Split(bufio.ScanLines) 
				var lofA []Article
				for scanner2.Scan() {	
					var a Article
					s := strings.Split(scanner2.Text(), ";")
					a.Id, _ = strconv.Atoi(s[0])
					a.Pv, _ = strconv.ParseFloat(s[1], 64)
					lofA = append(lofA, a )//.String())
					//fmt.Printf("lofA %v\n", lofA )
				}
			t.Articles = lofA
			t.Total, _ = strconv.ParseFloat(s[1], 64)
			t.PointsFid, _ = strconv.ParseFloat(s[2], 64)
			t.CompteurCA, _ = strconv.ParseFloat(s[3], 64)
			lofT = append(lofT, t )//.String())
			//fmt.Printf("lofT %v\n", lofT )
		}	

		t.Tickets = lofT

		lofTs = append(lofTs, t )//.String())
		//fmt.Printf("lofTs %v\n", lofTs )	
	}

	f, err := os.Create(cwd + "/"+TplFileName + "-compiled.txt")
    if err != nil {panic(err)}
	defer f.Close()
	err = t.Execute(f, lofTs)
	if err != nil {
		fmt.Printf("exec error: %v\n", err)
	}	
    bar.Increment()
    bar.FinishPrint("Done ! .... result in " + cwd + "\\"+TplFileName + "-compiled.txt")
}

func Str2Bool(a, b string) bool{
    if a == b {
        return true
    }
    return false
}

/* ne peut pas fonctionner meme {{$tot:= addPv $article.Pv $tot }} met toujours 5 + la valeur de tot a jour */
func addPvFunc(a float64, b float64) float64 {
	return a+b
}


func add(a float64, b float64) float64 {
	return a+b
}

func trunc(a float64) float64 {
	return math.Trunc(a)
}

func forFunc(start, end, step int) []int {
	//fmt.Printf("%v\n", ( end - start ) / step  )
	iCnt := int ( ( end - start ) / step ) + 1 // int() pas obligatoire
	//iCnt := int(math.Trunc( ( end - start ) / step ) ) +1
	s := make([]int, iCnt)
	for i, _:= range s {
		s[i] = start + step * i
	}
	
	return s
}

/*
* tx in percent ie 15 = 15%
*/
func discount(val float64, tx float64) float64 {
	return val - (val * ( tx / 100) ) 
}

func floatFunc(i int) float64 {
	return float64(i)
}

func enumStrFunc(what ...string) [] string{
	//fmt.Printf("%v\n", what )
	return what
}
func enumIntFunc(what ...int) [] int{
	//fmt.Printf("%v\n", what )
	return what
}

func includeFunc(in string, relative bool) string {
	var f []byte
	var err error
	if relative == false {
		f, err = ioutil.ReadFile(in) 
	} else {
			cwd, _ := os.Getwd()
		f, err = ioutil.ReadFile(fmt.Sprintf("%s%s%s",cwd,string(os.PathSeparator),in))
	}
	if err != nil {panic(err)}
	return fmt.Sprintf("%s",f)
}

/* 
http://play.golang.org/p/SufZdsx-1v pour trouver peut etre une solution 
issu de la discussion https://groups.google.com/forum/#!topic/golang-nuts/MUzNZ9dbHrg
*/

/*
Pour faire le total dans une feuille ....
{{range $idx, $ticket := $.Tickets }}{{$fTotal := cell 0.0}}{{range $article := $ticket.Articles }}{{ $sum := add $fTotal.Get $article.Pv }}{{$_ := $fTotal.Set $sum }}{{end}}Total :{{$fTotal.Get}} {{end}}
*/

type Cell struct{ v interface{} }


func NewCell(v ...interface{}) (*Cell, error) {
	switch len(v) {
	case 0:
		return new(Cell), nil
	case 1:
		return &Cell{v[0]}, nil
	default:
		return nil, fmt.Errorf(" wrong number of args: want 0 or 1, got %v", len(v))
	}
}
func (c *Cell) Set(v interface{}) *Cell { c.v = v; return c }
func (c *Cell) Get() interface{}        { return c.v }