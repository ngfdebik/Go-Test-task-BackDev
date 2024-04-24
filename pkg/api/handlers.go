package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"example.com/m/internal/services/generator"
	"example.com/m/internal/user"
	"example.com/m/internal/user/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	guid, ok := vars["GUID"]
	if !ok {
		fmt.Errorf("error")
	}

	usr := user.New(guid)
	service := service.NewService(*service.NewStorage())
	flag, u := service.GuidExist(context.Background(), usr.GUID)
	if flag {
		time, _ := strconv.ParseInt(u.AccessIssuedAt, 10, 64)
		rsp, _, err := generator.GetTokens(*usr, time)
		if err != nil {
			fmt.Errorf("error: %v", err)
		}

		rsp.Refrersh = u.RefreshToken

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rsp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	} else {
		rsp, u, err := generator.GetTokens(*usr, jwt.TimeFunc().Unix())
		if err != nil {
			fmt.Errorf("error: %v", err)
		}
		byteRef := []byte(u.RefreshToken)
		hash, _ := bcrypt.GenerateFromPassword(byteRef, 10)
		u.RefreshToken = string(hash)
		_, err = service.Create(context.Background(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rsp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acc, ok := vars["Access"]
	if !ok {
		fmt.Errorf("error")
	}
	ref, ok := vars["Refresh"]
	if !ok {
		fmt.Errorf("error")
	}

	service := service.NewService(*service.NewStorage())
	u, err := service.FindRefresh(context.Background(), ref)
	if err != nil {
		fmt.Errorf("error")
	}
	usr := user.New(u.GUID)

	time, _ := strconv.ParseInt(u.AccessIssuedAt, 10, 64)
	token, u, err := generator.GetTokens(*usr, time)

	if token.Access == acc {
		rsp, u, err := generator.GetTokens(*usr, jwt.TimeFunc().Unix())
		if err != nil {
			fmt.Errorf("error: %v", err)
		}
		byteRef := []byte(u.RefreshToken)
		hash, _ := bcrypt.GenerateFromPassword(byteRef, 10)
		u.RefreshToken = string(hash)

		err = service.Update(context.Background(), u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rsp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		fmt.Errorf("error")
	}
}
