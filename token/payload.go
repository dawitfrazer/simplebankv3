package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//differnt types of error retruned by the VerifyToken function 
var (
	ErrInvlaidToken = errors.New("token is invlaid")
	ErrExpiredToken = errors.New("token has expired")
)
//Payload contains the payload data of the token
type Payload struct {
	ID 					uuid.UUID 	`json:"id"`
	Username 		string 			`json:"username"`
	IssuedAt		time.Time		`json:"issued_at"`
	ExpiredAt		time.Time		`json:"expired_at"`
}

// NewPayload cretaes a new token payload with a specific username and duration 
func NewPaylaod(username string, duration time.Duration)(*Payload, error){
		tokenID, err := uuid.NewRandom()
		if err != nil{
			return nil, err
		}

		payload := &Payload{
			ID:						tokenID,
			Username: 		username,
			IssuedAt: 		time.Now(),
			ExpiredAt: 		time.Now().Add(duration),		
		} 

		return payload, nil
}

func(payloaf *Payload) Valid() error{
	if time.Now().After(payloaf.ExpiredAt){
		return  ErrExpiredToken
	}

	return nil
}