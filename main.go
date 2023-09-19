

package main

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "fmt"
    "github.com/google/uuid"
  
    "unicode"
    "math"
    "strings"
    "time"
)


// Receipt Items struct
type Item     struct {
    ShortDescription string  `json:"shortDescription"`
    Price            float64 `json:"price"`
} 

// Receipt struct
type receipt    struct {
    ID            string    `json:"id"`
    Retailer      string    `json:"retailer"`
    PurchaseDate  string    `json:"purchaseDate"`
    PurchaseTime  string    `json:"purchaseTime"`
    Total         float64   `json:"total"`
    Items         []Item    `json:"items"`
}

// Receipt ID response struct
type receiptIDresp    struct {
    ID            string    `json:"id"`
}

// Receipt Points response struct
type receiptPoints    struct {
    Points            int    `json:"points"`
}

var receipts []receipt


func main() {
    router := gin.Default()
 // router.ForwardedByClientIP = true
 // router.SetTrustedProxies([]string{"127.0.0.1"})
 // router.GET("/receipts/process/:id", getReceiptByID)

    router.GET("/receipts/process", getReceipts)           // not in specs, but convenient for QA
    router.GET("/receipts/:id/points", getReceiptByID)
    router.POST("/receipts/process", postReceipt)

    router.Run("localhost:8080")
}


// getReceipts responds with the list of all receipts as JSON.
func getReceipts(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, receipts)
}


// postReceipt adds a receipt from JSON received in the request body.
func postReceipt(c *gin.Context) {
    var newReceipt receipt
    var uid string

    // Call BindJSON to bind the received JSON to newReceipt.
    if err := c.BindJSON(&newReceipt); err != nil {
        return
    }

    _, err := time.Parse(time.RFC3339, newReceipt.PurchaseDate + "T" + newReceipt.PurchaseTime + ":00Z")
    if err != nil {
      fmt.Println("Error ",err)
      panic("Invalid Date Error ")
    }

    uuidWithHyphen := uuid.New()
    fmt.Println(uuidWithHyphen)
    uid = uuidWithHyphen.String()
    newReceipt.ID = uuidWithHyphen.String()

    var newReceiptID = receiptIDresp {ID: uid}

    // Add the new receipt to the slice.
    receipts = append(receipts, newReceipt)
    c.IndentedJSON(http.StatusCreated, newReceiptID)
 // c.IndentedJSON(http.StatusCreated, newReceipt)
}


// getReceiptByID locates the receipt whose ID value matches the id
// parameter sent by the client, then returns that receipt as a response
func getReceiptByID(c *gin.Context) {
    id := c.Param("id")

    // Loop over the list of receipts until
    // receipt with ID matching parameter
    // - No match, return: "receipt not found"
    for _, a := range receipts {
        if a.ID == id {
            // c.IndentedJSON(http.StatusOK, a)
            ReceiptPoints(c, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "receipt not found"})
}


func ReceiptPoints(r *gin.Context, ptsReceipt receipt) {
    var Pts int
    var FloatWk float64

    fmt.Println("ptsReceipt is: ", ptsReceipt)
    fmt.Println("ptsReceipt Items are: ", ptsReceipt.Items)

    // Award Points for Alphanumeric characters in Retailer name
    Pts = 0
    for _, r := range ptsReceipt.Retailer {
      if unicode.IsLetter(r) || unicode.IsDigit(r) {
        Pts++
      }
    }
    // fmt.Println("Pts AlpNm: ", Pts)

    // Award Points for even Dollar price
    if math.Round(ptsReceipt.Total) == ptsReceipt.Total {Pts+=50}
    // fmt.Println("Pts evenDollar: ", Pts)

    // Award Points for multiple of .25 price
    FloatWk = ptsReceipt.Total / 0.25
    if math.Round(FloatWk) == FloatWk {Pts+=25}
    // fmt.Println("Pts even25: ", Pts)

    // Award Points for receipt item pairs
    Pts += (len(ptsReceipt.Items) / 2) * 5
    // fmt.Println("Pts itemPairs: ", Pts)

    // Award Points for items with trimmed descriptions lengths = multiples of 3
    for _,ri := range ptsReceipt.Items {
        fmt.Println("ptsReceipt Items are: ", ri)
        if FloatWk = float64(len(strings.TrimSpace(ri.ShortDescription))) / 3; math.Round(FloatWk) == FloatWk {
            Pts+= int(math.Ceil(ri.Price * 0.2))
            // fmt.Println("Pts itemDscrPrc Itm: ", ri, " m.C ",math.Ceil(ri.Price * 0.2), 
            //            " I(m.C) ",int(math.Ceil(ri.Price * 0.2)), "Pts", Pts)
        }
    }
    // fmt.Println("Pts itemDscrPrice: ", Pts)


    // Award Points for odd Day and Time between 2pm and 4pm (non-inclusive)
    DtTm, err := time.Parse(time.RFC3339, ptsReceipt.PurchaseDate + "T" + ptsReceipt.PurchaseTime + ":00Z")
    if err != nil {
      fmt.Println("error ",err)
      panic("Invalid Date error ")

    }
    // fmt.Println("DtTm ",DtTm)

    day := DtTm.Day()     
    // fmt.Println("day ",day)
    hour := DtTm.Hour()
    // fmt.Println("hour ",hour)
    minute := DtTm.Minute()
    // fmt.Println("minute ",minute)

    // Award Points for odd Day
    if FloatWk = float64(day) / 2; math.Round(FloatWk) != FloatWk { Pts += 6 }
    // fmt.Println("Pts oddDay: ", Pts)        

    // Award Points for time between 2pm and 4pm (non-inclusive)  
    if hour == 15 || ( hour == 14 && minute > 0) { Pts += 10 }  // 2pm & 4pm not awarded 10pts
    fmt.Println("Pts time2to4: ", Pts)        

    // Reciept Points Response
    var ReceiptPts = receiptPoints {Points: Pts}
    r.IndentedJSON(http.StatusOK, ReceiptPts)
//  r.IndentedJSON(http.StatusOK, ptsReceipt)
}

