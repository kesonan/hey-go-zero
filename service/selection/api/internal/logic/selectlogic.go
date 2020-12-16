package logic

import (
	"context"
	"fmt"
	"time"

	"hey-go-zero/common/errorx"
	"hey-go-zero/service/course/rpc/course"
	"hey-go-zero/service/course/rpc/courseservice"
	"hey-go-zero/service/selection/api/internal/svc"
	"hey-go-zero/service/selection/api/internal/types"
	"hey-go-zero/service/selection/model"
	"hey-go-zero/service/user/rpc/userservice"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SelectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectLogic(ctx context.Context, svcCtx *svc.ServiceContext) SelectLogic {
	return SelectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectLogic) Select(userId int64, req types.SelectCourseId) error {
	if req.Id <= 0 {
		return errorx.NewInvalidParameterError("id")
	}

	key := fmt.Sprintf("%v", req.Id)
	lock := redis.NewRedisLock(l.svcCtx.BizRedis, key)
	lock.SetExpire(3)
	ok, err := lock.Acquire()
	if err != nil {
		return err
	}

	// todo: 这里和中间件中有重复逻辑，都查询了一次user，开发人员可以自己优化一下。
	userInfo, err := l.svcCtx.UserService.FindOne(l.ctx, &userservice.UserReq{Id: req.Id})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return errorx.NewDescriptionError("用户不存在")
		}

		return errorx.NewDescriptionError(st.Message())
	}

	now := time.Now()
	selectionCourse, err := l.svcCtx.SelectionCourseModel.FindOne(req.Id)
	switch err {
	case nil:
	case model.ErrNotFound:
		return errSelectionCourseNotFound
	default:
		return err
	}

	selection, err := l.svcCtx.SelectionModel.FindOne(selectionCourse.SelectionId)
	switch err {
	case nil:
		if now.Before(time.Unix(selection.StartTime, 0)) {
			return errNotInSelectionTimeRange
		}

		if now.After(time.Unix(selection.EndTime, 0)) {
			return errSelectionExpired
		}

	case model.ErrNotFound:
		return errSelectionNotFound
	default:
		return err
	}

	courseInfo, err := l.svcCtx.CourseService.FindOne(l.ctx, &courseservice.IdReq{Id: selectionCourse.CourseId})
	if err != nil {
		st := status.Convert(err)
		if st.Code() == codes.NotFound {
			return errCourseNotFound
		}

		return errorx.NewDescriptionError(st.Message())
	}

	if courseInfo.GenderLimit != course.GenderLimit_NoLimit && courseInfo.GenderLimit != course.GenderLimit(userInfo.Gender) {
		return errorx.NewDescriptionError("性别不符合")
	}

	if !ok {
		return errorx.NewDescriptionError("当前选课人数较多，请稍后再试")
	}
	defer func() {
		_, err = lock.Release()
		logx.Error(err)
	}()

	endTime := time.Unix(selection.EndTime, 0)
	ok, err = l.trySelect(req.Id, userId, selection.MaxCredit, courseInfo, endTime)
	if err != nil {
		logx.Error(err)
		return errorx.NewDescriptionError("选课失败，请稍后再试")
	}

	if !ok {
		return errorx.NewDescriptionError("选课人数已满，请选择其他课程")
	}
	threading.GoSafe(func() {
		_, err = l.svcCtx.SelectionStudentModel.Insert(model.SelectionStudent{
			SelectionCourseId: req.Id,
			StudentId:         userId,
		})
		if err != nil {
			logx.Error(err)
		}
	})

	return nil
}

func (l *SelectLogic) trySelect(selectCourseId, userId, maxCredit int64, courseInfo *courseservice.Course, expireAt time.Time) (bool, error) {
	expire := int(expireAt.Sub(time.Now()).Seconds()) + 1
	userKey := fmt.Sprintf("biz#user#selected#status#%v#%v", userId, selectCourseId)
	ok, err := l.svcCtx.BizRedis.SetnxEx(userKey, "*", expire)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, errSelectionSelected
	}

	userCreditKey := fmt.Sprintf("biz#user#selected#credit#%v", userId)
	credit, err := l.svcCtx.BizRedis.Incrby(userCreditKey, courseInfo.Credit)
	if err != nil {
		_, _ = l.svcCtx.BizRedis.Del(userKey)
		return false, err
	}

	if credit > maxCredit {
		_, _ = l.svcCtx.BizRedis.Del(userKey)
		return false, errorx.NewDescriptionError(fmt.Sprintf("选择当前课程后，该学期你的学分已经超出总学分%d，请合理选择课程", maxCredit))
	}
	_ = l.svcCtx.BizRedis.Expire(userCreditKey, expire)

	courseKey := fmt.Sprintf("biz#course#selected#data#%v", selectCourseId)
	var field string
	switch courseInfo.GenderLimit {
	case course.GenderLimit_Male:
		field = "male"
	case course.GenderLimit_Female:
		field = "female"
	default:
		field = "all"
	}

	count, err := l.svcCtx.BizRedis.Hincrby(courseKey, field, 1)
	if err != nil {
		_, _ = l.svcCtx.BizRedis.Del(userKey)
		return false, err
	}

	if count > int(courseInfo.MemberLimit) {
		_, _ = l.svcCtx.BizRedis.Hincrby(courseKey, field, -1)
		_, _ = l.svcCtx.BizRedis.Del(userKey)
		return false, nil
	}

	_ = l.svcCtx.BizRedis.Expire(courseKey, expire)
	return true, nil
}
