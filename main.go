package main

import (
    "fmt"
    "github.com/muraokamasaki/search_engine/indices"
)

func main() {
    ii := indices.NewInvertedIndex()
    ii.BuildFromTextFile("example.txt")
    query := "matrix value singular"
    ansx := ii.Union(indices.Tokenize(query))
    fmt.Println(ansx)


    ki := indices.NewKGramIndex(3)
    ki.BuildFromTextFile("example.txt")
    o := ki.KGramOverlap("cohen", 4)
    fmt.Println(o)
    /*
    for _, q := range []string{"hello", "hell", "heler", "ap", "ape"} {
        answer := []string{}
        for i := 0; i < len(q) - 3; i++ {
            fmt.Println("checking", q[i : i + 3])
            ans, ok := m[q[i: i + 3]]
            if !ok {
                answer = []string{}
                break
            }
            for _, a := range ans {
                answer = append(answer, a)
            }
        }
        fmt.Println(q, answer)
    }
*/
    fmt.Println(indices.PrefixEditDistance("uniwer", "university", 2))
}