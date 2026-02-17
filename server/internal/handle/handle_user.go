package handle

import (
	g "server/internal/global"
	"server/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminGetAllUsers 获取所有用户列表 (管理员)
func (*UserAuth) AdminGetAllUsers(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	if currentUser == nil || currentUser.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	query := c.Query("query")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	db := GetDB(c)
	users, total, err := model.GetAllUsers(db, page, limit, query)
	if err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	// 转换为VO，去除敏感信息
	var userVOs []LoginVO
	for _, u := range users {
		userVOs = append(userVOs, LoginVO{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			UserGroup: u.UserGroup,
			Avatar:    u.Avatar,
			LastLogin: u.LastLogin,
		})
	}

	ReturnSuccess(c, gin.H{
		"list":  userVOs,
		"total": total,
	})
}

// AdminDeleteUser 删除指定用户 (管理员)
func (*UserAuth) AdminDeleteUser(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	if currentUser == nil || currentUser.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 不能删除自己
	if id == currentUser.ID {
		ReturnError(c, g.ErrRequest, nil)
		return
	}

	db := GetDB(c)
	if err := model.DeleteUser(db, id); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

type UpdateUserRoleReq struct {
	UserGroup string `json:"user_group" binding:"required"`
}

// AdminUpdateUserRole 修改用户权限 (管理员)
func (*UserAuth) AdminUpdateUserRole(c *gin.Context) {
	currentUser := GetCurrentUser(c)
	if currentUser == nil || currentUser.UserGroup != "admin" {
		ReturnError(c, g.ErrPermission, nil)
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	var req UpdateUserRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, g.ErrRequest, err)
		return
	}

	// 不能修改自己的权限
	if id == currentUser.ID {
		ReturnError(c, g.ErrRequest, nil)
		return
	}

	// 验证用户组是否合法
	validGroups := map[string]bool{"admin": true, "user": true, "guest": true}
	if !validGroups[req.UserGroup] {
		ReturnError(c, g.ErrRequest, nil)
		return
	}

	db := GetDB(c)
	if err := model.UpdateUserGroup(db, id, req.UserGroup); err != nil {
		ReturnError(c, g.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}
