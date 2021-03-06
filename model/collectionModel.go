package model

type CourseListModel struct {
	Id           uint32 `gorm:"column:id"`
	UserId       uint32 `gorm:"column:user_id"`
	CourseHashId string `gorm:"column:course_hash_id"`
}

// 增加Type，用于前端颜色分配
// 课表界面的课程清单信息
type CourseInfoInTableCollection struct {
	CourseId   string                    `json:"course_id"` // 课程hash id
	CourseName string                    `json:"course_name"`
	ClassSum   int                       `json:"class_sum"` // 课堂数
	Type       int8                      `json:"type"`      // 0-通必,1-专必,2-专选,3-通选,4-专业课,5-通核
	Classes    []*ClassInfoInCollections `json:"classes"`
}

// 选课清单内的课堂（教学班）信息
type ClassInfoInCollections struct {
	ClassId   string                        `json:"class_id"` // 教学班编号
	ClassName string                        `json:"class_name"`
	Teacher   string                        `json:"teacher"`
	Times     []*ClassTimeInfoInCollections `json:"times"`
	Places    []string                      `json:"places"`
}

type ClassTimeInfoInCollections struct {
	Time      string `json:"time"`       // 时间区间（节数），1-2
	Day       int8   `json:"day"`        // 星期几
	Weeks     string `json:"weeks"`      // 周次，2-17
	WeekState int8   `json:"week_state"` // 全周0,单周1,双周2
}

// 选课清单页面的课程信息
type CourseInfoForCollections struct {
	Id                  uint32    `json:"id"` // 数据库表中记录的id，自增id
	CourseId            string    `json:"course_id"`
	CourseName          string    `json:"course_name"`
	Teacher             string    `json:"teacher"`
	EvaluationNum       uint32    `json:"evaluation_num"`
	Rate                float32   `json:"rate"`
	AttendanceCheckType string    `json:"attendance_check_type"`
	ExamCheckType       string    `json:"exam_check_type"`
	Tags                *[]string `json:"tags"`
}
