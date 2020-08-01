package restaurantunit

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

type Restaurant struct {
	name           string
	Latitude       float64
	Longitude      float64
	address        string
	photoReference string
}

func GetRestaurant(Latitude, Longitude float64) *Restaurant {
	//var jsonObj map[string]interface{}
	//json.Unmarshal(getJSONFromLocation(Latitude, Longitude), &jsonObj)
	//Todo:https://ithelp.ithome.com.tw/articles/10205062?sc=iThelpR
	mapData := getJSONFromLocation(Latitude, Longitude)
	oneRestaurant := GetOneRestaurant(mapData)
	return oneRestaurant
}

func GetOneRestaurant(mapData string) *Restaurant {
	oneRestaurant := Restaurant{}
	results := gjson.Get(mapData, "results")
	if results.IsArray() {
		nowJSON := results.Array()[rand.Intn(len(results.Array()))].String()
		fmt.Println(nowJSON)
		businessStatus := gjson.Get(nowJSON, "business_status")
		if businessStatus.String() == "OPERATIONAL" {
			name := gjson.Get(nowJSON, "name")
			Latitude := gjson.Get(nowJSON, "geometry.location.lat")
			Longitude := gjson.Get(nowJSON, "geometry.location.lng")
			address := gjson.Get(nowJSON, "vicinity")
			photoReference := gjson.Get(nowJSON, "photos.0.photo_reference")
			log.Println("photoReference=", photoReference)
			//geometry := gjson.Get(nowJson ,"geometry")
			// log.Println("name=", name)
			// log.Println("Latitude =", Latitude, ", Longitude =", Longitude)
			Lat, err := strconv.ParseFloat(Latitude.String(), 8)
			Lon, err := strconv.ParseFloat(Longitude.String(), 8)
			if err != nil {
				return &oneRestaurant
			}
			oneRestaurant.name = name.String()
			oneRestaurant.Latitude = Lat
			oneRestaurant.Longitude = Lon
			oneRestaurant.address = address.String()
			oneRestaurant.photoReference = photoReference.String()
		}
	}
	/*
	   for i := 0 ; i < len(results.Array()) ; i++{
	       nowJson := results.Array()[i].String()
	       business_status:= gjson.Get(nowJson ,"business_status")
	       if business_status.String() == "OPERATIONAL" {
	           name := gjson.Get(nowJson ,"name")
	           geometry := gjson.Get(nowJson ,"geometry")
	           fmt.Println("name=",name)
	           fmt.Println("geometry=",geometry)
	           fmt.Println("====================")
	       }
	   }
	*/
	return &oneRestaurant
}

func getJSONFromLocation(Latitude, Longitude float64) string {
	radius := "200"
	googleURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?radius="
	googleURL += radius + "&type=restaurant"
	googleURL += "&location=" + fmt.Sprintf("%f", Latitude) + "," + fmt.Sprintf("%f", Longitude)
	googleURL += "&key=" + os.Getenv("GoogleKey")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	res, err := client.Get(googleURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	//gitfmt.Printf("%s", sitemap)
	status := gjson.Get(string(sitemap), "status")
	var mapResult string
	if status.String() == "OK" {
		mapResult = string(sitemap)
	} else {
		mapResult = ""
	}
	return mapResult
}
