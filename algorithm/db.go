package algorithm

import (
	"database/sql"
	"encoding/json"
	v "../variables"
	"log"
)

func GetUserAlgorithms (username string) []UserAlgorithm {
	var userAlgorithms []UserAlgorithm
	var userAlgorithm UserAlgorithm
	var str sql.NullString
	var id sql.NullInt32

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetUserAlgorithms:", err.Error())
	}
	defer db.Close()

	sql := "SELECT name, uniqueid " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1"

	qr, err := db.Query(
		sql,
		username)
	if err != nil {
		log.Fatal("Error querying postgres in GetUserAlgorithms:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str, &id)
		if err != nil {
			log.Fatal("Error scanning query row in GetUserAlgorithms:", err.Error())
		}
		if str.String == ""{
			return []UserAlgorithm{}
		}
		userAlgorithm = UserAlgorithm{
			str.String,
			id.Int32,
		}
		userAlgorithms = append(userAlgorithms, userAlgorithm)
	}
	return userAlgorithms
}
func GetUserAlgorithm (username string, algorithm int) string {
	var str sql.NullString

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetUserAlgorithm:", err.Error())
	}
	defer db.Close()

	sql := "SELECT algorithm " +
		"FROM stock_server.algorithms " +
		"WHERE username = $1 AND " +
		"uniqueid = $2"

	qr, err := db.Query(
		sql,
		username,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in GetUserAlgorithm:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str)
		if err != nil {
			log.Fatal("Error scanning query row in GetUserAlgorithm:", err.Error())
		}
		if str.String == ""{
			return ""
		} else {
			return str.String
		}
	}
	return ""
}
func GetAllAlgorithms () []Algorithm {
	var algorithms []Algorithm
	var str sql.NullString

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in GetAllAlgorithms:", err.Error())
	}
	defer db.Close()

	sql := "SELECT algorithm " +
		"FROM stock_server.algorithms "

	qr, err := db.Query(
		sql)
	if err != nil {
		log.Fatal("Error querying postgres in GetAllAlgorithms:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&str)
		if err != nil {
			log.Fatal("Error scanning query row in GetAllAlgorithms:", err.Error())
		} else {
			if str.String != "" {
				algorithms = append(algorithms, JsonToAlgorithm(str.String))
			}
		}
	}
	return algorithms
}
func UpdateAlgorithm (algorithm Algorithm) {
	var stringAlgorithm = AlgorithmToJson(algorithm)

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in UpdateAlgorithm:", err.Error())
	}
	defer db.Close()

	sql := "UPDATE stock_server.algorithms " +
		"SET algorithm = $1, name = $2" +
		"WHERE uniqueid = $3"

	_, err = db.Exec(
		sql,
		stringAlgorithm,
		algorithm.Name,
		algorithm.UniqueID)
	if err != nil {
		log.Fatal("Error querying postgres in UpdateAlgorithm:", err.Error())
	}

}
func InsertAlgorithm (username string, algorithm Algorithm) {
	var uniqueId sql.NullInt32
	var id int32

	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in InsertAlgorithm:", err.Error())
	}
	defer db.Close()
	sql := "SELECT uniqueid " +
		"FROM stock_server.algorithms " +
		"ORDER BY uniqueid DESC " +
		"LIMIT 1"

	qr, err := db.Query(sql)
	if err != nil {
		log.Fatal("Error querying postgres in InsertAlgorithm:", err.Error())
	}

	for qr.Next() {
		err = qr.Scan(&uniqueId)
		if err != nil {
			log.Fatal("Error scanning query row in InsertAlgorithm:", err.Error())
		}
	}
	sql = "INSERT INTO stock_server.algorithms " +
		"(algorithm, name, uniqueid, username) " +
		"VALUES ($1, $2, $3, $4)"

	if uniqueId.Valid == false {
		id = 1
	} else {
		id = uniqueId.Int32 + 1
	}
	algorithm.UniqueID = int(id)

	var stringAlgorithm = AlgorithmToJson(algorithm)
	_, err = db.Exec(
		sql,
		stringAlgorithm,
		algorithm.Name,
		id,
		username)
	if err != nil {
		log.Fatal("Error querying postgres in InsertAlgorithm:", err.Error())
	}
	AddToAlgorithmMap(algorithm)
	if algorithm.RunOnStart == true {
		StartAlgorithm(algorithm.UserName, algorithm.UniqueID)
	}
}
func DeleteAlgorithm (algorithm int, username string) {
	db, err := sql.Open(v.PGUser, v.ConnectionString)
	if err != nil {
		log.Fatal("Error connecting to postgres in DeleteAlgorithm:", err.Error())
	}
	defer db.Close()

	sql := "DELETE  " +
		"FROM stock_server.algorithms " +
		"WHERE uniqueid = $1"

	_, err = db.Exec(
		sql,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in DeleteAlgorithm:", err.Error())
	}

	sql = "DELETE " +
		"FROM stock_server.notifications " +
		"WHERE algorithm = $1"

	_, err = db.Exec(
		sql,
		algorithm)
	if err != nil {
		log.Fatal("Error querying postgres in DeleteAlgorithm:", err.Error())
	}


}
func AlgorithmToJson(algorithm Algorithm) string {
	str, err := json.Marshal(algorithm)
	if err != nil {
		log.Fatal("Error turning Algorithm into json in AlgorithmToJson: ", err.Error())
	}
	return string(str)

}
func JsonToAlgorithm(str string) Algorithm {
	var algorithm Algorithm
	err := json.Unmarshal([]byte(str), &algorithm)
	if err != nil {
		log.Fatal("Error turning json into Algorithm in JsonToAlgorithm: ", err.Error())
	}
	return algorithm
}