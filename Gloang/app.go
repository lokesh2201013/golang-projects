package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"io"
	"os/user"
	"strings"
	"bufio"
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

//opens file at filepath crreats if not exixting
func openFile(filepath string) *os.File{
	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	checkerr(err)
	return f
}



func stats(email string) {
	print("STATS")
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