// // gin is similar to nodemon
// // gin --all -i run main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://test:1234@cluster0.h4a8j.mongodb.net/TestDatabase?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 1000000*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	quickstartDatabase := client.Database("TestDatabase")
	peopleCollection := quickstartDatabase.Collection("Music")

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowOrigins:     []string{"*"},
	}))

	e.GET("/music", func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		peopleResultAll, err := peopleCollection.Find(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}

		var people []bson.M
		if err = peopleResultAll.All(ctx, &people); err != nil {
			log.Fatal(err)

		}
		fmt.Println(people)

		return json.NewEncoder(c.Response()).Encode(people)
	})

	type Music struct {
		LowBand  string `json:"lowBand" form:"lowBand" query:"lowBand"`
		LowPeak  string `json:"lowPeak" form:"lowPeak" query:"lowPeak"`
		LowGain  string `json:"lowGain" form:"lowGain" query:"lowGain"`
		HighBand string `json:"highBand" form:"highBand" query:"highBand"`
		HighPeak string `json:"highPeak" form:"highPeak" query:"highPeak"`
		HighGain string `json:"highGain" form:"highGain" query:"highGain"`
	}

	e.POST("/addMusic", func(c echo.Context) error {

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		music := Music{}
		defer c.Request().Body.Close()
		err := json.NewDecoder(c.Request().Body).Decode(&music)

		if err != nil {
			log.Fatalf("Failed reading the request body %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
		}

		// log.Printf("this is your music %#v", music.Name)
		// return c.String(http.StatusOK, "We got your Music!!!")

		// u := new(Music)
		// if err = c.Bind(u); err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Println(u.Name)

		peopleResult, err := peopleCollection.InsertOne(ctx, bson.D{
			{"LowBand", music.LowBand},
			{"LowPeak", music.LowPeak},
			{"LowGain", music.LowGain},
			{"HighBand", music.HighBand},
			{"HighPeak", music.HighPeak},
			{"HighGain", music.HighGain},
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(peopleResult.InsertedID)

		return json.NewEncoder(c.Response()).Encode(peopleResult.InsertedID)

	})

	type MusicID struct {
		ID string `json:"id" form:"id" query:"id"`
	}

	e.POST("/deleteMusic", func(c echo.Context) error {

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		musicID := MusicID{}
		defer c.Request().Body.Close()
		err := json.NewDecoder(c.Request().Body).Decode(&musicID)

		if err != nil {
			log.Fatalf("Failed reading the request body %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
		}

		log.Printf("this is yout music ID: %#v", musicID)

		// Declare a primitive ObjectID from a hexadecimal string
		idPrimitive, err := primitive.ObjectIDFromHex(musicID.ID)
		if err != nil {
			log.Fatal("primitive.ObjectIDFromHex ERROR:", err)
		}

		// Call the DeleteOne() method by passing BSON
		res, err := peopleCollection.DeleteMany(ctx, bson.M{"_id": idPrimitive})
		fmt.Println("DeleteOne Result TYPE:", reflect.TypeOf(res))

		if err != nil {
			log.Fatal("DeleteOne() ERROR:", err)
		}

		return c.String(http.StatusOK, "Music Deleted!!!")

	})

	e.Logger.Print("Listening on port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
