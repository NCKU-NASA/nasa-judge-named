package router
import (
    "github.com/gin-gonic/gin"
    
    "github.com/NCKU-NASA/nasa-judge-lib/schema/user"

    "github.com/NCKU-NASA/nasa-judge-named/middlewares/auth"
    "github.com/NCKU-NASA/nasa-judge-named/utils/errutil"
    "github.com/NCKU-NASA/nasa-judge-named/models/named"
)

var router *gin.RouterGroup

func Init(r *gin.RouterGroup) {
    router = r
    router.POST("/update", auth.CheckIsTrust, update)
    router.POST("/set", auth.CheckIsTrust, set)
}

func update(c *gin.Context) {
    var userdata user.User
    err := c.ShouldBindJSON(&userdata)
    if err != nil {
        errutil.AbortAndStatus(c, 400)
        return
    }
    userdata = user.User{
        Username: userdata.Username,
    }
    userdata.Fix()
    if userdata.Username == "" {
        errutil.AbortAndStatus(c, 400)
    }
    userdata, err = user.GetUser(userdata)
    if err != nil {
        errutil.AbortAndError(c, &errutil.Err{
            Code: 409,
            Msg: "username not exist",
        })
        return
    }
    named.SetRecord(userdata)
    c.String(200, "Success")
}

func set(c *gin.Context) {
    var record named.Record
    err := c.ShouldBindJSON(&record)
    if err != nil {
        errutil.AbortAndStatus(c, 400)
        return
    }
    record.Set()
    c.String(200, "Success")
}

