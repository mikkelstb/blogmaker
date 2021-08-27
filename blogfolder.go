package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type blogFolder struct {
	path string
}

/* getYears() reads the root catalouge and returns a slice of yearfolders */
func (bf *blogFolder) getYears() []string {
	dirs, err := ioutil.ReadDir(bf.path)
	if err != nil {
		log.Fatal(err)
	}

	year_pattern := regexp.MustCompile(`^\d{4}$`)

	var folders []string
	for _, folder := range dirs {
		if folder.IsDir() && year_pattern.MatchString(folder.Name()) {
			fmt.Println(folder.Name())
			folders = append(folders, folder.Name())
		}
	}
	return folders
}

func (bf *blogFolder) createYearFolder(year int) error {

	folderpath := fmt.Sprintf("./%v/%d", bf.path, year)
	if _, err := os.Stat(folderpath); os.IsNotExist(err) {

		err = os.Mkdir(folderpath, 0770)
		if err != nil {
			return error(err)
		}
	}
	return nil
}

/* writeEntry(Entry) Writes an entry to the file system
The filename is generated according to date and time.
And written into subfolder of current year
Ex: ./2021/2021_08_23_1523.xml */
func (bf *blogFolder) writeEntry(e Entry) {

	year, month, day := e.Created.Date()
	hour, min, _ := e.Created.Clock()

	filename := fmt.Sprintf("%v/%v/%v_%.2v_%.2v_%.2v%.2v.xml", bf.path, year, year, int(month), day, hour, min)
	xml_data, err := xml.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	err = bf.createYearFolder(year)
	if err == nil {
		err = os.WriteFile(filename, xml_data, 0600)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (bf *blogFolder) getAllEntries() []Entry {

	fmt.Println("OK!")

	var entries []Entry

	for _, year := range bf.getYears() {

		files, err := ioutil.ReadDir(fmt.Sprintf("%v/%v", bf.path, year))
		if err != nil {
			log.Fatal(err)
		}

		for _, filename := range files {
			file, err := ioutil.ReadFile(fmt.Sprintf("%v/%v/%v", bf.path, year, filename.Name()))
			if err != nil {
				log.Fatal(err)
			}
			e := Entry{}
			err = xml.Unmarshal(file, &e)
			if err != nil {
				log.Fatal(err)
			}
			entries = append(entries, e)
		}
	}
	var ea []Entry
	for index := len(entries) - 1; index >= 0; index-- {
		ea = append(ea, entries[index])
	}

	return ea
}

func (bf *blogFolder) getEntries(limit, offset int) []Entry {

	entries := bf.getAllEntries()
	if len(entries) <= offset {
		return []Entry{}
	}
	if len(entries) <= limit+offset {
		return entries[offset:]
	}
	return entries[offset : limit+offset]
}
