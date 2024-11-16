package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/bbrks/go-blurhash"
	"github.com/go-chi/chi/v5"
	"github.com/goccy/go-json"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"openchat/dto"
	"openchat/lib"
	"openchat/model"
	"openchat/utils/rcontext"
	"openchat/utils/xhttp"
)

var Image = &imageApi{}

type imageApi struct{}

func (imageApi) ListImages(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	images := make([]*model.Image, 0)
	if err := lib.DB.Where("user_id = ?", userId).Order("created_at DESC").Find(&images).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		presigner := s3.NewPresignClient(lib.S3)
		presignHttpRequest, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(image.StorageKey),
		}, func(options *s3.PresignOptions) {
			options.Expires = time.Hour
		})
		if err != nil {
			continue
		}
		presignedUrl, _ := url.Parse(presignHttpRequest.URL)
		image.Url = strings.ReplaceAll(presignHttpRequest.URL, presignedUrl.Host, os.Getenv("CDN_HOST"))
	}
	xhttp.RespJson(w, images, http.StatusOK)
}

func (imageApi) DeleteImage(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	imageId := chi.URLParam(r, "id")
	if err := lib.DB.Where("user_id = ? AND id = ?", userId, imageId).Delete(model.Image{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (imageApi) GenerateImageFromCf(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	var reqData dto.GenerateImageReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	prompt := reqData.Prompt
	cfAIGateway := os.Getenv("CF_AI_GATEWAY")
	modelId := "@cf/black-forest-labs/flux-1-schnell"
	data := map[string]interface{}{
		"prompt": prompt,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req, err := http.NewRequest("POST", cfAIGateway+"/"+modelId, strings.NewReader(string(jsonData)))
	token := os.Getenv("CF_API_TOKEN")
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var cfImageResponse dto.CfImageResponse
	if err := json.Unmarshal(body, &cfImageResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	imgData, err := base64.StdEncoding.DecodeString(cfImageResponse.Result.Image)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	imageFormat := "png"
	img, err := png.Decode(bytes.NewReader(imgData))
	if err != nil {
		img, err = jpeg.Decode(bytes.NewReader(imgData))
		if err != nil {
			log.Println("err: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageFormat = "jpeg"
	}
	filename := gonanoid.Must(8) + "." + imageFormat
	_, err = lib.S3.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(imgData),
		ContentType: aws.String("image/" + imageFormat),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hash, err := blurhash.Encode(4, 4, img)
	if err != nil {
		log.Println("err: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	image := model.Image{
		UserId:     userId,
		StorageKey: filename,
		Blurhash:   hash,
		Prompt:     reqData.Prompt,
	}
	if err = lib.DB.Create(&image).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	presigner := s3.NewPresignClient(lib.S3)
	presignHttpRequest, err := presigner.PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(filename),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Hour
	})
	if err != nil {
		return
	}
	presignedUrl, _ := url.Parse(presignHttpRequest.URL)
	image.Url = strings.ReplaceAll(presignHttpRequest.URL, presignedUrl.Host, os.Getenv("CDN_HOST"))
	xhttp.RespJson(w, image, http.StatusCreated)
}
