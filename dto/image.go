package dto

type GenerateImageReqData struct {
	Prompt string `json:"prompt"`
}

type CfImageResponse struct {
	Result ImageResult `json:"result"`
}

type ImageResult struct {
	Image string `json:"image"`
}
