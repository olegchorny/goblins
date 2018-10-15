package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Pallinder/sillyname-go"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

//
type game struct {
	Number string `json:"number" bson:"number"`
	Result string `json:"result" bson:"result"`
}

//
type player struct {
	Number string `json:"number" bson:"number"`
	Name   string `json:"name" bson:"name"`
	Games  []game `json:"games" bson:"games"`
}

//
type tournament struct {
	ID      string   `json:"id" bson:"id"`
	Players []player `json:"players" bson:"players"`
}

func tournamentHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case "POST":
		if err := r.ParseMultipartForm(0); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		//fmt.Fprintf(w, r.FormValue("player1"))

		if r.FormValue("delete") == "true" {
			deleteMongo((strings.Split(r.URL.Path, "/"))[2])
		} else {
			p1 := r.FormValue("player1")
			p2 := r.FormValue("player2")
			p3 := r.FormValue("player3")
			p4 := r.FormValue("player4")

			g11 := game{"1", r.FormValue("player1-1")}
			g12 := game{"2", r.FormValue("player1-2")}
			g13 := game{"3", r.FormValue("player1-3")}
			g14 := game{"4", r.FormValue("player1-4")}
			g21 := game{"1", r.FormValue("player2-1")}
			g22 := game{"2", r.FormValue("player2-2")}
			g23 := game{"3", r.FormValue("player2-3")}
			g24 := game{"4", r.FormValue("player2-4")}
			g31 := game{"1", r.FormValue("player3-1")}
			g32 := game{"2", r.FormValue("player3-2")}
			g33 := game{"3", r.FormValue("player3-3")}
			g34 := game{"4", r.FormValue("player3-4")}
			g41 := game{"1", r.FormValue("player4-1")}
			g42 := game{"2", r.FormValue("player4-2")}
			g43 := game{"3", r.FormValue("player4-3")}
			g44 := game{"4", r.FormValue("player4-4")}

			games1 := []game{g11, g12, g13, g14}
			games2 := []game{g21, g22, g23, g24}
			games3 := []game{g31, g32, g33, g34}
			games4 := []game{g41, g42, g43, g44}

			player1 := player{"1", p1, games1}
			player2 := player{"2", p2, games2}
			player3 := player{"3", p3, games3}
			player4 := player{"4", p4, games4}

			players := []player{player1, player2, player3, player4}

			id := (strings.Split(r.URL.Path, "/"))[2]
			t := tournament{id, players}

			updateMongo(t)
			//log.Println("tournament", (strings.Split(r.URL.Path, "/"))[2])
			// log.Println(p1, g11, g12, g13, g14)
			// log.Println(p2, g21, g22, g23, g24)
			// log.Println(p3, g31, g32, g33, g34)
			// log.Println(p4, g41, g42, g43, g44)

			//log.Println(t)
		}

	case "GET":

		id := (strings.Split(r.URL.Path, "/"))[2]

		filter := bson.NewDocument(bson.EC.String("id", id))

		client, err := mongo.NewClient("mongodb://mongo:27017")
		if err != nil {
			log.Fatal(err)
		}
		err = client.Connect(context.TODO())
		if err != nil {
			log.Fatal(err)
		}

		collection := client.Database("goblins").Collection("tournaments")

		cursor, err := collection.Find(context.Background(), filter)

		var t tournament

		for cursor.Next(context.Background()) {

			err := cursor.Decode(&t)
			if err != nil {
				//handle err
			}
		}

		parseTournament(t)
		http.ServeFile(w, r, "tournament.html")
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}

}

