package models

type TestResult struct {
	Wrong int
	Right int
}

type TestPageData struct {
	Topic       string
	Words       []*Word
	Result      *TestResult
	TestPassed  bool
	LearnPassed bool
}
