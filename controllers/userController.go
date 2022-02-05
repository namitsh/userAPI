package controller

// all the application logic would come here..
// means validation and calling of service is via controllers

import (
	"UserMicroservice/models"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type e struct {
	Err string `json:"error"`
}

var jwt_key = []byte(os.Getenv("JWT_KEY"))

// Claims CREATE A STRUCT THAT WOULD BE ENCODED TO JWT
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateAuthToken(email string) (string, error) {
	//create an expiration token
	expirationTime := time.Now().Add(5 * time.Minute)

	//create a claim
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create a token with algo HS256 and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// now create a JWT tokenString
	tokenString, err := token.SignedString(jwt_key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

var validate = validator.New()

func hashPassword(pwd string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	return string(bytes), err
}

func checkHashedPassword(pwd string, hashPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(pwd))
	return err == nil
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	defer r.Body.Close()

	// convert json into struct
	log.Println("Converting JSON into struct")
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(e{err.Error()})
		return
	}
	log.Println("JSON converted into struct")
	log.Println("Validating the values of body")
	//	validate the values
	validationErr := validate.Struct(user)
	if validationErr != nil {
		w.WriteHeader(400)
		log.Println(validationErr)
		json.NewEncoder(w).Encode(e{validationErr.Error()})
		return
	}
	log.Println("Successfully validated the values of body")
	log.Println("Checking whether it is existing user.")
	//	check if it is existing user.
	ok, err := models.CheckExistingUser(user.Email)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err)
		json.NewEncoder(w).Encode(e{Err: err.Error()})
		return
	}
	log.Println("No error found while checking the user")
	if ok {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(e{Err: "email id already exist"})
		return
	}
	log.Println("Successfully checked the existence of user: Not found")

	//	hash the password
	hashedPwd, err := hashPassword(user.Password)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(e{Err: err.Error()})
	}
	user.Password = hashedPwd
	//fill out the values
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	fmt.Println("Hello stop")
	token, err := generateAuthToken(user.Email)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(e{Err: err.Error()})
		return
	}
	user.Token = &token

	//create a user in database
	if err := user.SignUp(); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(e{Err: err.Error()})
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request)      {}
func LogoutHandler(w http.ResponseWriter, r *http.Request)     {}
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {}
func GetUserHandler(w http.ResponseWriter, r *http.Request)    {}
