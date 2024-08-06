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

func SignUp_db(user models.Identity) (bool, error) {
	isuser := false
	query := "SELECT 'exists' AS result FROM userid WHERE username = ? UNION SELECT 'notexists' AS resut LIMIT 1"
	rows, err := db.Query(query, user)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	for rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			return false, nil
		}
		if result == "exists" {
			isuser = true
		}
	}
	err = rows.Err()
	if err != nil {
		return false, nil
	}
	return isuser, nil
}

func LogIn_db(iden models.Identify) (bool, error) {
	query := "SELECT 'exists' AS result FROM userid WHERE username = $1 AND password = $2 UNION SELECT 'not exists' AS result LIMIT 1"
	rows, err := db.Query(query, iden.User, iden.Password)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var result string
		err := rows.Scan(&result)
		if err != nil {
			return false, err
		}
		if result == "exists" {
			found = true
		}
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	return found, nil

}
