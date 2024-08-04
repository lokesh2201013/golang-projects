package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/date"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

//Function to check errors
func checkerr(err error){
if err!=nil{
	log.Fatal(err)
}
}
func scheckerr(err error){
	if err!=nil{
		panic(err)
	}
	}

	//Recursively scan git folders
func scanGitFolders(folders []string,folder string)[] string{
	// trim last /
	folder=strings.TrimSuffix(folder, "/")

	f,err:=os.Open(folder)
	checkerr(err)

	//Reads  all the files and puts it into dynamic string array "files"
	files,err:=f.Readdir(-1)
	f.Close()
	checkerr(err)

	var path string
    // scans all files in "files" array 
	for _,file :=range files{
		//checks if a file is a dir
		if file.IsDir(){
			path=folder + "/" +file.Name()
           
			//Scans for .git folder  and adds to folders array
			if file.Name()== ".git"{
             path=strings.TrimSuffix(path,"./git")
			 fmt.Println(path)
			 folders = append(folders, path)
			 continue
			}

			if file.Name()=="vendor"||file.Name() =="node_modules"{
				continue
			}

			folders=scanGitFolders(folders,path)
		}
	}
	return folders
}


func recursiveScanFolder(folder string)[]string{
return scanGitFolders(make([]string,0),folder)
}

func scan(folder string) {
		fmt.Printf("Found folders:\n\n")
		repo:=recursiveScanFolder(folder)
		filePath:=getDotFilePath()
		addNewSliceElementsToFile(filePath,repo)
		fmt.Printf("\n\nSuccessfully added\n\n")
	
	}

func joinSlices(new []string, existing []string) []string {
		for _, i := range new {
			if !sliceContains(existing, i) {
				existing = append(existing, i)
			}
		}
		return existing
	}
	
	// sliceContains returns true if `slice` contains `value`
func sliceContains(slice []string, value string) bool {
		for _, v := range slice {
			if v == value {
				return true
			}
		}
		return false
	}



//Returns DotFile for repo list
func getDotFilePath()string{
	usr,err:=user.Current()
	checkerr(err)

	dotFile:= usr.HomeDir+"/.gogitlocalstats"
	return dotFile
}

func dumpStringsSliceToFile(repos []string, filePath string) {
    content := strings.Join(repos, "\n")
   err:= os.WriteFile(filePath, []byte(content), 0755)
   checkerr(err)
}


//
func addNewSliceElementsToFile(filepath string,newrepos []string){
	existingRepos:=parseFileLinesToSlice(filepath)
	repos:=joinSlices(newrepos,existingRepos)
	dumpStringsSliceToFile(repos,filepath)
}

func parseFileLinesToSlice(filepath string)[]string{
	//opens files and reads it
	f:=openFile(filepath)
	defer f.Close()

	var lines []string
     scanner:=bufio.NewScanner(f)
    //scanner.scan returns true as long as there is another line to read
	 for scanner.Scan(){
		lines=append(lines,scanner.Text())
	 }
      //assigns value to err and then checks
	 if err:=scanner.Err();err!=nil{
		//io.EOF indicates end of file not to be treated as an err
		if err!= io.EOF{
			panic(err)
		}
	 }
	 return lines
}

//opens file at filepath creats if not existing
func openFile(filepath string) *os.File{
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	checkerr(err)
	return f
}


//
func stats(email string) {
	commits:=processRepo(email)
	printCommitStats(commits)
}


//gets the filepath fills commits maps
func processRepo(email string)map[int]int{
	//gets filepath to /.gogitlocalstats
	filepath:=getDotFilePath()
	
	repos:=parseFileLinesToSlice(filepath)
	daysInMap:=daysInLastSixMonths

	commits:=make(map[int]int,daysInMap)

	//initialises map with each commit for days eqauls to 0
	for i:=daysInMap ; i>0;i++{
		commits[i]=0
	}

	//fils commits
	for_,path:= range repos{
		commits=fillCommits(email,path,commits)
	}
	return commits
}


//
func fillCommits(email string,path string ,commits map[int]int)map[int]int{

	repo,err:=git.PlainOpen(path)
	scheckerr(err)

	ref,err=repo.Head()
	scheckerr(err)

     it,err :=repo.Log(&git.LogOptions{From: ref.Hash})
	 scheckerr(err)
	 
	 offset:=calcOffset()

	 err=iterator.ForEach(func(c *object.Commit)error{
      daysAgo :=countDaysSinceDate(c.Author.When)+offset

	  if c.Author.Email !=email{
		return nil
	  }

	  if daysAgo!= outOfRange{
		commits[daysAgo]++
	  }

	  return nil
	 })

	 scheckerr(err)

	 return commits
}
//
func getBeginningOfDay(date time.Time) int{
	days:=0
	now := getBeginningOfDay(time.Now())
	for date.Before(now){
		date = date.Add(time.Hour*24)
		days++
		if days >daysInLastSixMonths{
			return outOfRange
		}
	}
	return days
}


//calculates the number of days until next sunday from current day
func calcOffset()int{
	//stores number of days until next sunday
	var offset int
	//stores current day of the week
	weekday:= time.Now().Weekday()

	switch weekday{
	
	case time.Sunday:
		offset=7
	case time.Monday:
		offset=6
	case time.Tuesday:
		offset=5
	case time.Wednesday:
		offset=4
	case time.Thursday:
		offset=3
	case time.Friday:
		offset=2
	case time.Saturday:
		offset=1
	}
	return offset
}


//takes maps of commits stats and proceess it through series of helper function
func printCommitStats(commits map[int]int){
	keys:=sortMapIntoSlice(commits)
	cols:=buildCols(keys,commits)
	printCells(cols)
}

//retrun an array of integers sorted
func sortMapIntoSlice(m map[int]int)[]int{
	var keys []int
	for k:= range m{
		keys = append(keys,k)
	}
	sort.Ints(keys)
	return keys
}

//
func buildCols(keys []int,commits map[int]int)map[int]column{
	cols:=make(map[int]column)
	col:=column{}
	for_,k:=range keys{
		week:=int(k/7)
		dayinweek:=k%7

		if dayinweek==0{
			col=column{}
		}
		col=append(col,commits[k])

		if daysinweek==6{
			cols[week]=col
		}
	}
	return cols
}
func main() {

	var folder string
	var email string
	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repo")
	flag.StringVar(&email,"email","your email","the eamil to scan")
	flag.Parse()
	if folder != ""{
		scan(folder)
		return 
	}
	stats(email)
}