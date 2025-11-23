package models

// Struktur Output Statistik Global
type GlobalStatistics struct {
    TotalAchievements int                    `json:"totalAchievements"`
    PointsDistribution []TopStudent          `json:"topStudents"`
    TypeDistribution   map[string]int        `json:"typeDistribution"`
    LevelDistribution  map[string]int        `json:"levelDistribution"` // Nasional, Internasional, dll
    TrendByYear        map[string]int        `json:"trendByYear"`
}

type TopStudent struct {
    StudentID   string `json:"studentId"`
    Name        string `json:"name"`
    ProgramStudy string `json:"programStudy"`
    TotalPoints int    `json:"totalPoints"`
}

// Struktur Output Statistik Per Mahasiswa
type StudentStatistics struct {
    StudentName      string         `json:"studentName"`
    TotalPoints      int            `json:"totalPoints"`
    TotalAchievements int           `json:"totalAchievements"`
    ByType           map[string]int `json:"byType"`
}