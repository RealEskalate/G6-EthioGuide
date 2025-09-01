package domain

import "time"

type AIChat struct {
    ID        string     
    UserID    string     
    Source    string    
    Request   string   
    Response  string    
    Timestamp time.Time 
}