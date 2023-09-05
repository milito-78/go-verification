package go_verification

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CodeGenerator interface {
	Generate() string
}

type NumberGenerator struct {
	length         int
	notZeroAtStart bool
}

func NewNumberGenerator(length int, notZeroAtStart bool) *NumberGenerator {
	return &NumberGenerator{length: length, notZeroAtStart: notZeroAtStart}
}

func (n NumberGenerator) Generate() string {
	rand.Seed(time.Now().UnixNano())
	chars := "0123456789"
	result := make([]byte, n.length)
	for i := 0; i < n.length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	if n.notZeroAtStart && result[0] == '0' {
		chars = "123456789"
		result[0] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

type AlphabetGenerator struct {
	length        int
	allCapital    bool
	allNonCapital bool
}

func NewAlphabetGenerator(length int, allCapital bool, allNonCapital bool) *AlphabetGenerator {
	return &AlphabetGenerator{length: length, allCapital: allCapital, allNonCapital: allNonCapital}
}

func (n AlphabetGenerator) Generate() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if n.allCapital {
		chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	} else if n.allNonCapital {
		chars = "abcdefghijklmnopqrstuvwxyz"
	}

	result := make([]byte, n.length)
	for i := 0; i < n.length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type WordGenerator struct {
	length int
}

func NewWordGenerator(length int) *WordGenerator {
	return &WordGenerator{length: length}
}

func (n WordGenerator) Generate() string {
	rand.Seed(time.Now().UnixNano())
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, n.length)
	for i := 0; i < n.length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type RegexGenerator struct {
	regex string
}

func NewRegexGenerator(regex string) *RegexGenerator {
	return &RegexGenerator{regex: regex}
}

func (r *RegexGenerator) Generate() string {
	regex := r.regex
	regex = regexp.MustCompile(`^/?\^?`).ReplaceAllString(regex, "")
	regex = regexp.MustCompile(`\$?/?$`).ReplaceAllString(regex, "")
	regex = regexp.MustCompile(`{(\d+)}`).ReplaceAllString(regex, "{$1,$1}")
	regex = r.replaceCustomMarks('?', "{0,1}", regex)
	regex = r.replaceCustomMarks('*', "{0,"+strconv.Itoa(r.randomDigitNotNull())+"}", regex)
	regex = r.replaceCustomMarks('+', "{1,"+strconv.Itoa(r.randomDigitNotNull())+"}", regex)
	regex = regexp.MustCompile(`(\[[^\]]+\])\{(\d+),(\d+)\}`).ReplaceAllStringFunc(regex, func(m string) string {
		parts := regexp.MustCompile(`(\[[^\]]+\])\{(\d+),(\d+)\}`).FindStringSubmatch(m)
		min, max := parts[2], parts[3]
		repeated := r.repeatString(parts[1], r.randomIntElement(r.rangeSlice(min, max)))
		return repeated
	})
	regex = regexp.MustCompile(`(\([^\)]+\))\{(\d+),(\d+)\}`).ReplaceAllStringFunc(regex, func(m string) string {
		parts := regexp.MustCompile(`(\([^\)]+\))\{(\d+),(\d+)\}`).FindStringSubmatch(m)
		min, max := parts[2], parts[3]
		repeated := r.repeatString(parts[1], r.randomIntElement(r.rangeSlice(min, max)))
		return repeated
	})
	regex = regexp.MustCompile(`(\\?.)\{(\d+),(\d+)\}`).ReplaceAllStringFunc(regex, func(m string) string {
		parts := regexp.MustCompile(`(\\?.)\{(\d+),(\d+)\}`).FindStringSubmatch(m)
		min, max := parts[2], parts[3]
		repeated := r.repeatString(parts[1], r.randomIntElement(r.rangeSlice(min, max)))
		return repeated
	})
	regex = regexp.MustCompile(`\((.*?)\)`).ReplaceAllStringFunc(regex, func(m string) string {
		parts := regexp.MustCompile(`\((.*?)\)`).FindStringSubmatch(m)
		elements := strings.Split(strings.ReplaceAll(parts[1], "(", ")"), "|")
		randomElement := r.randomElement(elements)
		return randomElement
	})
	regex = regexp.MustCompile(`\[([^\]]+)\]`).ReplaceAllStringFunc(regex, func(match string) string {
		inner := match[1 : len(match)-1]
		innerRe := regexp.MustCompile(`(\w|\d)\-(\w|\d)`)
		expandedInner := innerRe.ReplaceAllStringFunc(inner, func(submatch string) string {
			rangeParts := strings.Split(submatch, "-")
			if len(rangeParts) == 2 {
				start := rangeParts[0]
				end := rangeParts[1]
				return r.expandRange(start[0], end[0])
			}
			return submatch
		})
		return "[" + expandedInner + "]"
	})
	regex = regexp.MustCompile(`\[([^\]]+)\]`).ReplaceAllStringFunc(regex, func(match string) string {
		inner := match[1 : len(match)-1]
		elements := strings.Split(inner, "")
		randomIndex := rand.Intn(len(elements))
		return elements[randomIndex]
	})
	regex = regexp.MustCompile(`\\w`).ReplaceAllStringFunc(regex, func(m string) string {
		return r.randomLetter()
	})
	regex = regexp.MustCompile(`\\d`).ReplaceAllStringFunc(regex, func(m string) string {
		return strconv.Itoa(r.randomDigit())
	})
	regex = r.replaceCustomMarks('.', r.asciify("*"), regex)
	regex = strings.ReplaceAll(regex, "\\", "")

	return regex
}

func (r *RegexGenerator) randomDigitNotNull() int {
	return rand.Intn(9) + 1
}

func (r *RegexGenerator) rangeSlice(start, end string) []int {
	min, _ := strconv.Atoi(start)
	max, _ := strconv.Atoi(end)
	return r.rangeIntSlice(min, max)
}

func (r *RegexGenerator) rangeIntSlice(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func (r *RegexGenerator) randomIntElement(array []int) int {
	rand.Seed(time.Now().UnixNano())
	return array[rand.Intn(len(array))]
}

func (r *RegexGenerator) randomElement(array []string) string {
	rand.Seed(time.Now().UnixNano())
	return array[rand.Intn(len(array))]
}

func (r *RegexGenerator) repeatString(s string, times int) string {
	return strings.Repeat(s, times)
}

func (r *RegexGenerator) replaceRangeWithChars(start, end string) string {
	startChar := r.parseInt(start)
	endChar := r.parseInt(end)
	var chars []string
	for i := startChar; i <= endChar; i++ {
		chars = append(chars, fmt.Sprint(i))
	}
	return strings.Join(chars, "")
}

func (r *RegexGenerator) parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func (r *RegexGenerator) randomLetter() string {
	return string(rune(rand.Intn(122-97+1) + 97))
}

func (r *RegexGenerator) asciify(s string) string {
	return strings.ReplaceAll(s, "*", r.randomAscii())
}

func (r *RegexGenerator) randomAscii() string {
	return fmt.Sprint(rune(rand.Intn(126-33+1) + 33))
}

func (r *RegexGenerator) randomDigit() int {
	return rand.Intn(10)
}

func (r *RegexGenerator) replaceCustomMarks(pattern rune, replaceWith, input string) string {
	var result strings.Builder
	escaped := false

	for _, char := range input {
		if char == pattern && !escaped {
			result.WriteString(replaceWith)
		} else {
			if char == '\\' {
				escaped = !escaped
			} else {
				escaped = false
			}
			result.WriteRune(char)
		}
	}

	return result.String()
}

func (r *RegexGenerator) expandRange(start, end byte) string {
	var expanded []byte
	for i := start; i <= end; i++ {
		expanded = append(expanded, i)
	}
	return string(expanded)
}
