package main

import (
	"fmt"
	"log"
	"net/http"

	"bufio"
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	//"golang.org/x/text"

	"github.com/PuerkitoBio/goquery"
	"github.com/Unknwon/goconfig"

	"github.com/axgle/mahonia"
	"github.com/chromedp/chromedp"
)

var w_title = false

func set_ini_value(sec string, key string, value string) {
	cfg, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic("错误")
	}
	cfg.SetValue(sec, key, value)
	gerr := goconfig.SaveConfigFile(cfg, "conf/conf.ini")
	if gerr != nil {
		panic("错误")
	}
}
func get_ini_value(sec string, key string) (r string) {
	cfg, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic("错误")
	}
	value, err := cfg.GetValue(sec, key)
	if err != nil {
		panic(err.Error())
	}
	return value
}

func set_page(tp int, cp int) {
	cfg, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic("错误")
	}
	cfg.SetValue("progress", "totoalpage", strconv.Itoa(tp))
	cfg.SetValue("progress", "currentpage", strconv.Itoa(tp))
	gerr := goconfig.SaveConfigFile(cfg, "conf/conf.ini")
	if gerr != nil {
		panic("错误")
	}
}

func get_doc() (doc *goquery.Document, sec map[string]string) {
	cfg, err := goconfig.LoadConfigFile("conf/conf.ini")
	if err != nil {
		panic("错误")
	}
	sec_r, err := cfg.GetSection("url_point")

	res, err := http.Get(sec_r["url"])
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	dec := mahonia.NewDecoder("gbk")
	rd := dec.NewReader(res.Body)
	// Load the HTML document
	doc_r, err := goquery.NewDocumentFromReader(rd)
	if err != nil {
		log.Fatal(err)
	}
	return doc_r, sec_r
}

func get_page_locate(doc *goquery.Document, sec map[string]string) (cp int, tp int) {
	var currentpage, totalpage int
	var err error
	if sec["run"] == "false" {
		doc.Find(sec["totalpage"]).Find("div font").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				totalpage, err = strconv.Atoi(s.Text())
				if err != nil {
					log.Fatal(err)
				}
			} else {
				currentpage, err = strconv.Atoi(s.Text())
				if err != nil {
					log.Fatal(err)
				}
			}
		})
		set_ini_value("url_point", "run", "true")
	} else {
		currentpage, _ = strconv.Atoi(get_ini_value("progress", "currentpage"))
		totalpage, _ = strconv.Atoi(get_ini_value("progress", "totoalpage"))
	}
	return currentpage, totalpage
}

func get_data() {
	doc, sec := get_doc()
	cp, tp := get_page_locate(doc, sec)
	fmt.Printf("当前页%d/总页数%d\n", cp, tp)
	//show_mx1(doc, sec)
}

func show_mx1(str string, sec map[string]string) {
	var docmx *goquery.Document
	var docmx_err error
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(sec["tables"]).Find("tbody tr").Each(func(i int, s *goquery.Selection) {
		var title, data string
		if i == 0 && w_title == false {
			s.Find("td").Each(func(i0 int, s0 *goquery.Selection) {
				strtmp := s0.Text() + "\t"
				title += strtmp
			})
			title += "项目名称\t"
			title += "项目坐落\t"
			title += "开工时间\t"
			title += "竣工时间(预计)\t"
			title += "项目基本情况-用地面积\t"
			title += "项目基本情况-土地使用年限\t"
			title += "项目基本情况-土地用途\t"
			title += "项目基本情况-土地等级\t"
			title += "项目基本情况-建筑面积\t"
			title += "项目基本情况-容积率\t"
			title += "项目基本情况-房屋套数\t"
			title += "项目基本情况-房屋栋数\t"
			title += "项目基本情况-销售时间\t"
			title += "项目基本情况-他项权利情况\t"
			title += "项目证件情况-前期-建设用地规划许可证号\t"
			title += "项目证件情况-前期-国有土地使用证号\t"
			title += "项目证件情况-前期-建设工程规划许可证号\t"
			title += "项目证件情况-前期-施工许可证号\t"
			title += "项目证件情况-前期-商品房预售许可证号\t"
			title += "项目证件情况-前期-开发企业资质证号\t"
			title += "开发企业\t"
			title += "联系电话\t"
			title += "代理公司\t"
			title += "联系电话\t"
			title += "项目备案机关\t"
			w_title = true
			WirteTXT(title)
		}
		if i != 0 {
			col1_href, ie := s.Find("td a").Attr("href")
			if ie == true {
				data += sec["urlmx_head"] + col1_href + "\t"

				mx_url := convertUrlWithChineseToHex(sec["urlmx_head"] + col1_href)
				res_mx, res_mx_err := http.Get(mx_url)
				if res_mx_err != nil {
					// 错误处理
					log.Fatal(err)
				}
				defer res_mx.Body.Close()
				if res_mx.StatusCode != 200 {
					log.Fatalf("status code error: %d %s", res_mx.StatusCode, res_mx.Status)
				}

				dec_mx := mahonia.NewDecoder("gbk")
				rd_mx := dec_mx.NewReader(res_mx.Body)
				//rd_mx:=dec_mx.
				docmx, docmx_err = goquery.NewDocumentFromReader(rd_mx)
				if docmx_err != nil {
					log.Fatal(err)
				}

				//fmt.Printf("%s",docmx.Find("#txt_xmzl").Text())
			}
			s.Find("td").Each(func(ia int, sa *goquery.Selection) {
				data += sa.Text() + "\t"
			})

			data += docmx.Find("#txt_xmmc2").Text() + "\t"
			data += docmx.Find("#txt_xmzl").Text() + "\t"
			data += docmx.Find("#txt_kgsj").Text() + "\t"
			data += docmx.Find("#txt_jgsj").Text() + "\t"
			data += docmx.Find("#txt_ydmj").Text() + "\t"
			data += docmx.Find("#txt_tdsynx").Text() + "至"
			data += docmx.Find("#txt_tdsynx1").Text() + "\t"

			data += docmx.Find("#txt_tdyt").Text() + "\t"
			data += docmx.Find("#txt_tddj").Text() + "\t"
			data += docmx.Find("#txt_jzmj").Text() + "\t"
			data += docmx.Find("#txt_rjl").Text() + "\t"
			data += docmx.Find("#txt_fwts").Text() + "\t"
			data += docmx.Find("#txt_fwds").Text() + "\t"
			data += docmx.Find("#txt_xssj").Text() + "\t"
			data += docmx.Find("#txt_txqlqk").Text() + "\t"

			data += docmx.Find("#txt_jsydxkz").Text() + "\t"
			data += docmx.Find("#txt_gytdsyzh").Text() + "\t"
			data += docmx.Find("#txt_jsgcxkz").Text() + "\t"
			data += docmx.Find("#txt_sgxkz").Text() + "\t"
			data += docmx.Find("#txt_spfxkz").Text() + "\t"
			data += docmx.Find("#txt_kfqyzzzh").Text() + "\t"

			data += docmx.Find("#txt_kfqy").Text() + "\t"
			data += docmx.Find("#txt_lxdh").Text() + "\t"
			data += docmx.Find("#txt_dlgs").Text() + "\t"
			data += docmx.Find("#txt_dlgslxdh").Text() + "\t"
			data += docmx.Find("#txt_xmbajg").Text() + "\t"

			WirteTXT(data)
		}
	})
}