func deleteMongo(id string) {
	client, err := mongo.NewClient("mongodb://mongo:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("goblins").Collection("tournaments")
	//id, _ := collection.Count(context.Background(), nil)
	filter := bson.NewDocument(bson.EC.String("id", id))
	_, err = collection.DeleteOne(context.Background(), filter)
}

func updateMongo(t tournament) {
	client, err := mongo.NewClient("mongodb://mongo:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("goblins").Collection("tournaments")
	//id, _ := collection.Count(context.Background(), nil)

	doc, err := bson.Marshal(t)
	if err != nil {
		log.Fatalln(err)
	}
	// newItemDoc := bson.NewDocument(
	// 	bson.EC.SubDocumentFromElements(
	// 		t.ID,
	// 		bson.EC.SubDocumentFromElements(
	// 			t.Players,
	// 		),
	// 	),
	// )

	//	log.Println(t)
	//log.Println(t.ID)
	//	log.Println(doc)

	filter := bson.NewDocument(bson.EC.String("id", t.ID))
	_, err = collection.ReplaceOne(context.Background(), filter, doc)
	//log.Println(result.DeletedCount)
	//_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Disconnect(context.Background())
}

func writeMongo(t tournament) {
	client, err := mongo.NewClient("mongodb://mongo:27017")
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("goblins").Collection("tournaments")
	//id, _ := collection.Count(context.Background(), nil)

	doc, err := bson.Marshal(t)
	if err != nil {
		log.Fatalln(err)
	}
	// newItemDoc := bson.NewDocument(
	// 	bson.EC.SubDocumentFromElements(
	// 		t.ID,
	// 		bson.EC.SubDocumentFromElements(
	// 			t.Players,
	// 		),
	// 	),
	// )

	_, err2 := collection.InsertOne(context.Background(), doc)
	if err2 != nil {
		log.Fatalln(err2)
	}

	defer client.Disconnect(context.Background())
}

func readMongo() []tournament {
	client, err := mongo.NewClient("mongodb://mongo:27017")
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("goblins").Collection("tournaments")

	cur, err := collection.Find(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(context.Background())
	var tasks []tournament

	//err = ioutil.WriteFile("tasks.html", []byte(r))

	for cur.Next(context.Background()) {
		t := tournament{}
		err := cur.Decode(&t)
		if err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, t)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())
	//	log.Print(tasks)
	return tasks

}

func parseTournament(t tournament) {
	// parse template
	tpl, err := template.ParseFiles("tournament.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	//tasks := readMongo()
	//tasks = sortTasks(tasks)

	f, err := os.Create("tournament.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = tpl.Execute(f, t)
	if err != nil {
		log.Fatalln(err)
	}
}

func parseTable(t []tournament) {
	// parse template
	tpl, err := template.ParseFiles("table.gohtml")
	if err != nil {
		log.Fatalln(err)
	}

	//tasks := readMongo()
	//tasks = sortTasks(tasks)

	f, err := os.Create("table.html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	err = tpl.Execute(f, t)
	if err != nil {
		log.Fatalln(err)
	}
}

func tableHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseMultipartForm(0); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		createTournament()
		//fmt.Fprintf(w, r.FormValue("player1"))

	case "GET":
		t := readMongo()
		parseTable(t)

		http.ServeFile(w, r, "table.html")
	}
	// g1 := game{"1", "0:0"}
	// g2 := game{"2", "0:0"}
	// g3 := game{"3", "0:0"}
	// g4 := game{"4", "0:0"}

	// games0 := []game{g1, g2, g3, g4}
	// // games1 := []game{g4, g2, g3, g1}
	// // games2 := []game{g3, g4, g1, g2}
	// // games3 := []game{g3, g1, g4, g2}
	// // games4 := []game{g3, g2, g1, g4}

	// p1 := player{"1", "Oleg Chorny", games0}
	// p2 := player{"2", "Alex Petrenko", games0}
	// p3 := player{"3", "Denis Shchelkonogov", games0}
	// p4 := player{"4", "Stepan Maksimchuk", games0}

	// p5 := player{"1", "Player1", games0}
	// p6 := player{"2", "Player2", games0}
	// p7 := player{"3", "Player3", games0}
	// p8 := player{"4", "Player4", games0}

	// p9 := player{"1", "Anton", games0}
	// p10 := player{"2", "Elena", games0}
	// p11 := player{"3", "Alexey", games0}
	// p12 := player{"4", "Alex", games0}

	// players1 := []player{p1, p2, p3, p4}
	// players2 := []player{p5, p6, p7, p8}
	// players3 := []player{p1, p2, p9, p10, p11, p12}

	// t1 := tournament{"FNM1", players1}
	// t2 := tournament{"FNM2", players2}
	// t3 := tournament{"FNM3", players3}

	// // t := []tournament{t1, t2, t3}

	// writeMongo(t1)
	// writeMongo(t2)
	// writeMongo(t3)

}

func createTournament() {
	g1 := game{"1", "0:0"}
	g2 := game{"2", "0:0"}
	g3 := game{"3", "0:0"}
	g4 := game{"4", "0:0"}

	g := []game{g1, g2, g3, g4}

	p1 := player{"1", "Player1", g}
	p2 := player{"2", "Player2", g}
	p3 := player{"3", "Player3", g}
	p4 := player{"4", "Player4", g}

	p := []player{p1, p2, p3, p4}

	t := tournament{sillyname.GenerateStupidName(), p}

	writeMongo(t)

}

func main() {

	//fmt.Println(t)
	//r := mux.NewRouter()
	srv := http.NewServeMux()
	srv.Handle("/", http.FileServer(http.Dir(".")))
	srv.HandleFunc("/table.html", tableHandler)
	srv.HandleFunc("/tournament/", tournamentHandler)

	//r := mux.NewRouter()
	//r.HandleFunc("/{key}", tournamentHandler)

	log.Printf("server started")

	log.Fatal(http.ListenAndServe(":5000", srv))

}
