package service

import (
	"github.com/MuxiKeStack/muxiK-StackBackend/model"
)

// 课堂信息
type SearchCourseInfo struct {
	Id         uint32  //主键
	Name       string  //课程名称
	Credit     float32 //学分
	Teacher    string  //任课教师姓名
	CourseId   string  //课程编号
	ClassId    uint64  //课堂编号
	Type       uint8   //课程类型
	CreditType uint8   //学分类型
}

func SearchCourses(keyword string, page, limit uint64) ([]SearchCourseInfo, error) {
	courses := make([]SearchCourseInfo, 0)
	courseRows, err := model.AgainstAndMatchCourses(keyword, page, limit)
	if err != nil {
		return courses, err
	}
	defer courseRows.Close()

	for courseRows.Next() {
		var course SearchCourseInfo
		courseRows.Scan(&course)
		courses = append(courses, course)
	}
	return courses, nil
}