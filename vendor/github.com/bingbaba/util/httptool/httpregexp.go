package httptool

import (
	// "errors"
	//"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func MatchAll(reg *regexp.Regexp, body []byte) ([][]string, error) {
	var matched_strings [][]string

	//matched by regexp
	matched_bytes := reg.FindAllSubmatch(body, -1)
	for _, matched_items := range matched_bytes {
		matched_string := make([]string, 0)
		for i, matched_item := range matched_items {
			if i == 0 {
				continue
			}
			matched_string = append(matched_string, string(matched_item))
		}
		matched_strings = append(matched_strings, matched_string)
	}

	// if len(matched_strings) == 0 {
	// 	fmt.Println(string(body))
	// 	return matched_strings, errors.New("can't match the anyting by " + reg.String())
	// }

	return matched_strings, nil
}

func Match(reg *regexp.Regexp, body []byte) ([]string, error) {
	var matched_string []string

	//matched by regexp
	matched_bytes := reg.FindSubmatch(body)
	for i, matched_items := range matched_bytes {
		if i == 0 {
			continue
		}
		matched_string = append(matched_string, string(matched_items))
	}

	// if len(matched_string) == 0 {
	// 	fmt.Println(string(body))
	// 	return matched_string, errors.New("can't match the anyting!")
	// }
	return matched_string, nil
}

func HttpDoRegexp(req *http.Request, reg *regexp.Regexp) ([]string, error) {
	//http GET
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return nil, client_err
	}
	defer resp.Body.Close()

	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return nil, read_err
	}

	return Match(reg, resp_body)
}

func HttpDoAllRegexp(req *http.Request, reg *regexp.Regexp) ([][]string, error) {
	//http GET
	resp, client_err := httpclient.Do(req)
	if client_err != nil {
		return nil, client_err
	}
	defer resp.Body.Close()

	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return nil, read_err
	}

	return MatchAll(reg, resp_body)
}

func HttpGetAllRegexp(req_url string, reg *regexp.Regexp) ([][]string, error) {
	//http GET
	resp, client_err := httpclient.Get(req_url)
	if client_err != nil {
		return nil, client_err
	}
	defer resp.Body.Close()

	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return nil, read_err
	}

	return MatchAll(reg, resp_body)
}

func HttpGetRegexp(req_url string, reg *regexp.Regexp) ([]string, error) {
	//http GET
	resp, client_err := httpclient.Get(req_url)
	if client_err != nil {
		return nil, client_err
	}
	defer resp.Body.Close()

	//result
	resp_body, read_err := ioutil.ReadAll(resp.Body)
	if read_err != nil {
		return nil, read_err
	}

	return Match(reg, resp_body)
}
