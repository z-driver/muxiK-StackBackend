package service

import (
	"strings"

	"github.com/MuxiKeStack/muxiK-StackBackend/model"
	"github.com/MuxiKeStack/muxiK-StackBackend/util"

	"github.com/lexkong/log"
)

type ProducedCourseItem struct {
	CourseId string `json:"course_id"`
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	// Academic     string `json:"academic"`
	HasEvaluated bool `json:"has_evaluated"`
}

// Get one's all courses from XK.
func GetSelfCourseList(userId uint32, sid, pwd, year, term string) (*[]ProducedCourseItem, error) {
	originalCourses, err := util.GetSelfCoursesFromXK(sid, pwd, year, term)
	if err != nil {
		log.Error("GetSelfCoursesFromXK function error", err)
		return nil, err
	}

	var list []ProducedCourseItem
	//(*originalCourses.Items)[0].Jsxx = "2006982627/葛非,2006982646/彭熙,2006982670/刘明,2007980066/姚华雄"

	for _, item := range *originalCourses.Items {
		teacher := util.GetTeachersSqStrBySplitting(item.Jsxx)
		hashId := util.HashCourseId(item.Kch, teacher)
		info := ProducedCourseItem{
			CourseId: hashId,
			Name:     item.Kcmc,
			Teacher:  teacher,
			// Academic:     item.Kkxymc,
			HasEvaluated: model.HasEvaluated(userId, hashId),
		}
		list = append(list, info)
	}

	return &list, nil
}

func GetSelfCourseListFromLocal(userId uint32) (*[]ProducedCourseItem, error) {
	hashIdStr, err := model.GetSelfCoursesByUserId(userId)
	if err != nil {

		return nil, err
	}

	hashIds := strings.Split(hashIdStr, ",")
	var list []ProducedCourseItem

	for _, hashId := range hashIds {
		course := &model.UsingCourseModel{Hash: hashId}
		if err := course.GetByHash(); err != nil {
			return nil, err
		}
		item := ProducedCourseItem{
			CourseId:     hashId,
			Name:         course.Name,
			Teacher:      course.Teacher,
			HasEvaluated: model.HasEvaluated(userId, hashId),
		}
		list = append(list, item)
	}

	return &list, nil
}

func SavingCourseDataToLocal(userId uint32, list *[]ProducedCourseItem) error {
	var record = &model.SelfCourseModel{UserId: userId}
	ok, err := record.GetByUserId()
	if err != nil {
		return err
	}
	// 无记录则新添
	var hashIds []string
	for _, item := range *list {
		hashIds = append(hashIds, item.CourseId)
	}
	hashIdStr := strings.Join(hashIds, ",")

	// 不存在记录则新添记录
	if !ok {
		record.Courses = hashIdStr
		err = record.New()
	} else if len(hashIdStr) != len(record.Courses) {
		// 若存在且课程变化则更新
		record.Courses = hashIdStr
		err = record.Update()
	}

	if err != nil {
		return err
	}
	return nil
}
