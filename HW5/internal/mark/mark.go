package mark

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Student struct {
	Name string
	Mark int
}

type StudentsStatistic interface {
	SummaryByStudent(student string) (int, bool)     // default_value, false - если студента нет
	AverageByStudent(student string) (float32, bool) // default_value, false - если студента нет
	Students() []string
	Summary() int
	Median() int
	MostFrequent() int
}

type StudentStatisticsImpl struct {
	StudentInfos []StudentInfo
}

type StudentInfo struct {
	name string
	mark int
}

func ReadStudentsStatistic(reader io.Reader) (StudentsStatistic, error) {

	regex := regexp.MustCompile("([ \t])")

	studentInfos := make([]StudentInfo, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		splitted := regex.Split(line, -1)
		if len(splitted) < 2 {
			continue
		} else {
			studentName := strings.Join(splitted[:len(splitted)-1], " ")

			studentMark := 0
			if mark, err := strconv.Atoi(splitted[len(splitted)-1]); err == nil {
				if mark < 10 {
					studentMark = mark
				} else {
					continue
				}
			} else {
				continue
			}

			studentInfos = append(studentInfos, StudentInfo{
				name: studentName,
				mark: studentMark,
			})
		}
	}

	return StudentStatisticsImpl{StudentInfos: studentInfos}, scanner.Err()
}

func WriteStudentsStatistic(writer io.Writer, statistic StudentsStatistic) error {

	//statisticImpl := statistic.(StudentStatisticsImpl)

	students := statistic.Students()

	comparator := func(i, j int) bool {
		iSum, _ := statistic.SummaryByStudent(students[i])
		jSum, _ := statistic.SummaryByStudent(students[j])
		return iSum-jSum > 0
	}

	sort.Slice(students, comparator)

	_, err := writer.Write([]byte(fmt.Sprintf("%d\t%d\t%d\n", statistic.Summary(), statistic.Median(), statistic.MostFrequent())))
	if err != nil {
		return err
	}

	for i, student := range students {
		sum, _ := statistic.SummaryByStudent(student)
		avg, _ := statistic.AverageByStudent(student)
		avgString := fmt.Sprintf("%.2f", avg)
		avgSplitted := strings.Split(avgString, ".")
		avgResult := "Почему нельзя выводить число с 2 знаками после запятой всегда :с"
		if avgSplitted[1] == "00" {
			avgResult = avgSplitted[0]
		} else {
			avgResult = avgString
		}

		_, err = writer.Write([]byte(fmt.Sprintf("%s\t%d\t%s", student, sum, avgResult)))

		if i != len(students)-1 {
			_, err = writer.Write([]byte("\n"))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s StudentStatisticsImpl) SummaryByStudent(student string) (int, bool) {
	sum, _, found := s.summaryAndEntries(student)
	return sum, found
}

func (s StudentStatisticsImpl) AverageByStudent(student string) (float32, bool) {
	sum, n, found := s.summaryAndEntries(student)
	avg := float64(sum) / float64(n)
	avgRounded := math.Round(avg*100) / 100
	return float32(avgRounded), found
}

func (s StudentStatisticsImpl) summaryAndEntries(student string) (int, int, bool) {
	sum := 0
	n := 0
	found := false
	for _, studentInfo := range s.StudentInfos {
		if studentInfo.name == student {
			found = true
			sum += studentInfo.mark
			n += 1
		}
	}
	return sum, n, found
}

func (s StudentStatisticsImpl) Students() []string {
	addedStudents := make(map[string]struct{}, len(s.StudentInfos))
	students := make([]string, 0)
	for _, info := range s.StudentInfos {
		if _, ok := addedStudents[info.name]; ok == false {
			addedStudents[info.name] = struct{}{}
			students = append(students, info.name)
		}
	}
	return students
}

func (s StudentStatisticsImpl) Summary() int {
	sum := 0
	for _, studentInfo := range s.StudentInfos {
		sum += studentInfo.mark
	}
	return sum
}

func (s StudentStatisticsImpl) Median() int {
	marks := make([]int, 0)
	for _, info := range s.StudentInfos {
		marks = append(marks, info.mark)
	}
	sort.Ints(marks)

	return marks[(len(marks) / 2)]
}

func (s StudentStatisticsImpl) MostFrequent() int {
	studentMap := make(map[int]int)
	for _, studentInfo := range s.StudentInfos {
		_, ok := studentMap[studentInfo.mark]
		if ok {
			studentMap[studentInfo.mark] = studentMap[studentInfo.mark] + 1
		} else {
			studentMap[studentInfo.mark] = 1
		}
	}
	maxM := -1
	maxN := 0
	for mark, n := range studentMap {
		if n > maxN || (n == maxN && mark > maxM) {
			maxM = mark
			maxN = n
		}
	}
	return maxM
}
