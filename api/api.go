//authors: marceloalvarez39 and emv18
package api

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/marceloalvarez39/dc-final/controller"
)

type Workload struct {
	WorkloadID     string   //"workload_id"
	Filter         string   //"filter"
	WorkloadName   string   //"workload_name"
	Status         string   //"status"
	RunningJobs    int      //"running_jobs"
	FilteredImages []string //"filtered_images"
}

type Images struct {
	WorkloadID string
	ImageID    string
	ImageType  string
}

var Users = make(map[string]string)
var PASSWORD = "123"

func GetLogin(context *gin.Context) {
	user, password, hasAuth := context.Request.BasicAuth()
	if password == PASSWORD && hasAuth {
		token := generateToken(user)
		Users[token] = user
		context.JSON(200, gin.H{
			"message": "Hi " + user + ", welcome to the DPIP System.",
			"token":   token,
		})
	} else {
		context.Abort()
	}
}

func GetLogout(context *gin.Context) {
	token := getToken(context.Request.Header.Get("Authorization"))
	if user, exists := Users[token]; exists {
		context.JSON(200, gin.H{
			"message": "Bye " + user + ", your token has been revoked. We hope you have an excellent day. ",
		})
		delete(Users, user)
	} else {
		context.JSON(http.StatusConflict, gin.H{
			"message": "Your Token does not exist",
		})
	}
}

// UploadImage gets an image from a path
// and reads it in order to change it later.
// Each image has an id and is related to a
// worker id.

func UploadImage(context *gin.Context) {
	token := getToken(context.Request.Header.Get("Authorization"))

	if _, exists := Users[token]; exists {
		var image Images
		file, err := context.FormFile("data")
		if err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("get form error: %s", err.Error()))
			return
		}

		random := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
		image.ImageID = strconv.Itoa(random)
		image.ImageType = "original"

		if err := context.SaveUploadedFile(file, "image"+image.ImageID+".jpg"); err != nil {
			context.String(http.StatusBadRequest, fmt.Sprintf("upload file error: %s", err.Error()))
			return
		}
		fileName := file.Filename
		fileSize := file.Size
		context.JSON(200, gin.H{
			"message":  "An image has been successfully uploaded",
			"type":     image.ImageType,
			"filename": fileName,
			"size":     fileSize,
		})
	} else {
		context.Abort()
	}

}

func getImage_Download(context *gin.Context) {
	//downloads en image
}

func getWorkloadData(context *gin.Context) {
	//gets data of a workload
}

func GetStatus(context *gin.Context) {
	token := getToken(context.Request.Header.Get("Authorization"))
	if user, exists := Users[token]; exists {
		context.JSON(200, gin.H{
			"User":             "Hello " + user + ".",
			"System Name":      "Distributed Parallel Image Processing (DPIP) System",
			"Server Time":      time.Now().Format("2006-01-02 3:4:5"),
			"Active Workloads": "You have " + strconv.Itoa(len(controller.Workloads)) + " workloads active.",
		})
	} else {
		context.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  false,
			"message": "No username registered with the given token. Please check your token and try again or log in",
		})
	}
}

func MakeWorkloads(context *gin.Context) {
	token := getToken(context.Request.Header.Get("Authorization"))
	if _, exists := Users[token]; exists {
		var work Workload

		random := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
		work.WorkloadID = strconv.Itoa(random)
		work.WorkloadName = context.PostForm("workload_name")
		work.Filter = context.PostForm("filter")

		// Check if the name given by the user is already used.
		available := true
		for _, w := range controller.Workloads {
			if w.WorkloadName == work.WorkloadName {
				available = false
				break
			}
		}

		// If it is available, then we can create the workload
		if available {
			if len(controller.Workers) > 0 {
				work.Status = "running"
			} else {
				work.Status = "scheduling"
			}

			// In order to have certain order in our API, we decided to create folders for the images.
			// Processed Folder, for the images that have been processed by a filter
			processedFolder := "images/processed/" + work.WorkloadName + "/"
			_ = os.MkdirAll(processedFolder, 0755)

			// Not Processed Folder, for the images that haven't been processed by a filter
			notProcessedFolder := "images/notProcessed/" + work.WorkloadName + "/"
			_ = os.MkdirAll(notProcessedFolder, 0755)

			newWorkload := controller.Workload{
				WorkloadID:     work.WorkloadID,
				Filter:         work.Filter,
				WorkloadName:   work.WorkloadName,
				Status:         work.Status,
				RunningJobs:    0,
				FilteredImages: []string{},
			}

			controller.Workloads[fmt.Sprintf("%v", newWorkload.WorkloadID)] = newWorkload

			context.JSON(http.StatusOK, gin.H{
				"workload_id":     newWorkload.WorkloadID,
				"filter":          newWorkload.Filter,
				"workload_name":   newWorkload.WorkloadName,
				"status":          newWorkload.Status,
				"running_jobs":    newWorkload.RunningJobs,
				"filtered_images": newWorkload.FilteredImages,
			})
		} else {
			context.JSON(http.StatusConflict, gin.H{
				"message": "This workload already exists",
			})
		}
	}
}

// ----------------------------------- Tools ---------------------------------
// -----------------------------------------------------------------------------

func getToken(header string) string {
	re := regexp.MustCompile("<(.*?)>")
	match := re.Find([]byte(header))
	token := string(match[1 : len(match)-1])
	return token
}

func generateToken(object interface{}) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", object)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// -------- Start Function ---------------------

func Start() {
	router := gin.Default()
	router.POST("/login", GetLogin)
	router.DELETE("/logout", GetLogout)
	router.POST("/images", UploadImage)
	router.GET("/images/{image_id}", getImage_Download)
	router.GET("/status", GetStatus)
	router.POST("/workloads", MakeWorkloads)
	router.GET("workloads/{workloads_id}", getWorkloadData)
	router.Run(":8080")
}
