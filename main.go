package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"log"
	"net/http"
	"strconv"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	// Username - Global Username for testing
	Username = "John"
	// Password - Global Password for testing
	Password = "IamNotACat"

	// OTP Required Data --------

	// Issuer - The site that is creating the OTP Access Code
	Issuer = "localhost:8080"
	// AccountName - The account that the OTP Access Code is linked to
	AccountName = "John@localhost"
)

var (
	// QRCode - in memory representation of the QR Access Code
	QRCode bytes.Buffer
	// KeyOTP - The One Time Access key for validating the access code
	KeyOTP *otp.Key
)

// OTPRequest - The request payload that handles the OTP Access Code
type OTPRequest struct {
	Code string `json:"OTP"`
}

func main() {
	log.Println("Staring...")
	createQRCode()
	http.HandleFunc("/login", login)
	http.HandleFunc("/QrCode", qrCode)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Printf("%+v", r)
	var otp OTPRequest
	err := json.NewDecoder(r.Body).Decode(&otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(nil)
		return
	}

	if check(username, password, otp.Code) {
		w.WriteHeader(http.StatusAccepted)
		w.Write(nil)
		return
	}

	w.WriteHeader(http.StatusForbidden)
	w.Write(nil)
	return
}

// check - Check the username and password against some sort of authentication data store
func check(username, password, otp string) bool {
	if username == Username && password == Password {
		valid := totp.Validate(otp, KeyOTP.Secret())
		if valid {
			return true
		}
	}
	return false
}

func qrCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(QRCode.Bytes())))
	if _, err := w.Write(QRCode.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
	return
}

func createQRCode() {
	key, err := totp.Generate(
		totp.GenerateOpts{
			Issuer:      Issuer,
			AccountName: AccountName,
		})
	if err != nil {
		log.Fatal(err)
	}

	img, err := key.Image(200, 200)
	if err != nil {
		log.Fatal(err)
	}

	// Create a QR Code and store it in memory
	png.Encode(&QRCode, img)
	KeyOTP = key
}
