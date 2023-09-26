package initializers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secret = "RYz0umVR1Hxw6JiOdE6L5m7Ne5ZcJB6agbqmxdAlunVuchoPgI+dCWD6QzWJkckgliPzRGaC+55+raGaAKrOX0pQeh2NohNi/CBy7h7Qg+wsTNrdiKwTcG9tARtmOgAvdls6qcDjTwzPaoeA+HxC6ECrpb7zr5jnnJbKkCu3d4QwFp9hiEUWpcChANZciba3Zi3tLELew7fwuC/DjxmXQoq6kfRRhAD7vZuAk1kzRQkzYoX4OMgvgIrJyBx2wnzDpAsAers2j6FB2fuU10aNhZDnYZ68oWGKWv4QEM3bSDESJBqFVdHpgdoErCvmrqyC08lk0KUS+P0t0gsvMARgNg=="

const licenseurl = "https://103.14.123.116:55555/api/v1/licenseCode/validate"

//const licenseurl = "https://localhost:55555/api/v1/licenseCode/validate"

func checkLicenseMain() {

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}
	now_time := time.Now()
	issued_time := now_time.Add(time.Minute * (-1)).Unix()
	expiry_time := now_time.Add(time.Minute * 4).Unix()

	mac, err := getMacAddr()
	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":  os.Getenv("LICENSE_CODE"),
		"aud":  mac,
		"port": port,
		"iat":  issued_time,
		"exp":  expiry_time,
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {

		log.Fatal(err)
	}

	// convert body to JSON data
	jsonData, err := json.Marshal(map[string]string{
		"data": tokenString,
	})
	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}

	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{

				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Post(licenseurl, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {

		log.Fatal(errors.New("license code validation failure:" + resp.Status + "," + string(body)))
	}
	bodyMap := make(map[string]interface{})
	err = json.Unmarshal(body, &bodyMap)

	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatalln(err)

	}

	token1, err := jwt.Parse(bodyMap["data"].(string), func(token2 *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token2.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token2.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal("invalid " + err.Error())

	}

	if claims, ok := token1.Claims.(jwt.MapClaims); ok && token1.Valid {

		//check expiry

		if float64(expiry_time) != claims["exp"].(float64) {
			fmt.Println("license validation failure...exiting the application**********")
			log.Fatal("return token not valid")
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			fmt.Println("license validation failure...exiting the application**********")
			log.Fatal("return token not valid")
		}

		if claims["sta"] != "success" {
			fmt.Println("license validation failure...exiting the application**********")
			log.Fatal("return token not valid")
		}

		split := strings.Split(claims["sub"].(string), ",")

		if split[0] != os.Getenv("LICENSE_CODE") || split[1] != mac || split[2] != os.Getenv("PORT") {
			fmt.Println("license validation failure...exiting the application**********")
			log.Fatal("return token not valid")
		}

		fmt.Println("license is valid******************")

	} else {
		fmt.Println("license validation failure...exiting the application**********")
		log.Fatal("return token not valid")
	}

}

func getMacAddr() (string, error) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
				// Don't use random as we have a real address
				addr := i.HardwareAddr.String()
				return addr, nil
			}
		}

		return "", errors.New("not found an active MAC address")
	} else {

		return "", err
	}

}

func CheckLicense() {

	checkLicenseMain()
	go licenseCheckRoutine()

}

// following routine will check license validity every 6 hours
func licenseCheckRoutine() {

	done := make(chan bool)

	ticker := time.NewTicker(6 * time.Hour)

	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				checkLicenseMain()
			}
		}
	}()

	<-done

}
