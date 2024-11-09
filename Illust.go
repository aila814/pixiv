package pixiv

import (
	"errors"
	"fmt"
	"github.com/aila814/pixiv/filea"
	"github.com/tidwall/gjson"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// 获取插图
func (p *Pixiv) GetIllust(IllustID string) (Illust, Error) {
	var (
		Err Error
		i   Illust
	)
	Url := fmt.Sprintf("%s/v1/illust/detail?filter=for_android&illust_id=%s", ApiAddress, IllustID)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		Err.Err = err
		return i, Err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		Err.Err = err
		return i, Err
	}
	body := GetBody(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		Err.Code = resp.StatusCode
		Err.Body = body
		if gjson.Get(body, "error.message").String() != "" {
			Err.Err = errors.New(gjson.Get(body, "error.message").String())
		} else {
			Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
		}
		return i, Err
	}
	i.UserID = gjson.Get(body, "illust.user.id").String()
	i.ID = gjson.Get(body, "illust.id").String()
	i.Title = gjson.Get(body, "illust.title").String()
	i.Type = gjson.Get(body, "illust.type").String()
	i.Width = gjson.Get(body, "illust.width").Int()
	i.Height = gjson.Get(body, "illust.height").Int()
	i.Caption = gjson.Get(body, "illust.caption").String()
	pageUrl := gjson.Get(body, "illust.meta_single_page.original_image_url").String()
	i.Page = map[string]string{path.Base(pageUrl): pageUrl}
	gjson.Get(body, "illust.meta_pages").ForEach(func(key, value gjson.Result) bool {
		ImgUrl := value.Get("image_urls.original").String()
		i.Pages = append(i.Pages, map[string]string{path.Base(ImgUrl): ImgUrl})

		return true
	})
	Err.Code = resp.StatusCode
	Err.Body = body
	return i, Err
}

// 获取动图
func (p *Pixiv) GetIllustGif(IllustID string) (IllustGif, Error) {
	var (
		Err Error
		i   IllustGif
	)
	Url := fmt.Sprintf("%s/v1/ugoira/metadata?illust_id=%s", ApiAddress, IllustID)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		Err.Err = err
		return i, Err
	}
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234 ")
	req.Header.Set("Authorization", "Bearer "+p.AccessToken)

	resp, err := p.HttpClient.Do(req)
	if err != nil {
		Err.Err = err
		return i, Err
	}
	body := GetBody(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		Err.Code = resp.StatusCode
		Err.Body = body
		if gjson.Get(body, "error.message").String() != "" {
			Err.Err = errors.New(gjson.Get(body, "error.message").String())
		} else {
			Err.Err = errors.New(gjson.Get(body, "error.user_message").String())
		}
		return i, Err
	}
	i.ZipUrl = gjson.Get(body, "ugoira_metadata.zip_urls.medium").String()
	gjson.Get(body, "ugoira_metadata.frames").ForEach(func(key, value gjson.Result) bool {
		name := value.Get("file").String()
		delay := value.Get("delay").Int()
		i.Pages = append(i.Pages, map[string]int64{name: delay})
		return true
	})
	Err.Code = resp.StatusCode
	Err.Body = body
	return i, Err
}

// 获取插图数据
func (p *Pixiv) GetIllustByte(Url string) ([]byte, error) {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "i.pximg.net")
	req.Header.Set("referer", "https://app-api.pixiv.net/")
	req.Header.Set("accept-language", "zh_CN")
	req.Header.Set("User-Agent", "PixivAndroidApp/5.0.234")
	resp, err := p.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("服务器错误响应: %s,code: %d", resp.Status, resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// 图片转gif ·=·
func (p *Pixiv) ToGif(dir, outputFile string, data []map[string]int64) error {
	var frames []*image.Paletted
	var delays []int
	// 生成调色板
	palette := generatePalette(dir, data)

	var cmd []string
	var cmd2 []string
	cmd = append(cmd, "convert")
	cmd2 = append(cmd2, "convert")
	for _, m := range data {
		for filename, delay := range m {

			cmd = append(cmd, fmt.Sprintf("\\( %s -delay %d \\)", fmt.Sprintf("./%s/%s", dir, filename), int(delay/10)))
			cmd2 = append(cmd2, fmt.Sprintf("-delay %d %s ", int(delay/10), fmt.Sprintf("./%s/%s", dir, filename)))

			img, format, err := loadImage(fmt.Sprintf("./%s/%s", dir, filename))
			if err != nil {
				return fmt.Errorf("加载图片错误 %s: %v", filename, err)
			}
			fmt.Printf("已加载图像 %s 格式 %s\n", filename, format)

			// 转换图像为 Paletted 格式
			palettedImg := image.NewPaletted(img.Bounds(), palette)
			// 将图像绘制到 Paletted 图像上
			draw.Draw(palettedImg, palettedImg.Rect, img, image.Point{}, draw.Over)

			//================================
			//// 使用 image/palette.Plan9 来获得一个标准调色板
			//palettedImg := image.NewPaletted(img.Bounds(), palette.Plan9)
			// 将图像绘制到 Paletted 图像上
			//draw.Draw(palettedImg, palettedImg.Rect, img, image.Point{}, draw.Over)
			//draw.FloydSteinberg.Draw(palettedImg, img.Bounds(), img, image.Point{})
			//==========================
			// Add the image to the frames
			frames = append(frames, palettedImg)
			// Add the delay to the delays slice
			delays = append(delays, int(delay/10))

			//fmt.Printf("  %s: %d\n", key, value)
		}
	}
	cmd = append(cmd, "dir.gif")
	cmd2 = append(cmd2, "dir.gif")
	filea.WriteFile("./cmd.txt", strings.Join(cmd, " "))
	filea.WriteFile("./cmd2.txt", strings.Join(cmd2, " "))
	// 创建一个 GIF 文件
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()
	// Create the GIF with the frames and delays
	err = gif.EncodeAll(f, &gif.GIF{
		Image:     frames,
		Delay:     delays,
		LoopCount: 0, // Loop forever
	})
	if err != nil {
		return fmt.Errorf("错误编码 GIF: %v", err)

	}
	return nil
}

// 从文件加载图像
func loadImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", fmt.Errorf("无法打开图像文件: %v", err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)

	if err != nil {
		return nil, "", fmt.Errorf("无法解码图像: %v", err)
	}

	return img, "jpeg", nil
}

// generatePalette 生成调色板
func generatePalette(dir string, data []map[string]int64) color.Palette {
	colorMap := make(map[color.Color]struct{})

	for _, m := range data {
		for filename, _ := range m {
			img, _, err := loadImage(fmt.Sprintf("./%s/%s", dir, filename))
			if err != nil {
				panic(err) // 在实际应用中，你可能想要更好的错误处理
			}

			// 遍历图像的每一个像素，收集颜色
			bounds := img.Bounds()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					colora := img.At(x, y)
					colorMap[colora] = struct{}{}
				}
			}
		}
	}

	// 将颜色转换为调色板
	palette := make(color.Palette, 0, len(colorMap))
	for c := range colorMap {
		palette = append(palette, c)
	}

	// 限制调色板的大小，最多256种颜色
	if len(palette) > 256 {
		palette = palette[:256]
	}

	return palette
}
