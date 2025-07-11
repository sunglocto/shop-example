package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID uint64
	ShortName string
	LongName string
	Description string
	Rating float64
	Image string
	Stripe string
	Reviews []Review
	Price float64
}

type Review struct {
	ReviewID uint64
	ReviewAuthor string
	ReviewRating float64
	ReviewDescription string
	ReviewTitle string
}

type Stats struct {
	ProductCount uint64
}

var LoadedProducts []Product

func loadProducts() error{
	log.SetPrefix("[INT] ")
	warn := false
	prods := []Product{}
	directory, err := os.ReadDir("./products/")
	if err != nil {
		return err
	}
	for _,v := range directory {
		if !v.IsDir() {
			if strings.ContainsAny(v.Name(), "abcdefghiklmpqrtuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ ") {
				log.Printf("\033[33mWARNING: Product files must only have numbers in their name and no letters (%s)\033[0m\n", v.Name())
				warn = true
			}
			bytes, err := os.ReadFile(fmt.Sprintf("./products/%s", v.Name()))
			if err != nil {
				return err
			}
			newProduct := Product{}
			err = json.Unmarshal(bytes, &newProduct)
			if err != nil {
				return err
			}
			if strconv.Itoa(int(newProduct.ID)) != strings.TrimSuffix(v.Name(),".json") {
				log.Printf("\033[33mWARNING: Product names must match their ID (%s)\n\033[0m", v.Name())
				warn = true
			}
			prods = append(prods, newProduct)
			if int(newProduct.ID) > len(prods) {
				log.Printf("\033[31mERROR: Product IDs must be in order! (%s)\033[0m\n", v.Name())	
				return errors.New("Products must be in order")
			}
		} else {
			log.Printf("\033[33mWARNING: Directories should not be in the products folder (%s)\033[0m\n", v.Name())	
			warn = true
		}
	}
	LoadedProducts = prods
	if !warn {
		log.Println("\033[32mRefreshed all products without any warnings.\033[0m")
	}
	return nil
}

func rootPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"products": LoadedProducts,
	})
}

func refreshProducts(c *gin.Context) {
	err := loadProducts()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %s",err.Error()))
	} else {
		c.String(http.StatusOK, "Successfully refreshed products")
	}
}

func productPage(c *gin.Context) {
	id := c.Param("id")
	i, _ := strconv.Atoi(id)
	if i >= len(LoadedProducts) {
		c.String(http.StatusNotFound, "Product does not exist")
		return
	}
	prod := LoadedProducts[i]
	c.HTML(http.StatusOK, "productpage.tmpl", gin.H{
		"ID": prod.ID,
		"ShortName": prod.ShortName,
		"LongName": prod.LongName,
		"Image": prod.Image,
		"Stripe": prod.Stripe,
		"Description": prod.Description,
		"Reviews": prod.Reviews, // TODO: Make a seperate view for reviews and place it on the page dynamically w/ HTMX
		"Price": prod.Price,
	})
}

func reloadJob() { // This function will run concurrently
	for {
		time.Sleep(60 * time.Second)
		loadProducts()
	}
}

func main() {
	fmt.Println("Running")
	/* Get all products from product folder
	and load them into LoadedProducts     */
	err := loadProducts()
	if err != nil {
		panic(err)
	}

	/* Begin refreshing the products every minute */
	go reloadJob()
	
	r := gin.Default()
	r.LoadHTMLGlob("./views/*")
	r.StaticFS("/css", http.Dir("./static/css/"))
	r.StaticFS("/img", http.Dir("./ugc/product_images"))
	r.GET("/", rootPage)
	r.GET("/refresh", refreshProducts)
	r.GET("/product/:id", productPage)

	r.Run()
}
