// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "strconv"
)

type Command struct {
    Cmd string
    Message string
    Data interface{}
}

type Card struct{
    Id int
    Text string
    Name string
    X string
    Y string
    Votes []string
}

type Wall struct{
    Title string
    Owner string
    Cards map[string]Card
    Index int
}

func NewWall (t string, o string) (Wall) {
    return Wall{
        Title: t,
        Owner: o,
        Cards: make(map[string]Card),
        Index: 0,
    }
}

func (w *Wall) AddCard (text string, name string) (*Card){
    c := Card{
        Id: w.Index + 1,
        Text: text,
        Name: name,
        X: "40",
        Y: "50",
        Votes: make([]string, 0),
    }
    c.Votes = append(c.Votes,name)
    w.Cards[strconv.Itoa(c.Id)] = c
    w.Index = w.Index + 1
    return &c
}
func (w *Wall) IsVoted (id string, name string) bool {
    c, _ := w.Cards[id]
    return (stringInSlice(name, c.Votes) != -1)
}

func (w *Wall) MoveCard (id string, name string, x string , y string) (*Card, bool) {
    c, ok := w.Cards[id]
    if(ok){
        c.X = x
        c.Y = y
        w.Cards[string(id)] = c
        return &c, true
    }
    return nil, false
}

func (w *Wall) PlusCard (id string, name string) bool {
    c, ok := w.Cards[id]
    if(ok){
        if stringInSlice(name, c.Votes) != -1 {
            return false
        }
        c.Votes = append(c.Votes, name)
        w.Cards[string(id)] = c
        return true
    }
    return false
}

func (w *Wall) UnplusCard (id string, name string) bool {
    c, ok := w.Cards[id]
    if(ok){
        if i := stringInSlice(name, c.Votes); i == -1 {
            return false
        }else{
            c.Votes = append(c.Votes[:i], c.Votes[i+1:]...)
            w.Cards[string(id)] = c
            return true
        }
    }
    return false

}

func (w *Wall) RemoveCard (id int) {
    delete(w.Cards, string(id))
}

func stringInSlice(a string, list []string) int {
    for i, b := range list {
        if b == a {
            return i
        }
    }
    return -1
}

