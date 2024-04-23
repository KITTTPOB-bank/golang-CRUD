package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/KITTTPOB-bank/Carnan_backend/pkg/Key"
	"github.com/KITTTPOB-bank/Carnan_backend/pkg/dbconfig"
	"github.com/KITTTPOB-bank/Carnan_backend/pkg/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db = dbconfig.DBconnect().Database("carpartdata")

type Cardata models.Cardata
type Carname models.Carname
type Carpartdataall models.Carpartdataall
type Formdatacar models.Formdatacar
type Carpart models.Carpart
type Cariamge models.Cariamge
type Dropcarpart models.Dropcarpart
type Addcarpart models.Addcarpart
type FormJsonfilecarpart models.FormJsonfilecarpart
type FormJsonfileMycarpart models.FormJsonfileMycarpart

func GetAWS3() *s3.Client {
	var aws_access_key_id, aws_secret_access_key = Key.GetIDAWSs3()
	appCreds := credentials.NewStaticCredentialsProvider(
		aws_access_key_id,
		aws_secret_access_key,
		"",
	)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-southeast-2"),
		config.WithCredentialsProvider(appCreds),
	)
	if err != nil {
		fmt.Println(err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return s3Client
}
func Getallcar(c *gin.Context) {
	var getalldata []Cardata
	cursor, err := db.Collection("car").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(context.TODO(), &getalldata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, getalldata)
}
func Getallcaradmin(c *gin.Context) {
	cursor, err := db.Collection("car").Find(context.TODO(), bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var getcaradmindata []bson.M

	if err := cursor.All(context.TODO(), &getcaradmindata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, getcaradmindata)
}

func Addcar(c *gin.Context) {
	var form Formdatacar
	// var name string
	getawss3 := GetAWS3()

	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	file, err2 := form.Car_image.Open()
	if err2 != nil {
		fmt.Println(err2)
	} else {

		folderPath := "imagepeviewcar"
		bucketName := "carimageapp"
		objectKey := folderPath + "/" + form.Car_image.Filename

		_, err := getawss3.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
			Body:   file,
			ACL:    "public-read",
		})
		if err != nil {
			fmt.Println(err)
		}

	}
	defer file.Close()
	filename := "https://carimageapp.s3.ap-southeast-2.amazonaws.com/imagepeviewcar/" + form.Car_image.Filename

	findOptions := options.FindOne().SetSort(map[string]int{"_id": -1})
	type ID struct {
		Id int `bson:"_id"`
	}
	var idmaxcar ID
	err := db.Collection("car").FindOne(context.Background(), bson.D{}, findOptions).Decode(&idmaxcar)
	if err != nil {
		fmt.Println("Error finding document:", err)
		return
	}
	addnewcar := bson.D{
		{Key: "_id", Value: idmaxcar.Id + 1},
		{Key: "name", Value: form.Name},
		{Key: "brand", Value: form.Brand},
		{Key: "year", Value: form.Year},
		{Key: "desc", Value: form.Desc},
		{Key: "frontbumper_list", Value: []int64{0}},
		{Key: "rearbumper_list", Value: []int64{0}},
		{Key: "grille_list", Value: []int64{0}},
		{Key: "headlamp_list", Value: []int64{0}},
		{Key: "backuplamp_list", Value: []int64{0}},
		{Key: "mirror_list", Value: []int64{0}},
		{Key: "door_list", Value: []int64{0}},
		{Key: "car_image", Value: filename},
	}

	db.Collection("car").InsertOne(context.TODO(), addnewcar)
}
func Dropcar(c *gin.Context) {
	delid := c.Param("delid")
	id, err := strconv.Atoi(delid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	deleteResult, err := db.Collection("car").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	json.NewEncoder(c.Writer).Encode(deleteResult)
}

func Getallcarname(c *gin.Context) {
	var getalldata []Carname
	cursor, err := db.Collection("car").Find(context.TODO(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := cursor.All(context.TODO(), &getalldata); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, getalldata)
}

func GetMyCar(c *gin.Context) {
	type Getallcarpart struct {
		ID    int32   `bson:"_id" json:"_id"`
		Code  string  `json:"code"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	var getcarpart Carpartdataall
	var frontbumper_list []Getallcarpart
	var rearbumper_list []Getallcarpart
	var backuplamp_list []Getallcarpart
	var headlamp_list []Getallcarpart
	var mirror_list []Getallcarpart
	var door_list []Getallcarpart
	var grille_list []Getallcarpart
	var notfrontbumper_list []Getallcarpart
	var notrearbumper_list []Getallcarpart
	var notbackuplamp_list []Getallcarpart
	var notheadlamp_list []Getallcarpart
	var notmirror_list []Getallcarpart
	var notdoor_list []Getallcarpart
	var notgrille_list []Getallcarpart
	brand := c.Param("brand")
	name := c.Param("name")
	year := c.Param("year")
	yearuse, _ := strconv.Atoi(year)
	check := c.Param("check")
	err := db.Collection("car").FindOne(context.TODO(), bson.M{"brand": brand, "name": name, "year": yearuse}).Decode(&getcarpart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	cursor1, _ := db.Collection("frontbumper").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Frontbumper_list}})
	cursor2, _ := db.Collection("rearbumper").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Rearbumper_list}})
	cursor3, _ := db.Collection("backuplamp").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Backuplamp_list}})
	cursor4, _ := db.Collection("headlamp").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Headlamp_list}})
	cursor5, _ := db.Collection("mirror").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Mirror_list}})
	cursor6, _ := db.Collection("door").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Door_list}})
	cursor7, _ := db.Collection("grille").Find(context.Background(), bson.M{"_id": bson.M{"$in": getcarpart.Grille_list}})

	cursor1.All(context.TODO(), &frontbumper_list)
	cursor2.All(context.TODO(), &rearbumper_list)
	cursor3.All(context.TODO(), &backuplamp_list)
	cursor4.All(context.TODO(), &headlamp_list)
	cursor5.All(context.TODO(), &mirror_list)
	cursor6.All(context.TODO(), &door_list)
	cursor7.All(context.TODO(), &grille_list)

	if check == "1" {
		cursor8, _ := db.Collection("frontbumper").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Frontbumper_list}})
		cursor9, _ := db.Collection("rearbumper").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Rearbumper_list}})
		cursor10, _ := db.Collection("backuplamp").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Backuplamp_list}})
		cursor11, _ := db.Collection("headlamp").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Headlamp_list}})
		cursor12, _ := db.Collection("mirror").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Mirror_list}})
		cursor13, _ := db.Collection("door").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Door_list}})
		cursor14, _ := db.Collection("grille").Find(context.Background(), bson.M{"_id": bson.M{"$nin": getcarpart.Grille_list}})

		cursor8.All(context.TODO(), &notfrontbumper_list)
		cursor9.All(context.TODO(), &notrearbumper_list)
		cursor10.All(context.TODO(), &notbackuplamp_list)
		cursor11.All(context.TODO(), &notheadlamp_list)
		cursor12.All(context.TODO(), &notmirror_list)
		cursor13.All(context.TODO(), &notdoor_list)
		cursor14.All(context.TODO(), &notgrille_list)
	}
	GoData := bson.M{
		"frontbumper":     frontbumper_list,
		"rearbumper":      rearbumper_list,
		"grille":          grille_list,
		"headlamp":        headlamp_list,
		"backuplamp":      backuplamp_list,
		"mirror":          mirror_list,
		"door":            door_list,
		"all_frontbumper": notfrontbumper_list,
		"all_rearbumper":  notrearbumper_list,
		"all_grille":      notgrille_list,
		"all_headlamp":    notheadlamp_list,
		"all_backuplamp":  notbackuplamp_list,
		"all_mirror":      notmirror_list,
		"all_door":        notdoor_list,
	}
	c.JSON(http.StatusOK, GoData)
}
func Addpart(c *gin.Context) {
	var carpartdata = &Carpart{}
	if err := json.NewDecoder(c.Request.Body).Decode(&carpartdata); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	check := checktype(carpartdata.Type)
	yearuse, _ := strconv.Atoi(carpartdata.Year)

	filter := bson.M{"name": carpartdata.Name, "brand": carpartdata.Brand, "year": yearuse}
	update := bson.M{"$push": bson.M{check: bson.M{"$each": carpartdata.IDcarpart}}}

	_, err := db.Collection("car").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}
}
func DropcarPart(c *gin.Context) {
	var Dropcarpart = &Dropcarpart{}
	if err := json.NewDecoder(c.Request.Body).Decode(&Dropcarpart); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return

	}
	yearuse, _ := strconv.Atoi(Dropcarpart.Year)
	check := checktype(Dropcarpart.Type)
	filter := bson.M{"name": Dropcarpart.Name, "brand": Dropcarpart.Brand, "year": yearuse}

	var getdata Carpartdataall

	err := db.Collection("car").FindOne(context.TODO(), filter).Decode(&getdata)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}
	var newarray []int
	switch check {
	case "frontbumper_list":
		newarray = findAndRemoveItem(getdata.Frontbumper_list, Dropcarpart.Delid)
	case "rearbumper_list":
		newarray = findAndRemoveItem(getdata.Rearbumper_list, Dropcarpart.Delid)
	case "grille_list":
		newarray = findAndRemoveItem(getdata.Grille_list, Dropcarpart.Delid)
	case "mirror_list":
		newarray = findAndRemoveItem(getdata.Mirror_list, Dropcarpart.Delid)
	case "headlamp_list":
		newarray = findAndRemoveItem(getdata.Grille_list, Dropcarpart.Delid)
	case "backuplamp_list":
		newarray = findAndRemoveItem(getdata.Mirror_list, Dropcarpart.Delid)
	default:
		newarray = findAndRemoveItem(getdata.Door_list, Dropcarpart.Delid)
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: check, Value: newarray},
		}},
	}

	_, err2 := db.Collection("car").UpdateOne(context.TODO(), filter, update)
	if err2 != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}
}

func findAndRemoveItem(list []int, id int) []int {
	var index int
	for i, item := range list {
		if item == id {
			index = i
			break
		}
	}
	return append(list[:index], list[index+1:]...)
}
func checktype(typenum int) (result string) {
	if typenum == 0 {
		result = "frontbumper_list"
	} else if typenum == 1 {
		result = "rearbumper_list"
	} else if typenum == 2 {
		result = "grille_list"
	} else if typenum == 3 {
		result = "mirror_list"
	} else if typenum == 4 {
		result = "headlamp_list"
	} else if typenum == 5 {
		result = "backuplamp_list"
	} else if typenum == 6 {
		result = "door_list"
	}
	return result
}

func Getcarpartall(c *gin.Context) {
	type Getallcarpart struct {
		ID    int64   `bson:"_id" json:"_id"`
		Code  string  `json:"code"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	var frontbumper_list []Getallcarpart
	var rearbumper_list []Getallcarpart
	var backuplamp_list []Getallcarpart
	var headlamp_list []Getallcarpart
	var mirror_list []Getallcarpart
	var door_list []Getallcarpart
	var grille_list []Getallcarpart

	cursor1, _ := db.Collection("frontbumper").Find(context.TODO(), bson.M{})
	cursor2, _ := db.Collection("rearbumper").Find(context.TODO(), bson.M{})
	cursor3, _ := db.Collection("backuplamp").Find(context.TODO(), bson.M{})
	cursor4, _ := db.Collection("headlamp").Find(context.TODO(), bson.M{})
	cursor5, _ := db.Collection("mirror").Find(context.TODO(), bson.M{})
	cursor6, _ := db.Collection("door").Find(context.TODO(), bson.M{})
	cursor7, _ := db.Collection("grille").Find(context.TODO(), bson.M{})

	cursor1.All(context.TODO(), &frontbumper_list)
	cursor2.All(context.TODO(), &rearbumper_list)
	cursor3.All(context.TODO(), &backuplamp_list)
	cursor4.All(context.TODO(), &headlamp_list)
	cursor5.All(context.TODO(), &mirror_list)
	cursor6.All(context.TODO(), &door_list)
	cursor7.All(context.TODO(), &grille_list)

	Carpartall := bson.M{
		"all_frontbumper": frontbumper_list,
		"all_rearbumper":  rearbumper_list,
		"all_grille":      grille_list,
		"all_headlamp":    headlamp_list,
		"all_backuplamp":  backuplamp_list,
		"all_mirror":      mirror_list,
		"all_door":        door_list,
	}

	c.JSON(http.StatusOK, Carpartall)
}
func Addallpart(c *gin.Context) {
	var newcarpart = &Addcarpart{}
	if err := json.NewDecoder(c.Request.Body).Decode(&newcarpart); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	newprice, _ := strconv.ParseFloat(newcarpart.Price, 64)
	check := checktypecarpartall(newcarpart.Type)
	findOptions := options.FindOne().SetSort(map[string]int{"_id": -1})

	type ID struct {
		Id int `bson:"_id"`
	}
	var idmaxcar ID
	err := db.Collection(check).FindOne(context.Background(), bson.D{}, findOptions).Decode(&idmaxcar)

	newcarparttoData := bson.D{
		{Key: "_id", Value: idmaxcar.Id + 1},
		{Key: "code", Value: newcarpart.Code},
		{Key: "name", Value: newcarpart.Name},
		{Key: "price", Value: newprice},
	}
	_, err2 := db.Collection(check).InsertOne(context.TODO(), newcarparttoData)
	if err2 != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}
}
func DeleteMypart(c *gin.Context) {
	id := c.Param("id")
	typecheck := c.Param("typesend")
	fmt.Println(id)
	fmt.Println(typecheck)

	delid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	typedata, err := strconv.Atoi(typecheck)
	check := checktypecarpartall(typedata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	deleteResult, err := db.Collection(check).DeleteOne(context.TODO(), bson.M{"_id": delid})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	fmt.Println(err)
	json.NewEncoder(c.Writer).Encode(deleteResult)
}

func AddJsonfilecarpart(c *gin.Context) {
	var form FormJsonfilecarpart
	typecheck := c.Param("typesend")
	typedata, _ := strconv.Atoi(typecheck)
	check := checktypecarpartall(typedata)

	//หาไอดีสูงสุด
	findOptions := options.FindOne().SetSort(map[string]int{"_id": -1})
	type ID struct {
		Id int `bson:"_id"`
	}
	var idmaxcar ID
	err2 := db.Collection(check).FindOne(context.Background(), bson.D{}, findOptions).Decode(&idmaxcar)
	if err2 != nil {
		fmt.Println("Error finding document:", err2)
		return
	}

	err := c.ShouldBind(&form)
	if err != nil {
		fmt.Println(err)
	}

	file, err3 := form.FileJson.Open()
	if err3 != nil {
		fmt.Println("Error opening file:", err3)
		return
	}
	defer file.Close()

	fileContent, err4 := io.ReadAll(file)
	if err4 != nil {
		fmt.Println("Error reading file content:", err4)
		return
	}

	type Document struct {
		Name  string  `json:"name"`
		Code  string  `json:"code"`
		Price float64 `json:"price"`
	}
	var documents []Document

	_ = json.Unmarshal(fileContent, &documents)

	for i, item := range documents {
		addmewpart := bson.D{
			{Key: "_id", Value: idmaxcar.Id + i + 1},
			{Key: "code", Value: item.Code},
			{Key: "name", Value: item.Name},
			{Key: "price", Value: item.Price},
		}
		_, err := db.Collection(check).InsertOne(context.TODO(), addmewpart)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// var items []interface{}
	// for i, item := range documents {
	// 	item._id = idmaxcar.Id + i
	// 	items = append(items, item)
	// }

	// db.Collection(check).InsertMany(context.TODO(), items)
}
func checktypecarpartall(typenum int) (result string) {
	if typenum == 0 {
		result = "frontbumper"
	} else if typenum == 1 {
		result = "rearbumper"
	} else if typenum == 2 {
		result = "grille"
	} else if typenum == 3 {
		result = "mirror"
	} else if typenum == 4 {
		result = "headlamp"
	} else if typenum == 5 {
		result = "backuplamp"
	} else if typenum == 6 {
		result = "door"
	} else {
		result = "body"

	}
	return result
}
func AddJsonfileMycarpart(c *gin.Context) {
	var form FormJsonfileMycarpart
	// var Mycarpart = &Carpartdataall{}
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}
	file, err3 := form.FileJson.Open()
	if err3 != nil {
		fmt.Println("Error opening file:", err3)
		return
	}
	defer file.Close()

	fileContent, err4 := io.ReadAll(file)
	if err4 != nil {
		fmt.Println("Error reading file content:", err4)
		return
	}
	type Document struct {
		Code string `json:"code"`
	}
	var documents []Document

	_ = json.Unmarshal(fileContent, &documents)
	checkpart := checktypecarpartall(form.Typecheck)
	type GetId struct {
		ID int `bson:"_id"`
	}
	var getid GetId
	var keepid []interface{}
	fmt.Println(documents)
	for _, item := range documents {
		err := db.Collection(checkpart).FindOne(context.Background(), bson.M{"code": item.Code}).Decode(&getid)
		fmt.Println(err)
		keepid = append(keepid, getid.ID)
	}
	fmt.Println(keepid)

	checkpartincar := checktype(form.Typecheck)

	filter := bson.M{"name": form.Name, "brand": form.Brand, "year": form.Year}

	_, err := db.Collection("car").UpdateOne(context.TODO(), filter, bson.M{"$push": bson.M{checkpartincar: bson.M{"$each": keepid}}})
	if err != nil {
		fmt.Println(err)
		return
	}

}
func Uploadcarimage(c *gin.Context) {
	c.Request.ParseMultipartForm(100)
	var data Cariamge
	c.Request.ParseMultipartForm(100 * 1024 * 1024)
	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, "bad request: %v", err)
		return
	}

	var files []*multipart.FileHeader

	form := c.Request.MultipartForm
	if form != nil {
		files = form.File["files"]
	}

	if len(files) == 0 {
		c.String(http.StatusBadRequest, "no files uploaded")
		return
	}
	yearuse := strconv.Itoa(data.Year)
	check := checktypecarpartall(data.Typecheck)
	folderPath := "Trendmodelpicture/" + data.Brand + " " + data.Name + " " + yearuse + "/" + check
	bucketName := "carimageapp"

	getawss3 := GetAWS3()

	for _, file := range files {
		fileopen, _ := file.Open()
		objectKey := folderPath + "/" + file.Filename
		_, err := getawss3.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: &bucketName,
			Key:    &objectKey,
			Body:   fileopen,
			ACL:    "public-read",
		})
		if err != nil {
			fmt.Println(err)
		}
		defer fileopen.Close()
	}

}

func Getaws3image(c *gin.Context) {
	name := c.Param("name")
	brand := c.Param("brand")
	year := c.Param("year")
	objectKey := "Trendmodelpicture/" + brand + " " + name + " " + year
	bucketName := "carimageapp"
	getawss3 := GetAWS3()

	res, err := getawss3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix: &objectKey,
	})
	if err != nil {
		fmt.Println(err)

	}

	imageURIsByCategory := make(map[string][]string)
	s3Folders := map[string]string{
		"frontbumper": objectKey + "/frontbumper",
		"rearbumper":  objectKey + "/rearbumper",
		"headlamp":    objectKey + "/headlamp",
		"backuplamp":  objectKey + "/backuplamp",
		"door":        objectKey + "/door",
		"body":        objectKey + "/body",
	}

	for category := range s3Folders {
		for _, obj := range res.Contents {
			imageURI := fmt.Sprintf("https://%s.s3.ap-southeast-2.amazonaws.com/%s", bucketName, *obj.Key)
			imageURIsByCategory[category] = append(imageURIsByCategory[category], imageURI)
		}
	}
	c.JSON(http.StatusOK, imageURIsByCategory)

}
