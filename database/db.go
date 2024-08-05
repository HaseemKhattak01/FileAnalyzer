package database

import (
	"FileReader/models"
	"fmt"
)

func CreateUser(vowels, spaces, capitalletters, smallleters, Words int) error {
	fmt.Println("In create user function")
	fmt.Println(vowels, spaces, capitalletters, smallleters, Words)

	_, err := db.Exec("INSERT INTO filerecords (vowels, spaces, capitalletters, smallleters, words) VALUES($1, $2, $3, $4, $5)", vowels, spaces, capitalletters, smallleters, Words)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

func Update(up models.UpdateField) {
	var query string
	if up.Field == "vowels" {
		query = fmt.Sprintf("UPDATE public.filerecords SET vowels=$1 WHERE id=$2")
	} else if up.Field == "capitalletters" {
		query = fmt.Sprintf("UPDATE public.filerecords SET capitalletters=$1 WHERE id=$2")
	} else if up.Field == "smallleters" {
		query = fmt.Sprintf("UPDATE public.filerecords SET smallleters=$1 WHERE id=$2")
	} else if up.Field == "words" {
		query = fmt.Sprintf("UPDATE public.filerecords SET words=$1 WHERE id=$2")
	} else if up.Field == "spaces" {
		query = fmt.Sprintf("UPDATE public.filerecords SET spaces=$1 WHERE id=$2")
	}

	row := db.QueryRow(query, up.Value, up.Id)
	fmt.Println(row)
}

func DeleteRecords(id int) (*models.DBResults, error) {
	var cal models.DBResults
	query := "DELETE FROM filerecords WHERE id =$1"
	result, err := db.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(result)
	return &cal, nil
}

func Getdata() ([]models.Results, error) {
	rows, err := db.Query("SELECT * FROM filerecords")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Results

	for rows.Next() {
		var stat models.Results
		if err := rows.Scan(&stat.Id, &stat.Vowels, &stat.Capitalletters, &stat.Smallletters, &stat.Spaces, &stat.Words); err != nil {
			return nil, err
		}
		results = append(results, stat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil

}
