package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

var fcmClient *messaging.Client

// Body request
type NotificationRequest struct {
	Title      string                   `json:"title"` // Title in notification
	Body       string                   `json:"body"`  //  Body in notification
	Token      string                   `json:"token" validate:"required"`
	ImageURL   string                   `json:"imageUrl"`   // "Image URL" in notification
	Data       map[string]string        `json:"data"`       // Custom data (optional)
	Android    *messaging.AndroidConfig `json:"android"`    // Android
	APNS       *messaging.APNSConfig    `json:"ios"`        // Ios
	Webpush    *messaging.WebpushConfig `json:"webpush"`    // Web Push
	Topic      string                   `json:"topic"`      // Topic
	FCMOptions *messaging.FCMOptions    `json:"fcmoptions"` // Fcm Options
	Condition  string                   `json:"condition"`  // Condition
}
type Name struct {
	Name string `json:"project_id"`
}

// Initial Firebase
func initFirebase(serviceAccountPath string) {
	opt := option.WithCredentialsFile(serviceAccountPath)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error when init App Firebase: %v", err.Error())
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("Error when initial FCM client: %v", err.Error())
	}

	fcmClient = client

	var data []byte
	data, _ = ioutil.ReadFile(serviceAccountPath)

	var str Name
	_ = json.Unmarshal(data, &str)

	log.Printf("Firebase appId: %s and FCM client initial success!", str.Name)
}

func sendNotification(c *gin.Context) {
	var req NotificationRequest

	// Validator body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parse error in request body!"})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required!"})
		return
	}

	// create message body from request
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    req.Title,
			Body:     req.Body,
			ImageURL: req.ImageURL,
		},
		Token:      req.Token,
		Data:       req.Data,
		Android:    req.Android,
		APNS:       req.APNS,
		Webpush:    req.Webpush,
		Topic:      req.Topic,
		Condition:  req.Condition,
		FCMOptions: req.FCMOptions,
	}

	response, err := fcmClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("Send message error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Send message error!."})
		return
	}

	log.Printf("FCM response: %s", response)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Send message successfully!",
		"response": response,
	})
}

func sendNotificationMultipleDevice(c *gin.Context) {
	// Body request
	type NotificationRequest struct {
		Title      string                   `json:"title"` // Title in notification
		Body       string                   `json:"body"`  //  Body in notification
		Token      []string                 `json:"token" validate:"required"`
		ImageURL   string                   `json:"imageUrl"`   // "Image URL" in notification
		Data       map[string]string        `json:"data"`       // Custom data (optional)
		Android    *messaging.AndroidConfig `json:"android"`    // Android
		APNS       *messaging.APNSConfig    `json:"ios"`        // Ios
		Webpush    *messaging.WebpushConfig `json:"webpush"`    // Web Push
		Topic      string                   `json:"topic"`      // Topic
		FCMOptions *messaging.FCMOptions    `json:"fcmoptions"` // Fcm Options
		Condition  string                   `json:"condition"`  // Condition
	}

	var req NotificationRequest

	// Validator body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parse error in request body!"})
		return
	}

	if len(req.Token) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required!"})
		return
	}

	fmt.Println(req.Title)
	// create message body from request

	responses := make([]string, 0)

	for _, token := range req.Token {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title:    req.Title,
				Body:     req.Body,
				ImageURL: req.ImageURL,
			},
			Token:      token,
			Data:       req.Data,
			Android:    req.Android,
			APNS:       req.APNS,
			Webpush:    req.Webpush,
			Topic:      req.Topic,
			Condition:  req.Condition,
			FCMOptions: req.FCMOptions,
		}
		response, err := fcmClient.Send(context.Background(), message)

		if err != nil {
			responses = append(responses, err.Error())
		} else {
			responses = append(responses, response)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Send message successfully!",
		"response": responses,
	})
}

func sendNotificationBulk(c *gin.Context) {
	var req []NotificationRequest

	// Validator body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parse error in request body!"})
		return
	}

	// create message body from request
	messages := []*messaging.Message{}
	for _, v := range req {
		if v.Token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "token is required!"})
			return
		}

		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title:    v.Title,
				Body:     v.Body,
				ImageURL: v.ImageURL,
			},
			Token:      v.Token,
			Data:       v.Data,
			Android:    v.Android,
			APNS:       v.APNS,
			Webpush:    v.Webpush,
			Topic:      v.Topic,
			Condition:  v.Condition,
			FCMOptions: v.FCMOptions,
		}
		messages = append(messages, message)

	}

	responses := []string{}
	errResponse := []error{}

	for _, m := range messages {

		response, err := fcmClient.Send(context.Background(), m)
		responses = append(responses, response)
		errResponse = append(errResponse, err)

		if err != nil {
			log.Printf("Send message error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Send message error!.", "data": errResponse})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Send message successfully!",
		"response": responses,
	})
}

func main() {
	// Định nghĩa tham số -p cho đường dẫn tệp
	serviceAccountPath := flag.String("p", "config/serviceAccount.json", "Path to serviceAccount.json file")
	flag.Parse()

	fmt.Println("        /\\         /\\    ")
	fmt.Println("       /  \\       /  \\   ")
	fmt.Println("      /    \\_____/    \\  ")
	fmt.Println("     /                  \\ ")
	fmt.Println("    /      AQUILA        \\")
	fmt.Println("   /        FCM           \\")
	fmt.Println("  /                       /")
	fmt.Println(" /_______________________/ ")
	fmt.Println("      |            |       ")
	fmt.Println("      |            |       ")
	fmt.Println("      |            |       ")
	fmt.Println("     /             \\      ")
	fmt.Println("    /_______________\\     ")
	fmt.Println("********************************************")
	fmt.Println("*      FCM send message debug service      *")
	fmt.Println("* Author: Pham Van Khang                   *")
	fmt.Println("* License: MIT                             *")
	fmt.Println("********************************************")

	initFirebase(*serviceAccountPath)

	gin.SetMode("release")
	router := gin.Default()

	router.POST("/send", sendNotification)
	router.PUT("/send", sendNotificationMultipleDevice)

	router.POST("/send/bulk", sendNotificationBulk)

	log.Println("Server listening at http://127.0.0.1:8080")
	log.Println("Server listening at http://::1:8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error when run server: %v", err)
	}
}