func convertUrlWithChineseToHex(url string) (rs string) {
	r := []rune(url)
	//fmt.Println(r)
	strSlice := []string{}
	rstrSlice := ""
	cnstr := ""
	for i := 0; i < len(r); i++ {
		if r[i] <= 40869 && r[i] >= 19968 {
			cnstr = cnstr + string(r[i])
			strSlice = append(strSlice, cnstr)

			mne := mahonia.NewEncoder("gbk")
			gbks := mne.ConvertString(string(r[i]))
			bgbks := []byte(gbks)

			for _, v := range bgbks {
				rstrSlice += "%"
				rstrSlice += fmt.Sprintf("%X", v)
			}

		} else {
			rstrSlice += string(r[i])
		}
		//fmt.Println("r[", i, "]=", r[i], "string=", string(r[i]))
	}
	if 0 == len(strSlice) {
		//无中文，需要跳过，后面再找规律
	}
	//fmt.Println("原字符串:", url, "\t提取出的中文字符串:", cnstr)
	//fmt.Println(rstrSlice)

	return rstrSlice
}

func WirteTXT(txt string) {
	f, err := os.OpenFile("fdc_data.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("os Create error: ", err)
		return
	}
	defer f.Close()

	bw := bufio.NewWriter(f)
	bw.WriteString(txt + "\n")
	bw.Flush()
}

func main() {
	//get_data()
	_, sec := get_doc()
	cp, tp := get_page_locate(get_doc())
	fmt.Printf("%d,%d\n", cp, tp)

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// navigate to a page, wait for an element, click
	var res string
	wd, _ := os.Getwd()
	fmt.Println(wd)
	set_ini_value("url_point", "run", "true")
	set_ini_value("progress", "totoalpage", strconv.Itoa(tp))
	for i := cp; i <= tp; i++ {
		fmt.Printf("读取%d/%d\n", i, tp)
		set_ini_value("progress", "currentpage", strconv.Itoa(i))
		err := chromedp.Run(ctx,
			chromedp.Navigate(sec["url"]),
			//chromedp.Sleep(3*time.Second), // 等待
			// wait for footer element is visible (ie, page is loaded)
			chromedp.WaitVisible("#form1", chromedp.ByQuery),
			// find and click "Expand All" link
			//chromedp.Sleep(10*time.Second), // 等待
			chromedp.SetValue(`#AspNetPager1_input`, strconv.Itoa(i), chromedp.ByID),
			//chromedp.Click(click_sect, chromedp.ByQuery),
			chromedp.Click(`#AspNetPager1_btn`, chromedp.ByID),
			//chromedp.Sleep(10*time.Second), // 等待
			// retrieve the value of the textarea
			chromedp.OuterHTML(sec["tables"], &res),
		)
		if err != nil {
			log.Fatal(err)
		}
		show_mx1(res, sec)

	}

}
