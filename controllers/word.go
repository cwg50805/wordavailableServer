package controllers

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// HeartBeat function
// @Summary Get available words
// @produce application/json
// @Param length query string true "Enter the length of the word here"
// @Param rightPos query string false "Right alphabet at right position. Enter with the format like this -> 0a1p Number stands for digit, and 0 means the first character."
// @Param wrongPos query string false "Right alphabet at wrong position. Enter with the format like this -> 0a1p Number stands for digit, and 0 means the first character."
// @Param wrongAlphabet query string false "Wrong alphabet. Enter with the format like this -> ap . ap means the answer doesn't contain alphabet 'a' and 'p'. "
// @Success 200 {array} array
// @Router /api/v1/words [get]
func GetWords(c *gin.Context) {
	file, err := os.Open("./data/words_alpha.txt")
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}
	}()

	lengthString, exist := c.GetQuery("length")
	if !exist {
		c.JSON(500, gin.H{
			"error": "Enter length",
		})
		return
	}
	lengthInt, _ := strconv.Atoi(lengthString)

	scanner := bufio.NewScanner(file)
	availableWords := []string{}
	for scanner.Scan() {
		if len(scanner.Text()) != lengthInt {
			continue
		}

		if wrongAlphabet, exist := c.GetQuery("wrongAlphabet"); exist {
			wrongAlphabet = strings.ToLower(wrongAlphabet)
			r, _ := regexp.Compile("[^" + wrongAlphabet + "]{" + lengthString + "}")
			if match := r.MatchString(scanner.Text()); !match {
				continue
			}
		}

		if wrongPos, exist := c.GetQuery("wrongPos"); exist {
			wrongPos = strings.ToLower(wrongPos)
			var notContain []string
			for index := 0; index < lengthInt; index++ {
				notContain = append(notContain, "")
			}
			atLeastContain := ""
			for index := 0; index < len(wrongPos)/2; index++ {
				i, _ := strconv.Atoi(string(wrongPos[2*index]))
				notContain[i] += string(wrongPos[2*index+1])
				atLeastContain += string(wrongPos[2*index+1])
			}
			regexString := ""
			for _, digit := range notContain {
				if digit == "" {
					regexString += "."
				} else {
					regexString += "[^" + digit + "]"
				}
			}
			r, _ := regexp.Compile(regexString)
			if match := r.MatchString(scanner.Text()); !match {
				continue
			}
			flag := false
			for _, value := range atLeastContain {
				r, _ := regexp.Compile("[" + string(value) + "]")
				if match := r.MatchString(scanner.Text()); !match {
					flag = true
				}
			}
			if flag {
				continue
			}
		}

		if rightPos, exist := c.GetQuery("rightPos"); exist {
			rightPos = strings.ToLower(rightPos)
			var mustContain []string
			for index := 0; index < lengthInt; index++ {
				mustContain = append(mustContain, "")
			}
			for index := 0; index < len(rightPos)/2; index++ {
				i, _ := strconv.Atoi(string(rightPos[2*index]))
				mustContain[i] += string(rightPos[2*index+1])
			}
			regexString := ""
			for _, digit := range mustContain {
				if digit == "" {
					regexString += "."
				} else {
					regexString += "[" + digit + "]"
				}
			}
			r, _ := regexp.Compile(regexString)
			if match := r.MatchString(scanner.Text()); !match {
				continue
			}
		}
		availableWords = append(availableWords, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": availableWords,
	})
}
