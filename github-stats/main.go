package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type langCount struct {
	lang  string
	count int
}

type yearSort struct {
	year  int
	count int
}

func sortByYear(ya YearActivity) []yearSort {
	var sortedEntries []yearSort
	for year, count := range ya {
		sortedEntries = append(sortedEntries, yearSort{year, count})
	}

	sort.Slice(sortedEntries, func(i, j int) bool {
		return sortedEntries[i].year < sortedEntries[j].year
	})

	return sortedEntries
}

type UserTable struct {
	users []UserData
	app   *tview.Application
	table *tview.Table
}

type UserData struct {
	Username     string
	Name         string
	Followers    int
	Languages    Languages
	TotalForks   int
	ReposCreated YearActivity
	ReposUpdated YearActivity
}

func NewUserTable(users []UserData) *UserTable {
	app := tview.NewApplication()
	table := tview.NewTable().SetBorders(true)
	headers := []string{"Username", "Name", "Followers", "Languages", "Total Forks", "Repos Created", "Repos Updated"}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(tview.Styles.SecondaryTextColor).
			SetAlign(tview.AlignCenter).
			SetMaxWidth(20).
			SetSelectable(false)
		table.SetCell(0, i, cell)
	}
	return &UserTable{users: users, app: app, table: table}
}

func (ut *UserTable) RenderTable() error {
	row := 1
	for _, user := range ut.users {
		var langCounts []langCount
		for lang, count := range user.Languages{
			langCounts = append(langCounts, langCount{lang, count})
		}
		sort.Slice(langCounts, func(i, j int) bool {
			return langCounts[i].count > langCounts[j].count
		})
		var resultLangs []string
		for _, lc := range langCounts {
			resultLangs = append(resultLangs, lc.lang)
		}
		langsStr := strings.Join(resultLangs, ", ")

		var creationStr, updateStr string
		sortedCreationCounts := sortByYear(user.ReposCreated)
		for i, entry := range sortedCreationCounts {
			if i > 0 {
				creationStr += ", "
			}
			creationStr += fmt.Sprintf("%d: %d", entry.year, entry.count)
		}
		sortedUpdateCounts := sortByYear(user.ReposUpdated)
		for i, entry := range sortedUpdateCounts {
			if i > 0 {
				updateStr += ", "
			}
			updateStr += fmt.Sprintf("%d: %d", entry.year, entry.count)
		}

		ut.table.SetCell(row, 0, tview.NewTableCell(user.Username).SetAlign(tview.AlignCenter).SetMaxWidth(15).SetExpansion(1))
		ut.table.SetCell(row, 1, tview.NewTableCell(user.Name).SetAlign(tview.AlignCenter).SetMaxWidth(15).SetExpansion(1))
		ut.table.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("%d", user.Followers)).SetAlign(tview.AlignCenter).SetMaxWidth(10).SetExpansion(1))
		ut.table.SetCell(row, 3, tview.NewTableCell(langsStr).SetAlign(tview.AlignLeft).SetMaxWidth(20).SetExpansion(1))
		ut.table.SetCell(row, 4, tview.NewTableCell(fmt.Sprintf("%d", user.TotalForks)).SetAlign(tview.AlignCenter).SetMaxWidth(10).SetExpansion(1))
		ut.table.SetCell(row, 5, tview.NewTableCell(creationStr).SetAlign(tview.AlignLeft).SetMaxWidth(20).SetExpansion(1))
		ut.table.SetCell(row, 6, tview.NewTableCell(updateStr).SetAlign(tview.AlignLeft).SetMaxWidth(20).SetExpansion(1))

		row++
	}

	if err := ut.app.SetRoot(ut.table, true).Run(); err != nil {
		return fmt.Errorf("error running table app: %v", err)
	}
	return nil
}

func fetchUserData(usernames []string) ([]UserData, error) {
	users := make([]UserData, 0)
	for _, username := range usernames {
		resp, err := fetchGitHub(fmt.Sprintf("https://api.github.com/users/%s", username))
		if err != nil {
			fmt.Printf("Error fetching user data for %s: %v\n", username, err)
			continue
		}

		var user User
		if err := json.Unmarshal(resp, &user); err != nil {
			fmt.Printf("Error decoding user data for %s: %v\n", username, err)
			continue
		}

		repos, err := user.Repos()
		if err != nil {
			fmt.Printf("Error loading repos data for %s: %v\n", username, err)
			continue
		}

		totalLangs := make(Languages)
		totalForks := 0
		creationCounts := make(YearActivity)
		updateCounts := make(YearActivity)
		for _, repo := range repos {
			langs, err := repo.Langs()
			if err != nil {
				fmt.Printf("Error loading languages data for %s: %v\n", repo.Name, err)
				continue
			}
			for lang, count := range langs {
				totalLangs[lang] += count
			}
			totalForks += repo.ForksCount
			creationYear, err := getYear(repo.CreatedAt)
			if err == nil {
				creationCounts[creationYear]++
			}
			updateYear, err := getYear(repo.UpdatedAt)
			if err == nil {
				updateCounts[updateYear]++
			}
		}

		users = append(users, UserData{
			Username:     user.Login,
			Name:         user.Name,
			Followers:    user.Followers,
			Languages:    totalLangs,
			TotalForks:   totalForks,
			ReposCreated: creationCounts,
			ReposUpdated: updateCounts,
		})
	}
	return users, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if len(os.Args) != 2 {
		panic("provide 1 arg, which should be a file")
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic("Failed to open file")
	}
	defer file.Close()

	usernames := make([]string, 0, 100)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		usernames = append(usernames, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	startTime := time.Now()
	users, err := fetchUserData(usernames)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch user data: %v", err))
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Time taken to load data: %v\n", elapsedTime)

	userTable := NewUserTable(users)
	if err := userTable.RenderTable(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
